package controllers

import (
	"context"
	"github.com/go-logr/logr"
	wso2v1beta1 "github.com/wso2/k8s-wso2is-operator/api/v1beta1"
	ctrl "sigs.k8s.io/controller-runtime"
)

func reconcileStatefulSet(r *Wso2IsReconciler, instance wso2v1beta1.Wso2Is, log logr.Logger, err error, ctx context.Context) (ctrl.Result, error) {
	// Define a new StatefulSet
	statefulSet := r.statefulSetForWso2Is(instance)
	log.Info("Creating a new StatefulSet", "StatefulSet.Namespace", statefulSet.Namespace, "StatefulSet.Name", statefulSet.Name)

	err = r.Create(ctx, statefulSet)
	if err != nil {
		log.Error(err, "Failed to create new StatefulSet", "StatefulSet.Namespace", statefulSet.Namespace, "StatefulSet.Name", statefulSet.Name)
		return ctrl.Result{}, err
	} else {
		log.Info("Successfully added new StatefulSet", "StatefulSet.Namespace", statefulSet.Namespace, "StatefulSet.Name", statefulSet.Name)
	}
	// StatefulSet created successfully - return and requeue
	return ctrl.Result{Requeue: true}, nil
}

func reconcileDeployment(r *Wso2IsReconciler, instance wso2v1beta1.Wso2Is, log logr.Logger, err error, ctx context.Context) (ctrl.Result, error) {
	// Define a new deployment
	dep := r.deploymentForWso2Is(instance)
	log.Info("Creating a new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)

	err = r.Create(ctx, dep)
	if err != nil {
		log.Error(err, "Failed to create new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
		return ctrl.Result{}, err
	} else {
		log.Info("Successfully added new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
	}
	// Deployment created successfully - return and requeue
	return ctrl.Result{Requeue: true}, nil
}

func reconcileSvc(r *Wso2IsReconciler, instance wso2v1beta1.Wso2Is, log logr.Logger, err error, ctx context.Context) (ctrl.Result, error) {
	// Define a new Service
	svc := r.addNewService(instance)
	log.Info("Creating a new Service", "Service.Namespace", svc.Namespace, "Service.Name", svc.Name)
	err = r.Create(ctx, svc)
	if err != nil {
		log.Error(err, "Failed to create new Service", "Service.Namespace", svc.Namespace, "Service.Name", svc.Name)
		return ctrl.Result{}, err
	} else {
		log.Info("Successfully created new Service", "Service.Namespace", svc.Namespace, "Service.Name", svc.Name)
	}
	// Service created successfully - return and requeue
	return ctrl.Result{Requeue: true}, nil
}

func reconcileCfg(r *Wso2IsReconciler, instance wso2v1beta1.Wso2Is, log logr.Logger, err error, ctx context.Context) (ctrl.Result, error) {
	// Define a new ConfigMap
	cfgMap := r.addConfigMap(instance, log)
	log.Info("Creating a new ConfigMap", "ConfigMap.Namespace", cfgMap.Namespace, "ConfigMap.Name", cfgMap.Name)
	err = r.Create(ctx, cfgMap)
	if err != nil {
		log.Error(err, "Failed to create new ConfigMap", "ConfigMap.Namespace", cfgMap.Namespace, "ConfigMap.Name", cfgMap.Name)
		return ctrl.Result{}, err
	} else {
		log.Info("Successfully created new ConfigMap", "ConfigMap.Namespace", cfgMap.Namespace, "ConfigMap.Name", cfgMap.Name)
	}
	// ConfigMap created successfully - return and requeue
	return ctrl.Result{Requeue: true}, nil
}

func reconcileSecret(r *Wso2IsReconciler, instance wso2v1beta1.Wso2Is, log logr.Logger, err error, ctx context.Context) (ctrl.Result, error) {
	// Define a new secret
	secret := r.addSecret(instance, log)
	log.Info("Creating a new Secret", "Secret.Namespace", secret.Namespace, "Secret.Name", secret.Name)
	err = r.Create(ctx, secret)
	if err != nil {
		log.Error(err, "Failed to create new Secret", "Secret.Namespace", secret.Namespace, "Secret.Name", secret.Name)
		return ctrl.Result{}, err
	} else {
		log.Info("Successfully created new Secret", "Secret.Namespace", secret.Namespace, "Secret.Name", secret.Name)
	}
	// Secret created successfully - return and requeue
	return ctrl.Result{Requeue: true}, nil
}

func reconcileSva(r *Wso2IsReconciler, instance wso2v1beta1.Wso2Is, log logr.Logger, err error, ctx context.Context) (ctrl.Result, error) {
	// Define a new ServiceAccount
	svc := r.addServiceAccount(instance)
	log.Info("Creating a new ServiceAccount", "ServiceAccount.Namespace", svc.Namespace, "ServiceAccount.Name", svc.Name)
	err = r.Create(ctx, svc)
	if err != nil {
		log.Error(err, "Failed to create new ServiceAccount", "ServiceAccount.Namespace", svc.Namespace, "ServiceAccount.Name", svc.Name)
		return ctrl.Result{}, err
	} else {
		log.Info("Successfully created new ServiceAccount", "ServiceAccount.Namespace", svc.Namespace, "ServiceAccount.Name", svc.Name)
	}
	// ServiceAccount created successfully - return and requeue
	return ctrl.Result{Requeue: true}, nil
}

func reconcileRole(r *Wso2IsReconciler, instance wso2v1beta1.Wso2Is, log logr.Logger, err error, ctx context.Context) (ctrl.Result, error) {
	// Define a new Role
	role := r.addRole(instance)
	log.Info("Creating a new Role", "Role.Namespace", role.Namespace, "Role.Name", role.Name)
	err = r.Create(ctx, role)
	if err != nil {
		log.Error(err, "Failed to create new Role", "Role.Namespace", role.Namespace, "Role.Name", role.Name)
		return ctrl.Result{}, err
	} else {
		log.Info("Successfully created new Role", "Role.Namespace", role.Namespace, "Role.Name", role.Name)
	}
	// Role created successfully - return and requeue
	return ctrl.Result{Requeue: true}, nil
}

func reconcileRoleBinding(r *Wso2IsReconciler, instance wso2v1beta1.Wso2Is, log logr.Logger, err error, ctx context.Context) (ctrl.Result, error) {
	// Define a new RoleBinding
	roleBinding := r.addRoleBinding(instance)
	log.Info("Creating a new RoleBinding", "RoleBinding.Namespace", roleBinding.Namespace, "RoleBinding.Name", roleBinding.Name)
	err = r.Create(ctx, roleBinding)
	if err != nil {
		log.Error(err, "Failed to create new RoleBinding", "RoleBinding.Namespace", roleBinding.Namespace, "RoleBinding.Name", roleBinding.Name)
		return ctrl.Result{}, err
	} else {
		log.Info("Successfully created new RoleBinding", "RoleBinding.Namespace", roleBinding.Namespace, "RoleBinding.Name", roleBinding.Name)
	}
	// RoleBinding created successfully - return and requeue
	return ctrl.Result{Requeue: true}, nil
}
