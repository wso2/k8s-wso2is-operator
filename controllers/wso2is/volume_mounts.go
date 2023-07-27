package wso2is

import (
	"fmt"
	wso2v1beta1 "github.com/wso2/k8s-wso2is-operator/api/v1beta1"
	corev1 "k8s.io/api/core/v1"
)

func MakeVolumeMounts(version string, instance wso2v1beta1.Wso2Is) []corev1.VolumeMount {

	volumeMounts := []corev1.VolumeMount{
		{
			Name:      instance.Name + "-pv",
			MountPath: fmt.Sprintf("/home/wso2carbon/wso2is-%s/repository/deployment/server/userstores", version),
		},
		{
			Name:        instance.Name + "-config",
			MountPath:   "/home/wso2carbon/wso2-config-volume/repository/conf/deployment.toml",
			SubPathExpr: "deployment.toml",
		},
		{
			Name:      instance.Name + "-secret",
			MountPath: fmt.Sprintf("/home/wso2carbon/wso2is-%s/repository/resources/security/controller-keystores", version),
			ReadOnly:  true,
		},
	}

	return volumeMounts
}
