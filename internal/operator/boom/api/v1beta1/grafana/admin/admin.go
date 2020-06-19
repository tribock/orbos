package admin

import (
	"github.com/caos/orbos/internal/secret"
	"reflect"
)

// Admin: Not defining the admin credentials results in an user admin with password admin.
type Admin struct {
	Username *secret.Secret `json:"username,omitempty" yaml:"username,omitempty"`
	Password *secret.Secret `yaml:"password,omitempty"`
	//Existing Secret containing username and password
	ExistingSecret *secret.ExistingIDSecret `json:"existingSecret,omitempty" yaml:"existingSecret,omitempty"`
}

func (a *Admin) MarshalYAML() (interface{}, error) {
	type Alias Admin
	return &Alias{
		Username:       secret.ClearEmpty(a.Username),
		Password:       secret.ClearEmpty(a.Password),
		ExistingSecret: a.ExistingSecret,
	}, nil
}

func ClearEmpty(x *Admin) *Admin {
	if x == nil {
		return nil
	}

	marshaled := Admin{
		Username:       secret.ClearEmpty(x.Username),
		Password:       secret.ClearEmpty(x.Password),
		ExistingSecret: x.ExistingSecret,
	}

	if reflect.DeepEqual(marshaled, Admin{}) {
		return nil
	}
	return &marshaled
}
