package controller

import (
	"context"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// Controller reconciles a custom resource
type Controller struct {
	client.Client
	Scheme *runtime.Scheme
}

// Reconcile handles the reconciliation loop for the custom resource
func (r *Controller) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	log.Info("Reconciling custom resource", "request", req)

	// Add your reconciliation logic here

	// Requeue after 1 minute if no error
	return ctrl.Result{RequeueAfter: time.Minute}, nil
}

// SetupWithManager sets up the controller with the Manager
func (r *Controller) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		// Add your custom resource type here
		// For example: For(&examplev1.YourCustomResource{}).
		Complete(r)
}
