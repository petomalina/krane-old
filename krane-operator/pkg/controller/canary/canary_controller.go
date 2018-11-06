package canary

import (
	"context"
	"log"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	kranev1alpha1 "github.com/petomalina/krane/krane-operator/pkg/apis/krane/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new Canary Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileCanary{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("canary-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource Canary
	err = c.Watch(&source.Kind{Type: &kranev1alpha1.Canary{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// Watch our testing and analysis pods
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &kranev1alpha1.Canary{},
	})
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileCanary{}

// ReconcileCanary reconciles a Canary object
type ReconcileCanary struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a Canary object and makes changes based on the state read
// and what is in the Canary.Spec
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileCanary) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	log.Printf("Reconciling Canary %s/%s\n", request.Namespace, request.Name)

	instance := &kranev1alpha1.Canary{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}

		return reconcile.Result{}, err
	}

	// initialize with the default test bootstrap phase and reconcile
	if instance.Status.Phase == "" {
		log.Printf("Bootsraping newly created Canary %s/%s\n", instance.Namespace, instance.Name)

		instance = instance.DeepCopy()
		instance.Status.Phase = kranev1alpha1.CanaryPhaseTest
		instance.Status.State = kranev1alpha1.CanaryStatusBootstrap

		return reconcile.Result{Requeue: true}, r.client.Update(context.TODO(), instance)
	}

	switch instance.Status.Phase {
	case kranev1alpha1.CanaryPhaseTest:
		return r.reconcileTestPhase(instance)
	case kranev1alpha1.CanaryPhaseAnalysis:
		return r.reconcileAnalysisPhase(instance)
	default:
		log.Printf("Unrecognized instance phase occured for %s/%s: %s", instance.Namespace, instance.Name, instance.Status.Phase)
		return reconcile.Result{}, nil
	}
}

func (r *ReconcileCanary) reconcileTestPhase(cr *kranev1alpha1.Canary) (reconcile.Result, error) {
	// get the current pod or create a new one
	pod, err := r.bootstrapTestPhase(cr)
	if err != nil || pod == nil {
		if err != nil {
			log.Printf("An error occured during test bootstrap %s/%s: %s\n", cr.Namespace, cr.Name, err)
		}
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}

func (r *ReconcileCanary) reconcileAnalysisPhase(cr *kranev1alpha1.Canary) (reconcile.Result, error) {
	// get the current pod or create a new one
	pod, err := r.bootstrapAnalysisPhase(cr)
	if err != nil || pod == nil {
		if err != nil {
			log.Printf("An error occured during analysis bootstrap %s/%s: %s\n", cr.Namespace, cr.Name, err)
		}
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}

func (r *ReconcileCanary) bootstrapTestPhase(cr *kranev1alpha1.Canary) (*corev1.Pod, error) {
	name := cr.Name + "-phase-test"

	found := &corev1.Pod{}
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: name, Namespace: cr.Namespace}, found)
	if err != nil && !errors.IsNotFound(err) {
		return nil, err
	} else if err == nil {
		return found, nil
	}

	log.Printf("Testing pod not found for %s/%s, creating\n", cr.Namespace, cr.Name)

	labels := map[string]string{
		"app": cr.Name,
	}
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name + "-pod",
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:    "busybox",
					Image:   "busybox",
					Command: []string{"sleep", "10"},
				},
			},
		},
	}

	if err := controllerutil.SetControllerReference(cr, pod, r.scheme); err != nil {
		return nil, err
	}

	err = r.client.Create(context.TODO(), pod)
	if err != nil {
		return nil, err
	}

	return pod, nil
}

func (r *ReconcileCanary) bootstrapAnalysisPhase(cr *kranev1alpha1.Canary) (*corev1.Pod, error) {
	name := cr.Name + "-phase-analysis"

	found := &corev1.Pod{}
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: name, Namespace: cr.Namespace}, found)
	if err != nil && !errors.IsNotFound(err) {
		return nil, err
	} else if err == nil {
		return found, nil
	}

	log.Printf("Testing pod not found for %s/%s, creating\n", cr.Namespace, cr.Name)

	labels := map[string]string{
		"app": cr.Name,
	}
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name + "-pod",
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:    "busybox",
					Image:   "busybox",
					Command: []string{"sleep", "10"},
				},
			},
		},
	}

	if err := controllerutil.SetControllerReference(cr, pod, r.scheme); err != nil {
		return nil, err
	}

	err = r.client.Create(context.TODO(), pod)
	if err != nil {
		return nil, err
	}

	return pod, nil
}
