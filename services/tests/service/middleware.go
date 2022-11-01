package service

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/Nerzal/gocloak/v11"
)

type contextKey string

const accessTokenContext = contextKey("accessToken")

func (app *HTTPApplication) allowedOriginValidator(r *http.Request, origin string) bool {
	_, ok := CorsManager.allowedOrigins[origin]
	return ok
}

func extractToken(r *http.Request) (string, error) {
	headerToken := r.Header.Get("Authorization")
	bodyToken := r.FormValue("access_token")

	if headerToken != "" {
		authHeader := strings.Split(headerToken, " ")
		if len(authHeader) != 2 {
			return "", fmt.Errorf("invalid Authorization header format")
		}
		return authHeader[1], nil
	} else if bodyToken != "" {
		return bodyToken, nil
	} else {
		return "", fmt.Errorf("no access token")
	}

}

func (app *HTTPApplication) JWTGuard(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.HTTP.InfoLog().Println("Middleware - JWTGuard")
		accessToken, err := extractToken(r)
		if err != nil {
			app.HTTP.ClientError(w, http.StatusForbidden)
			return
		}

		result, err := app.HTTP.Authorize(r.Context(), "RetrospectToken", accessToken)
		if result == nil {
			if err != nil {
				app.HTTP.ErrorLog().Println(err)
			}
			app.HTTP.ClientError(w, http.StatusForbidden)
			return
		}
		rptResult := result.(*gocloak.RetrospecTokenResult)
		if rptResult.Active != nil && !*rptResult.Active {
			app.HTTP.ClientError(w, http.StatusForbidden)
			return
		}

		var ctx context.Context
		if s, ok := r.Context().Value(accessTokenContext).(string); ok {
			fmt.Println(s)
			ctx = context.WithValue(r.Context(), accessTokenContext, "1")
		} else {
			ctx = context.WithValue(r.Context(), accessTokenContext, accessToken)
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *HTTPApplication) hasResourceScope(ctx context.Context, accessToken string, scopes ...string) (bool, error) {
	result, err := app.HTTP.Authorize(ctx, "RPT:perms", accessToken)
	if err == nil {
		userPerms := result.(*[]gocloak.RequestingPartyPermission)
		userScopes := strings.Join(*(*userPerms)[0].Scopes, " ")
		for _, resourceScope := range scopes {
			if strings.Contains(userScopes, resourceScope) {
				return true, nil
			}
		}
	}

	return false, err
}

func (app *HTTPApplication) withPermissionScope(scope string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			accessToken := r.Context().Value(accessTokenContext).(string)
			hasScope, err := app.hasResourceScope(r.Context(), accessToken, scope)
			if !hasScope || err != nil {
				app.HTTP.ClientError(w, http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func (app *HTTPApplication) withViewConstraints(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessToken := r.Context().Value(accessTokenContext).(string)
		hasModifiableView, err := app.hasResourceScope(r.Context(), accessToken, "view-unpublished")
		if err != nil {
			app.HTTP.ServerError(w, err)
			return
		}
		if !hasModifiableView && len(r.URL.Query()) > 0 && r.ContentLength > 0 {
			app.HTTP.ClientError(w, http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
