package resources

import (
	"context"
	"fmt"
	"strings"

	"github.com/Nerzal/gocloak/v11"
)

func CreateTestSession(kc *Keycloak) {
	ctx := context.Background()
	token, err := kc.Client.LoginAdmin(ctx, "admin", "pass", kc.Realm)
	if err != nil {
		return
	}

	viewScope, err := kc.Client.GetScope(ctx, token.AccessToken, kc.Realm, kc.ClientID, "view")
	if err != nil {
		return
	}

	tsViewScope, err := kc.Client.GetScope(ctx, token.AccessToken, kc.Realm, kc.ClientID, "testsession:view")
	if err != nil {
		return
	}

	userInfo, err := kc.Client.GetUserInfo(ctx, kc.UserAccessToken, kc.Realm)
	if err != nil {
		return
	}

	resource := gocloak.ResourceRepresentation{
		Name:           gocloak.StringP("TestSession"),
		DisplayName:    gocloak.StringP("Test Session"),
		ResourceScopes: &[]gocloak.ScopeRepresentation{*viewScope, *tsViewScope},
		Type:           gocloak.StringP(fmt.Sprintf("urn:%s:resources:testsession", strings.ToLower(kc.ClientID))),
		Owner: &gocloak.ResourceOwnerRepresentation{
			ID:   userInfo.Sub,
			Name: userInfo.Name,
		},
	}
	kc.Client.CreateResource(ctx, token.AccessToken, kc.Realm, kc.ClientID, resource)
}
