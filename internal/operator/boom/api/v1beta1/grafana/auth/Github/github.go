package github

import (
	"github.com/caos/orbos/internal/secret"
	"reflect"
)

type Auth struct {
	ClientID *secret.Secret `yaml:"clientID,omitempty"`
	//Existing secret with the clientID
	ExistingClientIDSecret *secret.Existing `json:"existingClientIDSecret,omitempty" yaml:"existingClientIDSecret,omitempty"`
	ClientSecret           *secret.Secret   `yaml:"clientSecret,omitempty"`
	//Existing secret with the clientSecret
	ExistingClientSecretSecret *secret.Existing `json:"existingClientSecretSecret,omitempty" yaml:"existingClientSecretSecret,omitempty"`
	// Organizations allowed to login
	AllowedOrganizations []string `json:"allowedOrganizations,omitempty" yaml:"allowedOrganizations,omitempty"`
	// TeamIDs where the user is required to have at least one membership
	TeamIDs []string `json:"teamIDs,omitempty" yaml:"teamIDs,omitempty"`
}

func ClearEmpty(g *Auth) *Auth {
	if g == nil {
		return nil
	}

	marshaled := Auth{
		ClientID:                   secret.ClearEmpty(g.ClientID),
		ExistingClientIDSecret:     g.ExistingClientIDSecret,
		ClientSecret:               secret.ClearEmpty(g.ClientSecret),
		ExistingClientSecretSecret: g.ExistingClientSecretSecret,
		AllowedOrganizations:       g.AllowedOrganizations,
		TeamIDs:                    g.TeamIDs,
	}
	if reflect.DeepEqual(marshaled, Auth{}) {
		return nil
	}
	return &marshaled
}

func (g *Auth) MarshalYAML() (interface{}, error) {
	type Alias Auth
	return &Alias{
		ClientID:                   secret.ClearEmpty(g.ClientID),
		ExistingClientIDSecret:     g.ExistingClientIDSecret,
		ClientSecret:               secret.ClearEmpty(g.ClientSecret),
		ExistingClientSecretSecret: g.ExistingClientSecretSecret,
		AllowedOrganizations:       g.AllowedOrganizations,
		TeamIDs:                    g.TeamIDs,
	}, nil
}
