package wso2is

import (
	"context"

	"github.com/go-logr/logr"
	wso2v1beta1 "github.com/wso2/k8s-wso2is-operator/api/v1beta1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
)

func reconcileIngress(r *Wso2IsReconciler, instance wso2v1beta1.Wso2Is, log logr.Logger, err error, ctx context.Context) (ctrl.Result, error) {
	// TODO:

	// Check for ingress
	ingressFound := networkingv1.Ingress{}
	//logger.Info("ingress name: " + ingName + "namespace: " + instance.Namespace)
	err = r.Get(ctx, types.NamespacedName{Name: instance.Name + "-ingress", Namespace: instance.Namespace}, &ingressFound)
	if err != nil && errors.IsNotFound(err) {
		log.Info("Unable to detect Ingress in your cluster. You may configure your own")
		return ctrl.Result{}, nil
	} else if err != nil {
		log.Error(err, "Failed to get Ingress")
		return ctrl.Result{}, err
	} else {
		log.Info("Found Ingress")
	}

	// Update ingress details in status
	if len(ingressFound.Status.LoadBalancer.Ingress) > 0 {
		instance.Status.IngressHostname = ingressFound.Status.LoadBalancer.Ingress[0].Hostname
	}

	return ctrl.Result{}, nil
}
