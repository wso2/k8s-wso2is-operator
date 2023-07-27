package wso2is

import (
	"context"
	"github.com/go-logr/logr"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	wso2v1beta1 "github.com/wso2/k8s-wso2is-operator/api/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

var _ = Describe("MakeVolumes", func() {
	It("should create volumes", func() {
		instance := wso2v1beta1.Wso2Is{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "wso2is",
				Namespace: "wso2-iam-system",
			},
		}

		volumes := MakeVolumes(instance)

		// Verify the length of the volumes slice
		expectedLength := 3
		Expect(len(volumes)).To(Equal(expectedLength))

		// Verify the content of each volume
		expectedVolumes := []corev1.Volume{
			{
				Name: "wso2is-pv",
				VolumeSource: corev1.VolumeSource{
					PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
						ClaimName: "wso2is-userstore-pvc",
					},
				},
			},
			{
				Name: "wso2is-config",
				VolumeSource: corev1.VolumeSource{
					ConfigMap: &corev1.ConfigMapVolumeSource{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: "wso2is-config",
						},
					},
				},
			},
			{
				Name: "wso2is-secret",
				VolumeSource: corev1.VolumeSource{
					Secret: &corev1.SecretVolumeSource{
						SecretName: "wso2is-secret",
					},
				},
			},
		}
		Expect(volumes).To(Equal(expectedVolumes))
	})
})

var _ = Describe("Reconcile Volume", func() {
	var (
		r            *Wso2IsReconciler
		instance     wso2v1beta1.Wso2Is
		log          logr.Logger
		ctx          context.Context
		err          error
		reconciled   ctrl.Result
		reconcileErr error
	)
	instance = wso2v1beta1.Wso2Is{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "wso2is",
			Namespace: "wso2-iam-system",
		},
		Spec: wso2v1beta1.Wso2IsSpec{
			Version: "1.0",
		},
	}

	BeforeEach(func() {
		scheme := runtime.NewScheme()
		_ = corev1.AddToScheme(scheme)
		_ = wso2v1beta1.AddToScheme(scheme)

		client := fake.NewClientBuilder().WithScheme(scheme).Build()

		log = logr.Discard()

		r = &Wso2IsReconciler{
			Client: client,
			Log:    log,
			Scheme: scheme,
		}

		ctx = context.TODO()
		err = nil
		reconcileErr = nil
	})

	Context("When PersistentVolumeClaim is not found", func() {
		BeforeEach(func() {
			reconciled, reconcileErr = reconcileVolume(r, instance, log, err, ctx)
		})

		It("Should return without error", func() {
			Expect(reconcileErr).To(HaveOccurred(), "Error during volume reconciliation")
		})

		It("Should not have any reconciliation result", func() {
			Expect(reconciled).To(Equal(ctrl.Result{}))
		})
	})

	Context("When PersistentVolumeClaim is found", func() {
		BeforeEach(func() {
			pvc := &corev1.PersistentVolumeClaim{
				ObjectMeta: metav1.ObjectMeta{
					Name:      instance.Name + "-userstore-pvc",
					Namespace: instance.Namespace,
				},
			}

			r.Client = fake.NewClientBuilder().WithScheme(r.Scheme).WithRuntimeObjects(&instance, pvc).Build()

			reconciled, reconcileErr = reconcileVolume(r, instance, log, err, ctx)
		})

		It("Should return without error", func() {
			Expect(reconcileErr).NotTo(HaveOccurred(), "Error during volume reconciliation")
		})

		It("Should not have any reconciliation result", func() {
			Expect(reconciled).To(Equal(ctrl.Result{}))
		})
	})
})
