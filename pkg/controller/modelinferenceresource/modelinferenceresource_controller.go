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

package modelinferenceresource

import (
	"context"
	"reflect"

	mirv1beta1 "github.com/rakelkar/mir/pkg/apis/mir/v1beta1"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
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

// Add creates a new ModelInferenceResource Controller and adds it to the Manager with default RBAC. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileModelInferenceResource{Client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

func use_az_secret() []v1.EnvVar {
	secret_name := "azcreds"

	return []v1.EnvVar{
		v1.EnvVar{
			Name: "AZ_CLIENT_ID",
			ValueFrom: &v1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: secret_name,
					},
					Key: "AZ_CLIENT_ID",
				},
			},
		},
		v1.EnvVar{
			Name: "AZ_CLIENT_SECRET",
			ValueFrom: &v1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: secret_name,
					},
					Key: "AZ_CLIENT_SECRET",
				},
			},
		},
		v1.EnvVar{
			Name: "AZ_TENANT_ID",
			ValueFrom: &v1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: secret_name,
					},
					Key: "AZ_TENANT_ID",
				},
			},
		},
		v1.EnvVar{
			Name: "AZ_SUBSCRIPTION_ID",
			ValueFrom: &v1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: secret_name,
					},
					Key: "AZ_SUBSCRIPTION_ID",
				},
			},
		},
		v1.EnvVar{
			Name: "AZ_SP_OBJECT_ID",
			ValueFrom: &v1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: secret_name,
					},
					Key: "AZ_SP_OBJECT_ID",
				},
			},
		},
	}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("modelinferenceresource-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to ModelInferenceResource
	err = c.Watch(&source.Kind{Type: &mirv1beta1.ModelInferenceResource{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create
	// Uncomment watch a Deployment created by ModelInferenceResource - change this for objects you create
	err = c.Watch(&source.Kind{Type: &appsv1.Deployment{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &mirv1beta1.ModelInferenceResource{},
	})
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileModelInferenceResource{}

// ReconcileModelInferenceResource reconciles a ModelInferenceResource object
type ReconcileModelInferenceResource struct {
	client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a ModelInferenceResource object and makes changes based on the state read
// and what is in the ModelInferenceResource.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  The scaffolding writes
// a Deployment as an example
// Automatically generate RBAC rules to allow the Controller to read and write Deployments
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=deployments/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=mir.k8s.io,resources=modelinferenceresources,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=mir.k8s.io,resources=modelinferenceresources/status,verbs=get;update;patch
func (r *ReconcileModelInferenceResource) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	// Fetch the ModelInferenceResource instance
	instance := &mirv1beta1.ModelInferenceResource{}
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

	// TODO(user): Change this to be the object type created by your controller
	ns := &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:   instance.Name + "-ingress",
			Labels: map[string]string{"mir": instance.Name},
		},
		Spec: v1.NamespaceSpec{},
	}

	if err := controllerutil.SetControllerReference(instance, ns, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	// TODO(user): Change this for the object type created by your controller
	// Check if the Deployment already exists
	{
		found := &v1.Namespace{}
		err = r.Get(context.TODO(), types.NamespacedName{Name: ns.Name, Namespace: "default"}, found)
		if err != nil && errors.IsNotFound(err) {
			log.Info("Creating Namespace", "namespace", "default", "name", ns.Name)
			err = r.Create(context.TODO(), ns)
		} else if err != nil {
			return reconcile.Result{}, err
		}
	}

	az_vars := use_az_secret()

	job_ns := instance.Namespace
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.Name + "-job",
			Namespace: job_ns,
			Labels:    make(map[string]string),
		},
		Spec: batchv1.JobSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: make(map[string]string),
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:   instance.Name + "-job",
					Labels: make(map[string]string),
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:    "azcmd",
							Image:   "kcorer/azcmd",
							Command: []string{"/entrypoint.sh", "/tm.sh", "rakelkar-delete-me", "eastus", "ugh-tm", "ahahahahahaha", "anep", "40.70.209.164"},
							Env:     az_vars,
						},
					},
					RestartPolicy: corev1.RestartPolicyOnFailure,
				},
			},
		},
	}

	if err := controllerutil.SetControllerReference(instance, job, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	// TODO(user): Change this for the object type created by your controller
	// Check if the Deployment already exists
	{
		found := &batchv1.Job{}
		err = r.Get(context.TODO(), types.NamespacedName{Name: job.Name, Namespace: job.Namespace}, found)
		if err != nil && errors.IsNotFound(err) {
			log.Info("Creating Job", "namespace", job.Namespace, "name", job.Name)
			err = r.Create(context.TODO(), job)
		} else if err != nil {
			return reconcile.Result{}, err
		}

		// TODO(user): Change this for the object type created by your controller
		// Update the found object and write the result back if there are any changes
		if !reflect.DeepEqual(job.Spec.Template.Spec, found.Spec.Template.Spec) {
			found.Spec.Template.Spec = job.Spec.Template.Spec
			log.Info("Updating Job", "namespace", job.Namespace, "name", job.Name)
			//err = r.Update(context.TODO(), found)
			if err != nil {
				return reconcile.Result{}, err
			}
			log.Info("Skipped update")
		}
	}

	// Define the desired Deployment object
	deploy := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.Name + "-deployment",
			Namespace: ns.Name,
			Labels:    map[string]string{"dnsPrefix": instance.Spec.DnsPrefix},
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"deployment": instance.Name + "-deployment"},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{"deployment": instance.Name + "-deployment"}},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "nginx",
							Image: "nginx",
						},
					},
				},
			},
		},
	}
	if err := controllerutil.SetControllerReference(instance, deploy, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	// TODO(user): Change this for the object type created by your controller
	// Check if the Deployment already exists
	found := &appsv1.Deployment{}
	err = r.Get(context.TODO(), types.NamespacedName{Name: deploy.Name, Namespace: deploy.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		log.Info("Creating Deployment", "namespace", deploy.Namespace, "name", deploy.Name)
		err = r.Create(context.TODO(), deploy)
		return reconcile.Result{}, err
	} else if err != nil {
		return reconcile.Result{}, err
	}

	// TODO(user): Change this for the object type created by your controller
	// Update the found object and write the result back if there are any changes
	if !reflect.DeepEqual(deploy.Spec, found.Spec) {
		found.Spec = deploy.Spec
		log.Info("Updating Deployment", "namespace", deploy.Namespace, "name", deploy.Name)
		err = r.Update(context.TODO(), found)
		if err != nil {
			return reconcile.Result{}, err
		}
	}
	return reconcile.Result{}, nil
}
