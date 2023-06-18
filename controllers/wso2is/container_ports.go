package wso2is

import (
	corev1 "k8s.io/api/core/v1"
)

func MakeContainerPorts() []corev1.ContainerPort {
	containerPortHttp := int32(9763)
	containerPortHttps := int32(9443)

	containerPorts := []corev1.ContainerPort{
		{
			ContainerPort: containerPortHttps,
			Protocol:      "TCP",
		},
		{
			ContainerPort: containerPortHttp,
			Protocol:      "TCP",
		},
	}

	return containerPorts
}
