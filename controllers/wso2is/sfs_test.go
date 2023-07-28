package wso2is

import (
	"context"
	"github.com/go-logr/logr"
	wso2v1beta1 "github.com/wso2/k8s-wso2is-operator/api/v1beta1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("StatefulSet", func() {
	cpuRequests := "1"
	memoryRequests := "4096Mi"
	cpuLimits := "2"
	memoryLimits := "8000Mi"
	expectedReplicas := int32(3)

	instance := wso2v1beta1.Wso2Is{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "wso2is",
			Namespace: "wso2-iam-system",
		},
		Spec: wso2v1beta1.Wso2IsSpec{
			Version:        "v1.0",
			Size:           expectedReplicas,
			TomlConfigFile: "",
			Resources: wso2v1beta1.Resources{
				Requests: wso2v1beta1.ResourceRequests{
					Cpu:    cpuRequests,
					Memory: memoryRequests,
				},
				Limits: wso2v1beta1.ResourceLimits{
					Cpu:    cpuLimits,
					Memory: memoryLimits,
				},
			},
		},
	}

	It("should define the correct StatefulSet", func() {
		r := &Wso2IsReconciler{
			Scheme: runtime.NewScheme(),
		}

		statefulSet := r.defineStatefulSet(instance)

		// Verify ObjectMeta properties
		Expect(statefulSet.ObjectMeta.Name).To(Equal("wso2is"), "Expected StatefulSet Name wso2is, got %s", statefulSet.ObjectMeta.Name)
		Expect(statefulSet.ObjectMeta.Namespace).To(Equal("wso2-iam-system"), "Expected StatefulSet Namespace wso2-iam-system, got %s", statefulSet.ObjectMeta.Namespace)

		// Verify Spec properties
		Expect(*statefulSet.Spec.Replicas).To(Equal(expectedReplicas), "Expected StatefulSet Replicas %d, got %d", expectedReplicas, *statefulSet.Spec.Replicas)
		expectedServiceName := "wso2is-service"
		Expect(statefulSet.Spec.ServiceName).To(Equal(expectedServiceName), "Expected StatefulSet ServiceName %s, got %s", expectedServiceName, statefulSet.Spec.ServiceName)

		// Verify Template properties
		template := statefulSet.Spec.Template

		// Verify PodSpec properties
		podSpec := template.Spec

		// Verify Containers
		expectedContainerPorts := MakeContainerPorts()
		expectedEnv := MakeEnvVars()
		expectedVolumeMounts := MakeVolumeMounts("v1.0", instance)
		expectedStartupProbe := MakeStartupProbe()
		expectedLivenessProbe := MakeLivenessProbe()
		expectedReadinessProbe := MakeReadinessProbe()

		container := podSpec.Containers[0]
		Expect(container.Name).To(Equal("wso2is"), "Expected Container Name wso2is, got %s", container.Name)
		Expect(container.Image).To(Equal("rukshanjs/wso2is:v6.1.0"), "Expected Container Image rukshanjs/wso2is:v6.1.0, got %s", container.Image)
		Expect(container.Ports).To(Equal(expectedContainerPorts), "Expected Container Ports %v, got %v", expectedContainerPorts, container.Ports)
		Expect(container.Env).To(Equal(expectedEnv), "Expected Container Env %v, got %v", expectedEnv, container.Env)
		Expect(container.VolumeMounts).To(Equal(expectedVolumeMounts), "Expected Container VolumeMounts %v, got %v", expectedVolumeMounts, container.VolumeMounts)
		Expect(container.StartupProbe).To(Equal(expectedStartupProbe), "Expected Container StartupProbe %v, got %v", expectedStartupProbe, container.StartupProbe)
		Expect(container.LivenessProbe).To(Equal(expectedLivenessProbe), "Expected Container LivenessProbe %v, got %v", expectedLivenessProbe, container.LivenessProbe)
		Expect(container.ReadinessProbe).To(Equal(expectedReadinessProbe), "Expected Container ReadinessProbe %v, got %v", expectedReadinessProbe, container.ReadinessProbe)
		Expect(container.Lifecycle).To(Equal(MakeContainerLifecycle()), "Expected Container Lifecycle %v, got %v", MakeContainerLifecycle(), container.Lifecycle)
		Expect(container.ImagePullPolicy).To(Equal(corev1.PullIfNotPresent), "Expected Container ImagePullPolicy IfNotPresent, got %s", container.ImagePullPolicy)

		// Verify ServiceAccountName
		expectedServiceAccountName := "wso2is-serviceaccount"
		Expect(podSpec.ServiceAccountName).To(Equal(expectedServiceAccountName), "Expected PodSpec ServiceAccountName %s, got %s", expectedServiceAccountName, podSpec.ServiceAccountName)
	})

	Context("StatefulSet reconciliation", func() {

		scheme := runtime.NewScheme()
		_ = corev1.AddToScheme(scheme)
		_ = appsv1.AddToScheme(scheme)

		client := fake.NewClientBuilder().WithScheme(scheme).Build()

		logger := logr.Discard()

		// Prepare reconciler
		recon := &Wso2IsReconciler{
			Client: client,
			Log:    logger,
			Scheme: scheme,
		}

		ctx := context.TODO()

		It("Should reconcile a StatefulSet", func() {
			// Create and reconcile the ConfigMap
			_, err := reconcileConfigMap(recon, instance, logger, nil, ctx)
			Expect(err).NotTo(HaveOccurred(), "Error during ConfigMap reconciliation")

			_, err = reconcileStatefulSet(recon, instance, logger, nil, ctx)
			Expect(err).NotTo(HaveOccurred(), "Error during StatefulSet reconciliation")

			statefulSet := &appsv1.StatefulSet{}
			err = recon.Client.Get(ctx, types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, statefulSet)
			Expect(err).NotTo(HaveOccurred(), "StatefulSet not got")
		})
	})

})
