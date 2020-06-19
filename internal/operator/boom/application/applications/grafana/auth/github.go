package auth

import (
	github "github.com/caos/orbos/internal/operator/boom/api/v1beta1/grafana/auth/Github"
	helper2 "github.com/caos/orbos/internal/utils/helper"
	"strings"
)

func GetGithubAuthConfig(spec *github.Auth) (map[string]string, error) {
	clientID, err := helper2.GetSecretValue(spec.ClientID, spec.ExistingClientIDSecret)
	if err != nil {
		return nil, err
	}

	clientSecret, err := helper2.GetSecretValue(spec.ClientSecret, spec.ExistingClientSecretSecret)
	if err != nil {
		return nil, err
	}

	if clientID == "" || clientSecret == "" {
		return nil, nil
	}

	teamIds := strings.Join(spec.TeamIDs, " ")
	allowedOrganizations := strings.Join(spec.AllowedOrganizations, " ")

	return map[string]string{
		"enabled":               "true",
		"allow_sign_up":         "true",
		"client_id":             clientID,
		"client_secret":         clientSecret,
		"scopes":                "user:email,read:org",
		"auth_url":              "https://github.com/login/oauth/authorize",
		"token_url":             "https://github.com/login/oauth/access_token",
		"api_url":               "https://api.github.com/user",
		"team_ids":              teamIds,
		"allowed_organizations": allowedOrganizations,
	}, nil
}
