/*


Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package wso2is

import (
	"context"
	wso2v1beta1 "github.com/wso2/k8s-wso2is-operator/api/v1beta1"
	"github.com/wso2/k8s-wso2is-operator/variables"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

// +kubebuilder:rbac:groups=iam.wso2.com,resources=wso2is,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=iam.wso2.com,resources=wso2is/status,verbs=get;update;patch

func (r *Wso2IsReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	// Get logger
	logger := r.Log.WithValues(req.Name, req.NamespacedName)
	logger.Info("\n-----------------------\nTriggered Reconcile Method\n-----------------------\n")
	//logger.Info(req.NamespacedName.Name)

	// Fetch the WSO2IS instance
	instance := wso2v1beta1.Wso2Is{}
	//err := r.Get(ctx, req.NamespacedName, &instance)
	// TODO: literal strings remove.
	err := r.Get(ctx, types.NamespacedName{Name: "wso2is", Namespace: variables.Wso2IsNamespace}, &instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			logger.Info("WSO2IS resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		logger.Error(err, "Failed to get WSO2IS Instance")
		return ctrl.Result{}, err
	}

	_, err = reconcileSva(r, instance, logger, err, ctx)
	if err != nil {
		return ctrl.Result{}, err
	}

	_, err = reconcileSvc(r, instance, logger, err, ctx)
	if err != nil {
		return ctrl.Result{}, err
	}

	_, err = reconcileVolume(r, instance, logger, err, ctx)
	if err != nil {
		return ctrl.Result{}, err
	}

	_, err = reconcileConfigMap(r, instance, logger, err, ctx)
	if err != nil {
		return ctrl.Result{}, err
	}

	_, err = reconcileRole(r, instance, logger, err, ctx)
	if err != nil {
		return ctrl.Result{}, err
	}

	_, err = reconcileRoleBinding(r, instance, logger, err, ctx)
	if err != nil {
		return ctrl.Result{}, err
	}

	_, err = reconcileSecret(r, instance, logger, err, ctx)
	if err != nil {
		return ctrl.Result{}, err
	}

	// TODO: not sure of following.

	_, err = reconcileIngress(r, instance, logger, err, ctx)
	if err != nil {
		return ctrl.Result{}, err
	}

	_, err = reconcileStatefulSet(r, instance, logger, err, ctx)
	if err != nil {
		return ctrl.Result{}, err
	}

	//updateStatus(r, instance, logger, err, ctx)

	return ctrl.Result{}, nil
}

func (r *Wso2IsReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&wso2v1beta1.Wso2Is{}).
		Owns(&appsv1.StatefulSet{}).
		Watches(
			&source.Kind{Type: &corev1.ConfigMap{}},
			//&handler.EnqueueRequestForObject{},
			handler.EnqueueRequestsFromMapFunc(func(a client.Object) []reconcile.Request {
				//if a.GetNamespace() == variables.Wso2IsNamespace && a.GetName() == "wso2is-configmap" {
				if a.GetNamespace() == variables.Wso2IsNamespace {
					return []reconcile.Request{
						{
							NamespacedName: types.NamespacedName{
								Name:      a.GetName(),
								Namespace: a.GetNamespace(),
							},
						},
					}
				}
				return []reconcile.Request{}
			}),
		).
		Complete(r)
}
