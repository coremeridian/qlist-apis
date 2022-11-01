package service

import (
	"context"
	"strings"

	"coremeridian.xyz/app/qlist/cmd/api/services/tests/models"
	"coremeridian.xyz/app/qlist/cmd/api/services/tests/service/resources"
	"github.com/Nerzal/gocloak/v11"
	"github.com/costal/go-misc-tools/httpapp"
)

type HTTPApplication struct {
	HTTP *httpapp.Application
	Endpoints
}

type Endpoints struct {
	Tests interface {
		Insert(*models.Test) (interface{}, error)
		Update(string, interface{}) (didUpdate bool, err error)
		Get(id string) (*models.Test, error)
		Latest(int, models.TestOptions) ([]*models.Test, error)
	}
	TestSessions interface {
		Insert(*models.TestSession) (interface{}, error)
		Update(string, interface{}) (didUpdate bool, err error)
		Get(keys *models.KeySet) (*models.TestSession, error)
		Latest(int) ([]*models.TestSession, error)
	}
}

func (app *HTTPApplication) Init() {
	app.HTTP.Authorize = app.setAuthorizationFunc()
}

func (app *HTTPApplication) setAuthorizationFunc() func(context.Context, ...interface{}) (interface{}, error) {
	client := gocloak.NewClient(app.HTTP.Domain())
	return func(ctx context.Context, keys ...interface{}) (result interface{}, err error) {
		actionKey := keys[0].(string)
		switch actionKey {
		case "RetrospectToken":
			accessToken := keys[1].(string)
			result, err = client.RetrospectToken(ctx, accessToken, app.HTTP.ClientID, app.HTTP.ClientSecret, app.HTTP.Realm)
		case "RPT:perms":
			accessToken := keys[1].(string)
			options := gocloak.RequestingPartyTokenOptions{
				Audience:     gocloak.StringP(app.HTTP.ClientID),
				Permissions:  &[]string{"Test"},
				ResponseMode: gocloak.StringP("permissions"),
			}
			result, err = client.GetRequestingPartyPermissions(ctx, accessToken, app.HTTP.Realm, options)
		case "resource:create:Test":
			testName := keys[1].(string)
			accessToken := keys[2].(string)
			result, err = app.serverResources(ctx, client, "create", "Test:"+testName, accessToken)
		}
		if err != nil {
			return nil, err
		}
		return
	}
}

func (app *HTTPApplication) serverResources(ctx context.Context, client gocloak.GoCloak, action, resource, accessToken string) (result interface{}, err error) {
	keycloak := resources.Keycloak{
		Client:          client,
		ClientID:        app.HTTP.ClientID,
		ClientSecret:    app.HTTP.ClientSecret,
		Realm:           app.HTTP.Realm,
		UserAccessToken: accessToken,
		ResourceName:    resource[strings.LastIndex(resource, ":")+1:],
	}
	switch action {
	case "create":
		result, err = resources.Create(ctx, resource, &keycloak)
	}

	return
}
