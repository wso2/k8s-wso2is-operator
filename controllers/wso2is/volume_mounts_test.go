package wso2is

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	wso2v1beta1 "github.com/wso2/k8s-wso2is-operator/api/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("MakeVolumeMounts", func() {
	It("should create volume mounts", func() {
		version := "v1.0"

		instance := wso2v1beta1.Wso2Is{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "wso2is",
				Namespace: "wso2-iam-system",
			},
			Spec: wso2v1beta1.Wso2IsSpec{},
		}

		volumeMounts := MakeVolumeMounts(version, instance)

		// Verify the length of the volumeMounts slice
		expectedLength := 3
		Expect(len(volumeMounts)).To(Equal(expectedLength))

		// Verify the content of each volumeMount
		expectedVolumeMounts := []corev1.VolumeMount{
			{
				Name:      "wso2is-pv",
				MountPath: fmt.Sprintf("/home/wso2carbon/wso2is-%s/repository/deployment/server/userstores", version),
			},
			{
				Name:        "wso2is-config",
				MountPath:   "/home/wso2carbon/wso2-config-volume/repository/conf/deployment.toml",
				SubPathExpr: "deployment.toml",
			},
			{
				Name:      "wso2is-secret",
				MountPath: fmt.Sprintf("/home/wso2carbon/wso2is-%s/repository/resources/security/controller-keystores", version),
				ReadOnly:  true,
			},
		}
		Expect(volumeMounts).To(Equal(expectedVolumeMounts))
	})
})
