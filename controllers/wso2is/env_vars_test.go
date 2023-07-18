package wso2is

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
)

var _ = Describe("EnvVars", func() {
	It("should create the correct environment variables", func() {
		envVars := MakeEnvVars()

		// Verify the length of the envVars slice
		expectedLength := 1
		Expect(len(envVars)).To(Equal(expectedLength))

		// Verify the content of the envVar
		expectedEnvVar := corev1.EnvVar{
			Name: "NODE_IP",
			ValueFrom: &corev1.EnvVarSource{
				FieldRef: &corev1.ObjectFieldSelector{
					FieldPath: "status.podIP",
				},
			},
		}
		envVar := envVars[0]
		Expect(envVar).To(Equal(expectedEnvVar))
	})
})
