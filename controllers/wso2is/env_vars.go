package wso2is

import (
	corev1 "k8s.io/api/core/v1"
)

func MakeEnvVars() []corev1.EnvVar {
	envVars := []corev1.EnvVar{
		{
			Name: "NODE_IP",
			ValueFrom: &corev1.EnvVarSource{
				FieldRef: &corev1.ObjectFieldSelector{
					FieldPath: "status.podIP",
				},
			},
		},
	}

	return envVars
}
