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
		return r.updateStatus(instance, kranev1alpha1.CanaryPhaseTest, kranev1alpha1.CanaryStateBootstrap)
	}

	// output variables for each phase
	var res reconcile.Result
	// var err error // also used, but declared above
	var state string

	switch instance.Status.Phase {
	case kranev1alpha1.CanaryPhaseTest:
		res, state, err = r.reconcileTestPhase(instance)
	case kranev1alpha1.CanaryPhaseAnalysis:
		res, state, err = r.reconcileAnalysisPhase(instance)
	default:
		log.Printf("Unrecognized instance phase occured for %s/%s: %s", instance.Namespace, instance.Name, instance.Status.Phase)
		return reconcile.Result{}, nil
	}

	// reschedule in case of error within the phase (no in-phase errors are reported here)
	if err != nil {
		return reconcile.Result{}, err
	}

	// update status and reconcile in case of status change
	if instance.Status.State != state {
		if state == kranev1alpha1.CanaryStateSuccess || state == kranev1alpha1.CanaryStateFailed {
			log.Printf("Canary analysis finished for %s/%s, result: %s", instance.Namespace, instance.Name, state)

			// if we just finished the test phase and canary is scheduled, change phases
			if instance.Spec.AnalysisPhase.Image != "" && instance.Status.Phase == kranev1alpha1.CanaryPhaseTest && state == kranev1alpha1.CanaryStateSuccess {
				return r.updateStatus(instance, kranev1alpha1.CanaryPhaseAnalysis, kranev1alpha1.CanaryStateBootstrap)
			}
		}

		return r.updateStatus(instance, instance.Status.Phase, state)
	}

	return res, nil
}

func (r *ReconcileCanary) updateStatus(cr *kranev1alpha1.Canary, phase string, state string) (reconcile.Result, error) {
	instance := cr.DeepCopy()
	instance.Status.Phase = phase
	instance.Status.State = state

	return reconcile.Result{Requeue: true}, r.client.Update(context.TODO(), instance)
}

func (r *ReconcileCanary) reconcileTestPhase(cr *kranev1alpha1.Canary) (reconcile.Result, string, error) {
	// get the current pod or create a new one
	pod, err := r.bootstrapTestPhase(cr)
	if err != nil || pod == nil {
		if err != nil {
			log.Printf("An error occured during TEST bootstrap %s/%s: %s\n", cr.Namespace, cr.Name, err)
		}
		return reconcile.Result{}, kranev1alpha1.CanaryStateBootstrap, err
	}

	// switch based on the pod internal status
	switch pod.Status.Phase {
	case corev1.PodRunning:
		return reconcile.Result{}, kranev1alpha1.CanaryStateInProgress, nil

	case corev1.PodSucceeded:
		return reconcile.Result{}, kranev1alpha1.CanaryStateSuccess, nil

	case corev1.PodFailed:
		return reconcile.Result{}, kranev1alpha1.CanaryStateFailed, nil
	}

	return reconcile.Result{}, cr.Status.Phase, nil
}

func (r *ReconcileCanary) reconcileAnalysisPhase(cr *kranev1alpha1.Canary) (reconcile.Result, string, error) {
	// get the current pod or create a new one
	pod, err := r.bootstrapAnalysisPhase(cr)
	if err != nil || pod == nil {
		if err != nil {
			log.Printf("An error occured during ANALYSIS bootstrap %s/%s: %s\n", cr.Namespace, cr.Name, err)
		}
		return reconcile.Result{}, kranev1alpha1.CanaryStateBootstrap, err
	}

	// switch based on the pod internal status
	switch pod.Status.Phase {
	case corev1.PodRunning:
		return reconcile.Result{}, kranev1alpha1.CanaryStateInProgress, nil

	case corev1.PodSucceeded:
		return reconcile.Result{}, kranev1alpha1.CanaryStateSuccess, nil

	case corev1.PodFailed:
		return reconcile.Result{}, kranev1alpha1.CanaryStateFailed, nil
	}

	return reconcile.Result{}, cr.Status.Phase, nil
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
			Name:      name,
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:    "tester",
					Image:   cr.Spec.TestPhase.Image,
					Command: cr.Spec.TestPhase.Cmd,
					Env: []corev1.EnvVar{
						{
							Name:  "KRANE_TARGET",
							Value: cr.Spec.Target,
						},
					},
				},
			},
			RestartPolicy: corev1.RestartPolicyNever,
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
			Name:      name,
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:    "analyzer",
					Image:   cr.Spec.AnalysisPhase.Image,
					Command: cr.Spec.AnalysisPhase.Cmd,
					Env: []corev1.EnvVar{
						{
							Name:  "KRANE_TARGET",
							Value: cr.Spec.Target,
						},
					},
				},
			},
			RestartPolicy: corev1.RestartPolicyNever,
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
