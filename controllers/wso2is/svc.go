package wso2is

import (
	"context"

	"github.com/go-logr/logr"
	wso2v1beta1 "github.com/wso2/k8s-wso2is-operator/api/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
)

func (r *Wso2IsReconciler) defineService(m wso2v1beta1.Wso2Is) *corev1.Service {
	servicePortHttp := int32(9763)
	servicePortHttps := int32(9443)

	ls := labelsForWso2IS(m.Name, m.Spec.Version)

	// Make Service type configurable
	//serviceType := corev1.ServiceTypeNodePort
	serviceType := corev1.ServiceTypeClusterIP
	if m.Spec.Configurations.ServiceType == "NodePort" {
		serviceType = corev1.ServiceTypeNodePort
	} else if m.Spec.Configurations.ServiceType == "ClusterIP" {
		serviceType = corev1.ServiceTypeClusterIP
	} else if m.Spec.Configurations.ServiceType == "LoadBalancer" {
		serviceType = corev1.ServiceTypeLoadBalancer
	}

	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name + "-service",
			Namespace: m.Namespace,
			Labels:    ls,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{{
				Name:     "servlet-http",
				Protocol: "TCP",
				Port:     servicePortHttp,
				TargetPort: intstr.IntOrString{
					IntVal: servicePortHttp,
				},
			}, {
				Name:     "servlet-https",
				Protocol: "TCP",
				Port:     servicePortHttps,
				TargetPort: intstr.IntOrString{
					IntVal: servicePortHttps,
				},
			}},
			Selector: labelsForWso2IS(m.Name, m.Spec.Version),
			Type:     serviceType,
		},
	}
	ctrl.SetControllerReference(&m, svc, r.Scheme)
	return svc
}

func reconcileSvc(r *Wso2IsReconciler, instance wso2v1beta1.Wso2Is, log logr.Logger, err error, ctx context.Context) (ctrl.Result, error) {
	svcDefinition := r.defineService(instance)
	svc := &corev1.Service{}
	err = r.Get(ctx, types.NamespacedName{Name: instance.Name + "-service", Namespace: instance.Namespace}, svc)
	if err != nil {
		if errors.IsNotFound(err) {
			log.Info("Service resource " + instance.Name + "-service" + " not found. Creating or re-creating service")
			err = r.Create(ctx, svcDefinition)
			if err != nil {
				log.Error(err, "Failed to create new Service", "Service.Namespace", svcDefinition.Namespace, "Service.Name", svcDefinition.Name)
				return ctrl.Result{}, err
			}
		} else {
			log.Info("Failed to get service resource " + instance.Name + "-service" + ". Re-running reconcile.")
			return ctrl.Result{}, err
		}
	} else {
		// Note: For simplication purposes Services are not updated - see deployment section
		log.Info("Found Service")
	}
	return ctrl.Result{}, nil
}
