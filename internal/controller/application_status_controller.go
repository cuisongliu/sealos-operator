package controller

import (
	"context"
	appsv1beta1 "github.com/cuisongliu/sealos-operator/api/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/retry"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func (r *ApplicationReconciler) statusReconcile(ctx context.Context, obj client.Object) error {
	log := log.FromContext(ctx)
	log.V(1).Info("update reconcile status controller application", "request", client.ObjectKeyFromObject(obj))
	application := &appsv1beta1.Application{}
	err := r.Get(ctx, client.ObjectKeyFromObject(obj), application)
	if err != nil {
		return client.IgnoreNotFound(err)
	}
	application.Status.Phase = appsv1beta1.ApplicationPending
	// Let's just set the status as Unknown when no status are available
	status := true
	for _, v := range application.Status.Conditions {
		if v.Status != metav1.ConditionTrue {
			status = false
			break
		}
	}
	if !status {
		application.Status.Phase = appsv1beta1.ApplicationError
	} else {
		application.Status.Phase = appsv1beta1.ApplicationInProcess
	}
	if application.Status.Phase == appsv1beta1.ApplicationInProcess {
		//TODO add logic here
	}

	return r.syncStatus(ctx, application)
}

func (r *ApplicationReconciler) syncStatus(ctx context.Context, application *appsv1beta1.Application) error {
	log := log.FromContext(ctx)
	// Let's just set the status as Unknown when no status are available
	if err := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		original := &appsv1beta1.Application{}
		if err := r.Get(ctx, client.ObjectKeyFromObject(application), original); err != nil {
			return err
		}
		original.Status = *application.Status.DeepCopy()
		return r.Client.Status().Update(ctx, original)
	}); err != nil {
		log.Error(err, "Failed to update application status")
		return err
	}
	return nil
}
