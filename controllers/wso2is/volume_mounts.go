package wso2is

import (
	"fmt"
	"github.com/wso2/k8s-wso2is-operator/variables"
	corev1 "k8s.io/api/core/v1"
)

func MakeVolumeMounts(version string) []corev1.VolumeMount {
	volumeMounts := []corev1.VolumeMount{
		{
			Name:      variables.PersistenVolumeName,
			MountPath: fmt.Sprintf("/home/wso2carbon/wso2is-%s/repository/deployment/server/userstores", version),
		},
		{
			Name:        variables.ConfigMapName,
			MountPath:   "/home/wso2carbon/wso2-config-volume/repository/conf/deployment.toml",
			SubPathExpr: "deployment.toml",
		},
		{
			Name:      variables.SecretName,
			MountPath: fmt.Sprintf("/home/wso2carbon/wso2is-%s/repository/resources/security/controller-keystores", version),
			ReadOnly:  true,
		},
	}

	return volumeMounts
}
