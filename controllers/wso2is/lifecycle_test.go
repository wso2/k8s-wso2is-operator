package wso2is

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ContainerLifecycle", func() {
	It("should create the correct container lifecycle", func() {
		lifecycle := MakeContainerLifecycle()

		// Verify PreStop properties
		preStop := lifecycle.PreStop
		Expect(preStop).NotTo(BeNil(), "PreStop is nil")

		// Verify ExecAction properties
		execAction := preStop.Exec
		Expect(execAction).NotTo(BeNil(), "ExecAction is nil")
		expectedCommand := []string{
			"sh",
			"-c",
			"${WSO2_SERVER_HOME}/bin/wso2server.sh stop",
		}
		Expect(execAction.Command).To(Equal(expectedCommand), "Expected ExecAction Command %v, got %v", expectedCommand, execAction.Command)
	})
})
