package wso2is

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

func MakeResourceRequirements() corev1.ResourceRequirements {
	resourceReqs := corev1.ResourceRequirements{
		Requests: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse("4Gi"),
			corev1.ResourceMemory: resource.MustParse("4096m"),
		},
		Limits: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse("4Gi"),
			corev1.ResourceMemory: resource.MustParse("4096m"),
		},
	}

	return resourceReqs
}