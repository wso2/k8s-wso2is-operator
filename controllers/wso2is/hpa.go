package wso2is

import (
	"context"

	"github.com/go-logr/logr"
	wso2v1beta1 "github.com/wso2/k8s-wso2is-operator/api/v1beta1"
	autoscalingv2 "k8s.io/api/autoscaling/v2"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
)

func (r *Wso2IsReconciler) defineHpa(m wso2v1beta1.Wso2Is) *autoscalingv2.HorizontalPodAutoscaler {
	minReplicas := int32(1)
	maxReplicas := int32(3)
	targetCpuUtilizationPercentage := int32(50)

	hpa := &autoscalingv2.HorizontalPodAutoscaler{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name + "-hpa",
			Namespace: m.Namespace,
		},
		Spec: autoscalingv2.HorizontalPodAutoscalerSpec{
			ScaleTargetRef: autoscalingv2.CrossVersionObjectReference{
				APIVersion: "apps/v1",
				Kind:       "StatefulSet",
				Name:       m.Name,
			},
			MinReplicas: &minReplicas,
			MaxReplicas: maxReplicas,
			Metrics: []autoscalingv2.MetricSpec{
				{
					Type: autoscalingv2.ResourceMetricSourceType,
					Resource: &autoscalingv2.ResourceMetricSource{
						Name: corev1.ResourceCPU,
						Target: autoscalingv2.MetricTarget{
							Type:               autoscalingv2.UtilizationMetricType,
							AverageUtilization: &targetCpuUtilizationPercentage,
						},
					},
				},
			},
		},
	}
	ctrl.SetControllerReference(&m, hpa, r.Scheme)
	return hpa
}
func reconcileHpa(r *Wso2IsReconciler, instance wso2v1beta1.Wso2Is, log logr.Logger, err error, ctx context.Context) (ctrl.Result, error) {
	hpaDefinition := r.defineHpa(instance)
	hpa := &corev1.Service{}
	err = r.Get(ctx, types.NamespacedName{Name: instance.Name + "-hpa", Namespace: instance.Namespace}, hpa)
	if err != nil {
		if errors.IsNotFound(err) {
			log.Info("Hpa resource " + instance.Name + "-hpa" + " not found. Creating or re-creating hpa")
			err = r.Create(ctx, hpaDefinition)
			if err != nil {
				log.Error(err, "Failed to create new Hpa", "Hpa.Namespace", hpaDefinition.Namespace, "Hpa.Name", hpaDefinition.Name)
				return ctrl.Result{}, err
			}
		} else {
			log.Info("Failed to get hpa resource " + instance.Name + "-hpa" + ". Re-running reconcile.")
			return ctrl.Result{Requeue: true}, err
		}
	} else {
		log.Info("Found Hpa")
	}
	return ctrl.Result{}, nil
}
