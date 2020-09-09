package cloudflare

import (
	"github.com/caos/orbos/internal/secret"
	"github.com/caos/orbos/internal/tree"
	"github.com/caos/orbos/mntr"
	"github.com/pkg/errors"
)

func SecretsFunc() secret.Func {
	return func(monitor mntr.Monitor, desiredTree *tree.Tree) (secrets map[string]*secret.Secret, err error) {
		defer func() {
			err = errors.Wrapf(err, "building %s failed", desiredTree.Common.Kind)
		}()

		desiredKind, err := parseDesired(desiredTree)
		if err != nil {
			return nil, errors.Wrap(err, "parsing desired state failed")
		}
		desiredTree.Parsed = desiredKind

		return getSecretsMap(desiredKind), nil
	}
}

func getSecretsMap(desiredKind *Desired) map[string]*secret.Secret {
	secrets := map[string]*secret.Secret{}
	return secrets
} /*
	if desiredKind.Spec != nil && desiredKind.Spec.Configuration != nil {
		conf := desiredKind.Spec.Configuration
		if conf.ConsoleEnvironmentJSON == nil {
			conf.ConsoleEnvironmentJSON = &secret.Secret{}
		}
		secrets["consoleenvironmentjson"] = conf.ConsoleEnvironmentJSON

		if conf.Tracing != nil {
			if conf.Tracing.ServiceAccountJSON == nil {
				conf.Tracing.ServiceAccountJSON = &secret.Secret{}
			}
			secrets["serviceaccountjson"] = conf.Tracing.ServiceAccountJSON
		}

		if conf.Secrets != nil {
			if conf.Secrets.Keys == nil {
				conf.Secrets.Keys = &secret.Secret{}
			}
			secrets["keys"] = conf.Secrets.Keys
		}

		if conf.Notifications != nil {
			if conf.Notifications.GoogleChatURL == nil {
				conf.Notifications.GoogleChatURL = &secret.Secret{}
			}
			secrets["googlechaturl"] = conf.Notifications.GoogleChatURL

			if conf.Notifications.Twilio.SID == nil {
				conf.Notifications.Twilio.SID = &secret.Secret{}
			}
			secrets["twiliosid"] = conf.Notifications.Twilio.SID

			if conf.Notifications.Twilio.AuthToken == nil {
				conf.Notifications.Twilio.AuthToken = &secret.Secret{}
			}
			secrets["twilioauthtoken"] = conf.Notifications.Twilio.AuthToken

			if conf.Notifications.Email.AppKey == nil {
				conf.Notifications.Email.AppKey = &secret.Secret{}
			}
			secrets["emailappkey"] = conf.Notifications.Email.AppKey
		}
	}
	return secrets
}
*/
