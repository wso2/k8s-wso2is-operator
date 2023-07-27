package wso2is

import (
	wso2v1beta1 "github.com/wso2/k8s-wso2is-operator/api/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

func MakeResourceRequirements(m wso2v1beta1.Wso2Is) corev1.ResourceRequirements {
	resourceReqs := corev1.ResourceRequirements{
		Requests: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse(m.Spec.Resources.Requests.Cpu),
			corev1.ResourceMemory: resource.MustParse(m.Spec.Resources.Requests.Memory),
		},
		Limits: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse(m.Spec.Resources.Limits.Cpu),
			corev1.ResourceMemory: resource.MustParse(m.Spec.Resources.Limits.Memory),
		},
	}

	return resourceReqs
}
