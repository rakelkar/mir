/*
Copyright 2019 Microsoft Corporation.

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

package modeldeploymentsource

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"

	"k8s.io/apimachinery/pkg/runtime/serializer"

	"net/http"

	"github.com/rakelkar/mir/pkg/apis/mir/v1beta1"
	mirv1beta1 "github.com/rakelkar/mir/pkg/apis/mir/v1beta1"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new ModelDeploymentSource Controller and adds it to the Manager with default RBAC. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileModelDeploymentSource{Client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("modeldeploymentsource-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to ModelDeploymentSource
	err = c.Watch(&source.Kind{Type: &mirv1beta1.ModelDeploymentSource{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create
	// Uncomment watch a Deployment created by ModelDeploymentSource - change this for objects you create
	err = c.Watch(&source.Kind{Type: &appsv1.Deployment{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &mirv1beta1.ModelDeploymentSource{},
	})
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileModelDeploymentSource{}

// ReconcileModelDeploymentSource reconciles a ModelDeploymentSource object
type ReconcileModelDeploymentSource struct {
	client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a ModelDeploymentSource object and makes changes based on the state read
// and what is in the ModelDeploymentSource.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  The scaffolding writes
// a Deployment as an example
// Automatically generate RBAC rules to allow the Controller to read and write Deployments
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=deployments/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=mir.k8s.io,resources=modeldeploymentsources,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=mir.k8s.io,resources=modeldeploymentsources/status,verbs=get;update;patch
func (r *ReconcileModelDeploymentSource) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	mir_service_host := os.Getenv("MIR_SERVICE_HOST")
	if len(mir_service_host) == 0 {
		mir_service_host = "MIR_SERVICE_HOST"
	}

	// Fetch the ModelDeploymentSource instance
	instance := &mirv1beta1.ModelDeploymentSource{}
	err := r.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Object not found, return.  Created objects are automatically garbage collected.
			// For additional cleanup logic use finalizers.
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	// TODO: stuff that needs to come labels set by mutating admission controller
	var mir_dns_prefix, mir string
	containingNs := &v1.Namespace{}
	err = r.Get(context.TODO(), types.NamespacedName{Name: instance.Namespace, Namespace: ""}, containingNs)
	if err != nil {
		return reconcile.Result{}, err
	}
	ok := false
	if mir, ok = containingNs.ObjectMeta.Labels["mir"]; !ok {
		return reconcile.Result{}, fmt.Errorf("namespace is missing mir label")
	}
	if mir_dns_prefix, ok = containingNs.ObjectMeta.Labels["mir-dns-prefix"]; !ok {
		return reconcile.Result{}, fmt.Errorf("namespace is missing mir-dns-prefix label")
	}

	sourceNsName := mir + "-" + instance.Name
	stuffNsName := sourceNsName + "-stuff"

	// Create a namespace for CRDs
	if result, err := r.createNamespace(instance, sourceNsName, sourceNsName, stuffNsName, mir, mir_dns_prefix); err != nil {
		return result, err
	}

	// Create a namespace for other stuff
	if result, err := r.createNamespace(instance, stuffNsName, sourceNsName, stuffNsName, mir, mir_dns_prefix); err != nil {
		return result, err
	}

	url := fmt.Sprintf("http://%s/mir/%s/modelsource/%s/models", mir_service_host, mir, instance.Name)
	msl, err := r.downloadModelList(instance, url)
	if err != nil {
		return reconcile.Result{}, err
	}

	if err := r.writeModelsToNamespace(instance, sourceNsName, msl); err != nil {
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}

func (r *ReconcileModelDeploymentSource) createNamespace(instance *mirv1beta1.ModelDeploymentSource, nsName string, sourceNsName string, stuffNsName string, mir string, mir_dns_prefix string) (reconcile.Result, error) {
	ns := &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: nsName,
			Labels: map[string]string{
				"mir":            mir,
				"mir-dns-prefix": mir_dns_prefix,
				"source-ns":      sourceNsName,
				"model-ns":       stuffNsName,
			},
		},
		Spec: v1.NamespaceSpec{},
	}

	if err := controllerutil.SetControllerReference(instance, ns, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	// Check if the Namespace already exists
	found := &v1.Namespace{}
	err := r.Get(context.TODO(), types.NamespacedName{Name: ns.Name, Namespace: ""}, found)
	if err != nil && errors.IsNotFound(err) {
		log.Info("Creating Namespace", "namespace", "default", "name", ns.Name)
		err = r.Create(context.TODO(), ns)
	}
	if err != nil {
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}

func (r *ReconcileModelDeploymentSource) downloadModelList(instance *mirv1beta1.ModelDeploymentSource, url string) (*mirv1beta1.ModelServiceList, error) {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		log.Info("Failed to download models for model source", "source", instance.Name, "url", url, "error", err.Error())
		return &v1beta1.ModelServiceList{}, nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		b, _ := ioutil.ReadAll(resp.Body)
		log.Info("Failed to download models for model source", "source", instance.Name, "url", url, "error", string(b))
		return &v1beta1.ModelServiceList{}, nil
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// attempt to pull models into the new namespace
	codecs := serializer.NewCodecFactory(r.scheme)
	deserializer := codecs.UniversalDeserializer()
	obj, groupVersionKind, err := deserializer.Decode(b, nil, nil)
	if err != nil {
		return nil, err
	}

	if groupVersionKind.Group != "mir" || groupVersionKind.Kind != "ModelService" {
		return nil, fmt.Errorf("url %s returned invalid type: %s", url, groupVersionKind)
	}

	msl := obj.(*mirv1beta1.ModelServiceList)
	return msl, nil
}

func (r *ReconcileModelDeploymentSource) writeModelsToNamespace(instance *mirv1beta1.ModelDeploymentSource, namespaceName string, msl *mirv1beta1.ModelServiceList) error {
	for _, ms := range msl.Items {
		ms.SetNamespace(namespaceName)

		if err := controllerutil.SetControllerReference(instance, &ms, r.scheme); err != nil {
			return err
		}

		// create if not found
		found := &mirv1beta1.ModelService{}
		err := r.Get(context.TODO(), types.NamespacedName{Name: ms.Name, Namespace: ms.Namespace}, found)
		if err != nil && errors.IsNotFound(err) {
			log.Info("Creating Model Service", "namespace", ms.Namespace, "name", ms.Name)
			err = r.Create(context.TODO(), &ms)
		}
		if err != nil {
			return err
		}
	}

	return nil
}
