package wso2is

import (
	"context"
	"github.com/go-logr/logr"
	wso2v1beta1 "github.com/wso2/k8s-wso2is-operator/api/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Secret", func() {
	It("should define the correct Secret", func() {
		r := &Wso2IsReconciler{
			Scheme: runtime.NewScheme(),
		}

		instance := wso2v1beta1.Wso2Is{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "wso2is",
				Namespace: "wso2-iam-system",
			},

			Spec: wso2v1beta1.Wso2IsSpec{
				KeystoreMounts: []wso2v1beta1.KeystoreMount{
					{
						Name: "keystore1",
						Data: "key1",
					},
					{
						Name: "keystore2",
						Data: "key2",
					},
				},
			},
		}

		secret := r.defineSecret(instance)

		// Verify ObjectMeta properties
		Expect(secret.ObjectMeta.Name).To(Equal("wso2is-secret"), "Expected Secret Name wso2is-secret, got %s", secret.ObjectMeta.Name)
		Expect(secret.ObjectMeta.Namespace).To(Equal("wso2-iam-system"), "Expected Secret Namespace wso2-iam-system, got %s", secret.ObjectMeta.Namespace)

		// Verify Data
		expectedData := map[string][]byte{
			"keystore1": []byte("key1"),
			"keystore2": []byte("key2"),
		}
		Expect(secret.Data).To(Equal(expectedData), "Expected Secret Data %v, got %v", expectedData, secret.Data)
	})

	Context("Secret reconciliation", func() {
		scheme := runtime.NewScheme()
		_ = corev1.AddToScheme(scheme)

		client := fake.NewClientBuilder().WithScheme(scheme).Build()

		logger := logr.Discard()

		// Prepare reconciler
		recon := &Wso2IsReconciler{
			Client: client,
			Log:    logger,
			Scheme: scheme,
		}

		instance := wso2v1beta1.Wso2Is{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "wso2is",
				Namespace: "wso2-iam-system",
			},
			Spec: wso2v1beta1.Wso2IsSpec{
				TomlConfigFile: "",
			},
		}

		ctx := context.TODO()

		It("Should reconcile a Secret", func() {

			// Create and reconcile the Secret
			_, err := reconcileSecret(recon, instance, logger, nil, ctx)
			Expect(err).NotTo(HaveOccurred(), "Error during Secret reconciliation")

			secret := &corev1.Secret{}
			err = recon.Client.Get(ctx, types.NamespacedName{Name: instance.Name + "-secret", Namespace: instance.Namespace}, secret)
			Expect(err).NotTo(HaveOccurred(), "Secret not got")
		})
	})

})
