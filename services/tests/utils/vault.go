package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type vault map[string]interface{}

var secrets vault

func GetAPIVault() (vault, error) {
	if secrets != nil {
		return secrets, nil
	}

	client := &http.Client{}

	data := struct {
		RoleID   string `json:"role_id"`
		SecretID string `json:"secret_id"`
	}{
		RoleID:   os.Getenv("APP_ROLE_ID"),
		SecretID: os.Getenv("APP_SECRET_ID"),
	}

	body, _ := json.Marshal(data)
	req, err := http.NewRequest("POST", os.Getenv("VAULT_LOGIN_URL"), bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}

	apiKey := os.Getenv("API_KEY")
	req.Header.Set("X-API-Key", apiKey)
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	} else if resp.StatusCode != 200 {
		return nil, fmt.Errorf("no token provided: %d", resp.StatusCode)
	}

	defer resp.Body.Close()
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var vaultAccess struct {
		Result struct {
			Token string `json:"client_token"`
		} `json:"auth"`
	}

	err = json.Unmarshal(body, &vaultAccess)
	if err != nil {
		fmt.Println(string(body), os.Getenv("API_KEY"))
		panic(err)
	}

	return getAPISecrets(client, vaultAccess.Result.Token)
}

func getAPISecrets(client *http.Client, token string) (vault, error) {
	req, err := http.NewRequest("GET", os.Getenv("VAULT_SECRET_URL"), nil)
	if err != nil {
		panic(err)
	}

	req.Header.Add("X-API-Key", os.Getenv("API_KEY"))
	req.Header.Add("X-Vault-Token", token)

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	} else if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		defer resp.Body.Close()
		return nil, fmt.Errorf("no secrets provided: %s - %d", string(body), resp.StatusCode)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var data struct {
		Secrets map[string]interface{} `json:"data"`
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		panic(err)
	}
	secrets = data.Secrets

	return secrets, err
}
