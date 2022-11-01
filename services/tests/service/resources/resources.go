// Package resources manages resources via keycloak
package resources

import (
	"context"
	"strings"

	"github.com/Nerzal/gocloak/v11"
)

type Keycloak struct {
	Client          gocloak.GoCloak
	ClientID        string
	ClientSecret    string
	Realm           string
	UserAccessToken string
	ResourceName    string
}

func Create(ctx context.Context, resource string, kc *Keycloak) (rep *gocloak.ResourceRepresentation, err error) {
	// Need to pool or cache redundant access to keycloak to limit network
	switch {
	case strings.Contains(resource, "Test"):
		rep, err = CreateTest(ctx, kc)
	}

	return
}

func getScopesFromReps(reps []*gocloak.ScopeRepresentation, scopes ...string) *[]gocloak.ScopeRepresentation {
	var matchedReps *[]gocloak.ScopeRepresentation
	for _, rep := range reps {
		for i := 0; i < len(reps); i++ {
			if scopes == nil || *rep.Name == scopes[i] {
				*matchedReps = append(*matchedReps, *rep)
				break
			}
		}
	}
	return matchedReps
}
