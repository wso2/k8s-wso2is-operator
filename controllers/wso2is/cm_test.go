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
)

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("WSO2 IS", func() {
	It("should return the correct ConfigMapName when TomlConfigFile is not set", func() {
		instance := wso2v1beta1.Wso2Is{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "wso2is",
				Namespace: "wso2-iam-system",
			},
		}

		configMapName := getConfigMapName(instance)
		expectedConfigMapName := instance.Name + "-config"
		Expect(configMapName).To(Equal(expectedConfigMapName))
	})

	It("should return the correct ConfigMapName when TomlConfigFile is set", func() {
		instance := wso2v1beta1.Wso2Is{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "wso2is",
				Namespace: "wso2-iam-system",
			},
			Spec: wso2v1beta1.Wso2IsSpec{
				TomlConfigFile: "custom-config",
			},
		}

		configMapName := getConfigMapName(instance)
		expectedConfigMapName := "custom-config"
		Expect(configMapName).To(Equal(expectedConfigMapName))
	})

	It("should calculate the ConfigMap hash without error", func() {
		configMap := &corev1.ConfigMap{}

		hash, err := calculateConfigMapHash(configMap)
		Expect(err).To(BeNil())
		Expect(hash).NotTo(BeEmpty())
	})

	It("should return the correct volume index", func() {
		volumes := []corev1.Volume{
			{
				Name: "volume1",
			},
			{
				Name: "volume2",
			},
		}

		// Check for a volume that exists
		volumeIndex := getVolumeIndex(volumes, "volume1")
		expectedVolumeIndex := 0
		Expect(volumeIndex).To(Equal(expectedVolumeIndex))

		// Check for a volume that does not exist
		volumeIndex = getVolumeIndex(volumes, "volume3")
		expectedVolumeIndex = -1
		Expect(volumeIndex).To(Equal(expectedVolumeIndex))
	})

	It("should define the ConfigMap", func() {
		// Prepare scheme for the fake client
		scheme := runtime.NewScheme()
		_ = corev1.AddToScheme(scheme)

		// Create an instance of your reconciler with the fake client
		r := &Wso2IsReconciler{
			Client: fake.NewClientBuilder().WithScheme(scheme).Build(),
			Scheme: scheme,
		}

		m := wso2v1beta1.Wso2Is{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "wso2is",
				Namespace: "wso2-iam-system",
			},
		}

		cm := r.defineConfigMap(m)
		Expect(cm).NotTo(BeNil(), "ConfigMap is nil")
	})

	Context("ConfigMap reconciliation", func() {

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

		It("Should recpncile a ConfigMap", func() {
			_, err := reconcileConfigMap(recon, instance, logger, nil, ctx)
			Expect(err).NotTo(HaveOccurred(), "Error during ConfigMap reconciliation")

			configMap := &corev1.ConfigMap{}
			err = recon.Client.Get(ctx, types.NamespacedName{Name: instance.Name + "-config", Namespace: instance.Namespace}, configMap)
			Expect(err).NotTo(HaveOccurred(), "ConfigMap not got")
		})

	})
})
