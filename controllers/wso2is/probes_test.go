package wso2is

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Probes", func() {
	It("should create the correct startup probe", func() {
		probe := MakeStartupProbe()

		// Verify ExecAction properties
		execAction := probe.ProbeHandler.Exec
		Expect(execAction).NotTo(BeNil(), "ExecAction is nil")
		expectedCommand := []string{
			"/bin/sh",
			"-c",
			"nc -z localhost 9443",
		}
		Expect(execAction.Command).To(Equal(expectedCommand), "Expected ExecAction Command %v, got %v", expectedCommand, execAction.Command)

		// Verify other Probe properties
		verifyProbeProperties(probe, 120, 10, 10)
	})

	It("should create the correct liveness probe", func() {
		probe := MakeLivenessProbe()

		// Verify ExecAction properties
		execAction := probe.ProbeHandler.Exec
		Expect(execAction).NotTo(BeNil(), "ExecAction is nil")
		expectedCommand := []string{
			"/bin/sh",
			"-c",
			"nc -z localhost 9443",
		}
		Expect(execAction.Command).To(Equal(expectedCommand), "Expected ExecAction Command %v, got %v", expectedCommand, execAction.Command)

		// Verify other Probe properties
		verifyProbeProperties(probe, 120, 10, 0)
	})

	It("should create the correct readiness probe", func() {
		probe := MakeReadinessProbe()

		// Verify HTTPGetAction properties
		httpGetAction := probe.ProbeHandler.HTTPGet
		Expect(httpGetAction).NotTo(BeNil(), "HTTPGetAction is nil")
		expectedPath := "/api/health-check/v1.0/health"
		Expect(httpGetAction.Path).To(Equal(expectedPath), "Expected HTTPGetAction Path %s, got %s", expectedPath, httpGetAction.Path)
		expectedPort := intstr.FromInt(9443)
		Expect(httpGetAction.Port).To(Equal(expectedPort), "Expected HTTPGetAction Port %d, got %d", expectedPort.IntValue(), httpGetAction.Port.IntValue())
		expectedScheme := corev1.URISchemeHTTPS
		Expect(httpGetAction.Scheme).To(Equal(expectedScheme), "Expected HTTPGetAction Scheme %s, got %s", expectedScheme, httpGetAction.Scheme)

		// Verify other Probe properties
		verifyProbeProperties(probe, 120, 10, 0)
	})
})

func verifyProbeProperties(probe *corev1.Probe, initialDelaySeconds, periodSeconds, failureThreshold int32) {
	Expect(probe.InitialDelaySeconds).To(Equal(initialDelaySeconds), "Expected InitialDelaySeconds %d, got %d", initialDelaySeconds, probe.InitialDelaySeconds)
	Expect(probe.PeriodSeconds).To(Equal(periodSeconds), "Expected PeriodSeconds %d, got %d", periodSeconds, probe.PeriodSeconds)
	Expect(probe.FailureThreshold).To(Equal(failureThreshold), "Expected FailureThreshold %d, got %d", failureThreshold, probe.FailureThreshold)
}
