package wso2is

import (
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	wso2v1beta1 "github.com/wso2/k8s-wso2is-operator/api/v1beta1"
	corev1 "k8s.io/api/core/v1"

	logr "github.com/go-logr/logr"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

var _ = Describe("updateStatus", func() {
	var (
		ctx        context.Context
		log        logr.Logger
		instance   wso2v1beta1.Wso2Is
		reconciler Wso2IsReconciler
		testClient client.Client
	)

	BeforeEach(func() {
		// Prepare the context, logger, and instance
		ctx = context.Background()
		log = logr.Logger{}
		instance = wso2v1beta1.Wso2Is{
			// TODO: fill with necessary fields
		}

		// Prepare the reconciler
		s := runtime.NewScheme()
		_ = scheme.AddToScheme(s)
		_ = wso2v1beta1.AddToScheme(s)

		testClient = fake.NewFakeClientWithScheme(s)
		reconciler = Wso2IsReconciler{
			Client: testClient,
			Log:    log,
			Scheme: s,
		}
	})

	Context("When updating the status", func() {
		It("should not return an error if pods are listed successfully", func() {
			// Mock the list of pods returned by client
			pod := &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-pod",
					Namespace: "wso2-iam-system",
					Labels:    map[string]string{"app": "wso2is"},
				},
			}
			err := testClient.Create(ctx, pod)
			Expect(err).ToNot(HaveOccurred())

			_, err = updateStatus(&reconciler, instance, log, nil, ctx)
			Expect(err).ToNot(HaveOccurred())
		})
	})
})
