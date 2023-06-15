package wso2is

import (
	"github.com/wso2/k8s-wso2is-operator/variables"
	corev1 "k8s.io/api/core/v1"
)

func MakeContainerPorts() []corev1.ContainerPort {
	containerPorts := []corev1.ContainerPort{
		{
			ContainerPort: variables.ContainerPortHttps,
			Protocol:      "TCP",
		},
		{
			ContainerPort: variables.ContainerPortHttp,
			Protocol:      "TCP",
		},
	}

	return containerPorts
}
