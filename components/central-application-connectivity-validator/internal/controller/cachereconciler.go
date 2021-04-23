package controller

import (
	"context"
	"github.com/kyma-project/kyma/common/logging/logger"
	v1alpha1 "github.com/kyma-project/kyma/components/application-operator/pkg/apis/applicationconnector/v1alpha1"
	gocache "github.com/patrickmn/go-cache"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type CacheReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	AppCache *gocache.Cache
	Log      *logger.Logger
}

func (r *CacheReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	//TODO: ??? will delete remove from all possible instances ?

	ctx := context.Background()
	var instance v1alpha1.Application
	if err := r.Get(ctx, req.NamespacedName, &instance); err != nil {
		err = client.IgnoreNotFound(err)
		if err != nil {

		} else {

		}
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

func (r *CacheReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.Application{}).
		Complete(r)
}
