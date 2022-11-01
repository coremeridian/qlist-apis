package resources

import (
	"context"
	"fmt"
	"strings"

	"github.com/Nerzal/gocloak/v11"
)

func CreateTest(ctx context.Context, kc *Keycloak) (rep *gocloak.ResourceRepresentation, err error) {
	token, err := kc.Client.LoginClient(ctx, kc.ClientID, kc.ClientSecret, kc.Realm)
	if err != nil {
		return
	}

	testTemplate, err := kc.Client.GetResourcesClient(ctx, token.AccessToken, kc.Realm, gocloak.GetResourceParams{
		Name: gocloak.StringP("Test"),
	})
	if err != nil {
		return
	}

	userInfo, err := kc.Client.GetUserInfo(ctx, kc.UserAccessToken, kc.Realm)
	if err != nil {
		return
	}

	resource := gocloak.ResourceRepresentation{
		Name:           gocloak.StringP(fmt.Sprintf("Test - %s", kc.ResourceName)),
		DisplayName:    gocloak.StringP(kc.ResourceName),
		ResourceScopes: testTemplate[0].ResourceScopes,
		Type:           gocloak.StringP(fmt.Sprintf("urn:%s:resources:test", strings.ToLower(kc.ClientID))),
		Owner: &gocloak.ResourceOwnerRepresentation{
			ID:   userInfo.Sub,
			Name: userInfo.Name,
		},
	}

	return kc.Client.CreateResourceClient(ctx, token.AccessToken, kc.Realm, resource)
}
