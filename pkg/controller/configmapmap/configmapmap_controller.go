package configmapmap

import (
	"context"

	operatorsv1alpha1 "github.com/joesonw/configmap-map-operator/pkg/apis/operators/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/workqueue"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_configmapmap")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new ConfigMapMap Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileConfigMapMap{
		client:       mgr.GetClient(),
		scheme:       mgr.GetScheme(),
		watchConfigs: map[string]map[string]types.NamespacedName{},
		watchSecret:  map[string]map[string]types.NamespacedName{},
		oldSpecs:     map[string]map[string]*operatorsv1alpha1.ConfigMapMapSpec{},
	}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("configmapmap-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource ConfigMapMap
	err = c.Watch(&source.Kind{Type: &operatorsv1alpha1.ConfigMapMap{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &corev1.ConfigMap{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &operatorsv1alpha1.ConfigMapMap{},
	})
	if err != nil {
		return err
	}
	reconciler := r.(*ReconcileConfigMapMap)

	err = c.Watch(&source.Kind{Type: &corev1.ConfigMap{}}, &handler.Funcs{
		CreateFunc: nil,
		UpdateFunc: func(event event.UpdateEvent, limitingInterface workqueue.RateLimitingInterface) {
			cm := event.ObjectNew.(*corev1.ConfigMap)
			if configs, ok := reconciler.watchConfigs[cm.Namespace]; ok {
				if ns, ok := configs[cm.Name]; ok {
					limitingInterface.Add(reconcile.Request{NamespacedName: ns})
				}
			}
			if secrets, ok := reconciler.watchSecret[cm.Namespace]; ok {
				if ns, ok := secrets[cm.Name]; ok {
					limitingInterface.Add(reconcile.Request{NamespacedName: ns})
				}
			}
		},
		DeleteFunc:  nil,
		GenericFunc: nil,
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileConfigMapMap implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileConfigMapMap{}

// ReconcileConfigMapMap reconciles a ConfigMapMap object
type ReconcileConfigMapMap struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client       client.Client
	scheme       *runtime.Scheme
	oldSpecs     map[string]map[string]*operatorsv1alpha1.ConfigMapMapSpec
	watchConfigs map[string]map[string]types.NamespacedName
	watchSecret  map[string]map[string]types.NamespacedName
}

// Reconcile reads that state of the cluster for a ConfigMapMap object and makes changes based on the state read
func (r *ReconcileConfigMapMap) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling ConfigMapMap")

	// Fetch the ConfigMapMap instance
	instance := &operatorsv1alpha1.ConfigMapMap{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	if resources, ok := r.oldSpecs[request.Namespace]; ok {
		spec := resources[request.Name]
		if spec != nil {
			for _, item := range spec.Data {
				if item.Kind == "configmap" || item.Kind == "cm" {
					if configs, ok := r.watchConfigs[item.Namespace]; ok {
						delete(configs, item.Name)
					}
				} else if item.Kind == "secret" {
					if secrets, ok := r.watchSecret[item.Namespace]; ok {
						delete(secrets, item.Name)
					}
				}
			}
		}
	}

	spec := instance.Spec
	cm := &corev1.ConfigMap{}
	err = r.client.Get(context.TODO(), client.ObjectKey{Namespace: spec.Namespace, Name: spec.Name}, cm)
	if err != nil && errors.IsNotFound(err) {
		cm = &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: spec.Namespace,
				Name:      spec.Name,
				Labels: map[string]string{
					"configmapmap": instance.Name,
				},
				OwnerReferences: []metav1.OwnerReference{*metav1.NewControllerRef(instance, schema.GroupVersionKind{
					Group:   operatorsv1alpha1.SchemeGroupVersion.Group,
					Version: operatorsv1alpha1.SchemeGroupVersion.Version,
					Kind:    "ConfigMapMap",
				})},
			},
			Data: map[string]string{},
		}
		cm.Data, err = r.createDataMap(instance)
		if err != nil {
			return reconcile.Result{}, err
		}

		if err := r.client.Create(context.TODO(), cm); err != nil {
			return reconcile.Result{}, err
		}
	} else if err == nil {
		cm.Data, err = r.createDataMap(instance)
		if err != nil {
			return reconcile.Result{}, err
		}

		if err := r.client.Update(context.TODO(), cm); err != nil {
			return reconcile.Result{}, err
		}
	}

	if err == nil || errors.IsNotFound(err) {
		if _, ok := r.oldSpecs[request.Namespace]; !ok {
			r.oldSpecs[request.Namespace] = map[string]*operatorsv1alpha1.ConfigMapMapSpec{}
		}
		r.oldSpecs[request.Namespace][request.Name] = spec.DeepCopy()
	}

	return reconcile.Result{}, err
}

func (r *ReconcileConfigMapMap) createDataMap(instance *operatorsv1alpha1.ConfigMapMap) (map[string]string, error) {
	spec := instance.Spec
	data := map[string]string{}
	for key, item := range spec.Data {
		if item.Kind == "cm" || item.Kind == "configmap" {
			config := &corev1.ConfigMap{}
			if err := r.client.Get(context.TODO(), client.ObjectKey{Namespace: item.Namespace, Name: item.Name}, config); err != nil {
				continue
			}
			value, ok := config.Data[item.SubPath]
			if !ok {
				continue
			}
			data[key] = value
			if _, ok := r.watchConfigs[item.Namespace]; !ok {
				r.watchConfigs[item.Namespace] = map[string]types.NamespacedName{}
			}
			r.watchConfigs[item.Namespace][item.Name] = getNamespaceName(instance)
		} else if item.Kind == "secret" {
			secret := &corev1.Secret{}
			if err := r.client.Get(context.TODO(), client.ObjectKey{Namespace: item.Namespace, Name: item.Name}, secret); err != nil {
				continue
			}
			bytes, ok := secret.Data[item.SubPath]
			if !ok {
				continue
			}
			data[key] = string(bytes)
			if _, ok := r.watchSecret[item.Namespace]; !ok {
				r.watchSecret[item.Namespace] = map[string]types.NamespacedName{}
			}
			r.watchSecret[item.Namespace][item.Name] = getNamespaceName(instance)
		}
	}

	return data, nil
}

func getNamespaceName(obj metav1.Object) types.NamespacedName {
	return types.NamespacedName{
		Namespace: obj.GetNamespace(),
		Name:      obj.GetName(),
	}
}
