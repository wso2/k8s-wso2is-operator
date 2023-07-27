package wso2is

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func MakeStartupProbe() *corev1.Probe {
	return &corev1.Probe{
		ProbeHandler: corev1.ProbeHandler{
			Exec: &corev1.ExecAction{
				Command: []string{
					"/bin/sh",
					"-c",
					"nc -z localhost 9443",
				},
			},
		},
		InitialDelaySeconds: 120,
		PeriodSeconds:       10,
		FailureThreshold:    10,
	}
}

func MakeLivenessProbe() *corev1.Probe {
	return &corev1.Probe{
		ProbeHandler: corev1.ProbeHandler{
			Exec: &corev1.ExecAction{
				Command: []string{
					"/bin/sh",
					"-c",
					"nc -z localhost 9443",
				},
			},
		},
		InitialDelaySeconds: 120,
		PeriodSeconds:       10,
	}
}

func MakeReadinessProbe() *corev1.Probe {
	return &corev1.Probe{
		ProbeHandler: corev1.ProbeHandler{
			HTTPGet: &corev1.HTTPGetAction{
				Path:   "/api/health-check/v1.0/health",
				Port:   intstr.FromInt(9443),
				Scheme: corev1.URISchemeHTTPS,
			},
		},
		InitialDelaySeconds: 120,
		PeriodSeconds:       10,
	}
}
