package wso2is

import (
	corev1 "k8s.io/api/core/v1"
)

func MakeContainerLifecycle() *corev1.Lifecycle {
	return &corev1.Lifecycle{
		PreStop: &corev1.LifecycleHandler{
			Exec: &corev1.ExecAction{
				Command: []string{
					"sh",
					"-c",
					"${WSO2_SERVER_HOME}/bin/wso2server.sh stop",
				},
			},
		},
	}
}
