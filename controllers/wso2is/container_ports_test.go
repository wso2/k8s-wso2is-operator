package wso2is

import (
	corev1 "k8s.io/api/core/v1"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ContainerPorts", func() {
	It("should create the correct container ports", func() {
		containerPorts := MakeContainerPorts()

		// Verify the length of the containerPorts slice
		expectedLength := 2
		Expect(len(containerPorts)).To(Equal(expectedLength))

		// Verify the ContainerPort and Protocol of each ContainerPort
		expectedContainerPorts := []corev1.ContainerPort{
			{
				ContainerPort: int32(9443),
				Protocol:      "TCP",
			},
			{
				ContainerPort: int32(9763),
				Protocol:      "TCP",
			},
		}
		for i, expectedContainerPort := range expectedContainerPorts {
			containerPort := containerPorts[i]
			Expect(containerPort.ContainerPort).To(Equal(expectedContainerPort.ContainerPort))
			Expect(containerPort.Protocol).To(Equal(expectedContainerPort.Protocol))
		}
	})
})
