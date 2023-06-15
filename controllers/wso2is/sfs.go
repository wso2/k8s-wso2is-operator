package wso2is

import (
	"context"
	"fmt"
	"github.com/go-logr/logr"
	wso2v1beta1 "github.com/wso2/k8s-wso2is-operator/api/v1beta1"
	"github.com/wso2/k8s-wso2is-operator/variables"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
)

func (r *Wso2IsReconciler) defineStatefulSet(m wso2v1beta1.Wso2Is) *appsv1.StatefulSet {
	ls := labelsForWso2IS(m.Name, m.Spec.Version)
	replicas := m.Spec.Size
	runasuser := int64(802)

	sfs := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name,
			Namespace: m.Namespace,
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas:    &replicas,
			ServiceName: variables.ServiceName,
			Selector: &metav1.LabelSelector{
				MatchLabels: ls,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      ls,
					Annotations: m.Spec.Template.Annotations,
				},
				Spec: corev1.PodSpec{
					Volumes: MakeVolumes(),
					Containers: []corev1.Container{{
						Name:            variables.DeploymentName,
						Image:           variables.ContainerImage,
						Ports:           MakeContainerPorts(),
						Env:             MakeEnvVars(),
						Resources:       MakeResourceRequirements(),
						VolumeMounts:    MakeVolumeMounts(m.Spec.Version),
						StartupProbe:    MakeStartupProbe(),
						LivenessProbe:   MakeLivenessProbe(),
						ReadinessProbe:  MakeReadinessProbe(),
						Lifecycle:       MakeContainerLifecycle(),
						ImagePullPolicy: "IfNotPresent",
						SecurityContext: &corev1.SecurityContext{
							RunAsUser: &runasuser,
						},
					}},
					ServiceAccountName: variables.ServiceAccountName,
					//HostAliases: []corev1.HostAlias{{
					//	IP:        "127.0.0.1",
					//	Hostnames: []string{m.Spec.Configurations.Host},
					//}},
				},
			},
			MinReadySeconds: 30,
		},
	}
	// Set WSO2IS instance as the owner and controller
	ctrl.SetControllerReference(&m, sfs, r.Scheme)
	return sfs
}

func reconcileStatefulSet(r *Wso2IsReconciler, instance wso2v1beta1.Wso2Is, log logr.Logger, err error, ctx context.Context) (ctrl.Result, error) {
	sfsDefinition := r.defineStatefulSet(instance)
	sfs := &appsv1.StatefulSet{}
	err = r.Get(ctx, types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, sfs)
	if err != nil {
		if errors.IsNotFound(err) {
			log.Info("StatefulSet resource " + instance.Name + " not found. Creating or re-creating statefulset")
			err = r.Create(ctx, sfsDefinition)
			if err != nil {
				log.Error(err, "Failed to create new StatefulSet", "StatefulSet.Namespace", sfsDefinition.Namespace, "StatefulSet.Name", sfsDefinition.Name)
				return ctrl.Result{}, err
			}
		} else {
			log.Info("Failed to get statefulset resource " + instance.Name + ". Re-running reconcile.")
			return ctrl.Result{}, err
		}
	} else {
		log.Info("Found StatefulSet")

		// Update replica status
		instance.Status.Replicas = fmt.Sprint(*sfs.Spec.Replicas)

		// Ensure the StatefulSet size is the same as the spec
		size := instance.Spec.Size

		if *sfs.Spec.Replicas != size {
			sfs.Spec.Replicas = &size
			err = r.Update(ctx, sfs)
			if err != nil {
				log.Error(err, "Failed to update StatefulSet", "StatefulSet.Namespace", sfs.Namespace, "StatefulSet.Name", sfs.Name)
				return ctrl.Result{}, err
			}
			// Spec updated - return and requeue
			return ctrl.Result{Requeue: true}, nil
		}
	}
	return ctrl.Result{}, nil
}
