package validator

import (
    "fmt"
    adm "k8s.io/api/admission/v1beta1"
    "github.com/kelseyhightower/envconfig"
)


type ReplicasValidator struct {
    Max int32
}

func NewReplicasValidator() ReplicasValidator {
    v := ReplicasValidator{}
    envconfig.MustProcess("replicas", &v)
    return v
}

func (v ReplicasValidator) Validate(req *adm.AdmissionRequest) *adm.AdmissionResponse {
    dep, ok := GetDeployment(req)
    if !ok {
        return NewResponse(true, "ok")
    }

    replicas := dep.Spec.Replicas

    if *replicas > v.Max {
        return NewResponse(false, fmt.Sprintf("replica count too high"))
    }

    return NewResponse(true, "ok")
}
