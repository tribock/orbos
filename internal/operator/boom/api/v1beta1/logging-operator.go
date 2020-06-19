package v1beta1

import "github.com/caos/orbos/internal/operator/boom/api/v1beta1/storage"

type LoggingOperator struct {
	//Flag if tool should be deployed
	//@default: false
	Deploy bool `json:"deploy" yaml:"deploy"`
	//Spec to define how the persistence should be handled
	FluentdPVC *storage.Spec `json:"fluentdStorage,omitempty" yaml:"fluentdStorage,omitempty"`
}
