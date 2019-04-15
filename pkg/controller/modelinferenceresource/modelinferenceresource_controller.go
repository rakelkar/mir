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
	"fmt"
	"os"
	"reflect"
	"strings"

	"k8s.io/apimachinery/pkg/util/intstr"

	mirv1beta1 "github.com/rakelkar/mir/pkg/apis/mir/v1beta1"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	extv1beta1 "k8s.io/api/extensions/v1beta1"
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

// to create the secret this fn needs on your cluster, run:
// kubectl create secret generic --namespace=${K8S_NAMESPACE} azcreds --from-literal=AZ_CLIENT_ID=${AZ_CLIENT_ID} --from-literal=AZ_CLIENT_SECRET=${AZ_CLIENT_SECRET} --from-literal=AZ_TENANT_ID=${AZ_TENANT_ID} --from-literal=AZ_SUBSCRIPTION_ID=${AZ_SUBSCRIPTION_ID}  --from-literal=AZ_SP_OBJECT_ID=${AZ_SP_OBJECT_ID}
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

	// init values TODO: do somewhere else
	mir_scale_unit := os.Getenv("MIR_SCALE_UNIT")
	if len(mir_scale_unit) == 0 {
		mir_scale_unit = "su1"
	}
	mir_location := os.Getenv("MIR_LOCATION")
	if len(mir_location) == 0 {
		mir_location = "eastus"
	}
	mir_tm_endpoint_address := os.Getenv("MIR_ENDPOINT_ADDRESS")
	if len(mir_tm_endpoint_address) == 0 {
		// some crappy address TODO: get this from the cluster
		mir_tm_endpoint_address = "40.70.209.164"
	}
	mir_dns_prefix := strings.ToLower(instance.Spec.DnsPrefix)
	mir_resource_group := fmt.Sprintf("mir-tms-%s-%s", mir_scale_unit, mir_location)
	mir_tm_name := mir_dns_prefix + "-mir-tm"
	mir_tm_endpoint_name := mir_scale_unit

	// add or execute finalizer
	myFinalizerName := "tm.finalizer.azureml.microsoft.com"
	if instance.ObjectMeta.DeletionTimestamp.IsZero() {
		// The object is not being deleted, so if it does not have our finalizer,
		// then lets add the finalizer and update the object.
		if !containsString(instance.ObjectMeta.Finalizers, myFinalizerName) {
			log.Info("add finalizer", "name", instance.Name)
			instance.ObjectMeta.Finalizers = append(instance.ObjectMeta.Finalizers, myFinalizerName)
			if err := r.Update(context.Background(), instance); err != nil {
				return reconcile.Result{}, err
			}
		}
	} else {
		// The object is being deleted
		if containsString(instance.ObjectMeta.Finalizers, myFinalizerName) {
			// our finalizer is present, so lets handle our external dependency
			// Ensure that delete implementation is idempotent and safe to invoke
			// multiple types for same object.
			log.Info("deleting the external dependencies", "name", instance.Name)
			cmd := []string{"/entrypoint.sh", "/ep_dtor.sh", mir_resource_group, mir_location, mir_tm_name, mir_dns_prefix, mir_tm_endpoint_name, mir_tm_endpoint_address}
			job := constructJob(instance.Name+"-"+string(instance.UID)+"-dtor", instance.Namespace, cmd)
			found := &batchv1.Job{}
			err := r.Get(context.TODO(), types.NamespacedName{Name: job.Name, Namespace: job.Namespace}, found)
			if err != nil && errors.IsNotFound(err) {
				log.Info("Creating Job", "namespace", job.Namespace, "name", job.Name)
				err = r.Create(context.TODO(), job)
			}
			if err != nil {
				// if fail to delete the external dependency here, return with error
				// so that it can be retried
				return reconcile.Result{}, err
			}

			// remove our finalizer from the list and update it.
			instance.ObjectMeta.Finalizers = removeString(instance.ObjectMeta.Finalizers, myFinalizerName)
			if err := r.Update(context.Background(), instance); err != nil {
				return reconcile.Result{}, err
			}
		}

		// Our finalizer has finished, so the reconciler can do nothing.
		return reconcile.Result{}, nil
	}

	// TODO(user): Change this to be the object type created by your controller
	ns := &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: instance.Name + "-ns",
			Labels: map[string]string{
				"mir":           instance.Name,
				"mir-dns-prfix": mir_dns_prefix,
			},
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
		err = r.Get(context.TODO(), types.NamespacedName{Name: ns.Name, Namespace: ""}, found)
		if err != nil && errors.IsNotFound(err) {
			log.Info("Creating Namespace", "namespace", "default", "name", ns.Name)
			err = r.Create(context.TODO(), ns)
		}
		if err != nil {
			return reconcile.Result{}, err
		}
	}

	cmd := []string{"/entrypoint.sh", "/tm.sh", mir_resource_group, mir_location, mir_tm_name, mir_dns_prefix, mir_tm_endpoint_name, mir_tm_endpoint_address}
	job := constructJob(instance.Name+"-"+string(instance.UID)+"-ctor", instance.Namespace, cmd)
	if err := controllerutil.SetControllerReference(instance, job, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	// TODO(user): Change this for the object type created by your controller
	// Check if the Deployment already exists
	{
		isFound := true
		found := &batchv1.Job{}
		err = r.Get(context.TODO(), types.NamespacedName{Name: job.Name, Namespace: job.Namespace}, found)
		if err != nil && errors.IsNotFound(err) {
			log.Info("Creating Job", "namespace", job.Namespace, "name", job.Name)
			isFound = false
			err = r.Create(context.TODO(), job)
		}
		if err != nil {
			return reconcile.Result{}, err
		}

		// TODO(user): Change this for the object type created by your controller
		// Update the found object and write the result back if there are any changes
		if isFound && !reflect.DeepEqual(job.Spec.Template.Spec.Containers[0].Command, found.Spec.Template.Spec.Containers[0].Command) {
			found.Spec.Template.Spec.Containers = job.Spec.Template.Spec.Containers
			log.Info("Deleting Updated Job", "namespace", job.Namespace, "name", job.Name)
			err = r.Delete(context.TODO(), found)
			if err != nil && !errors.IsNotFound(err) {
				return reconcile.Result{}, err
			}

			return reconcile.Result{Requeue: true}, nil
		}
	}

	// Define the desired Service object
	servicePort := 80
	service := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.Name + "-service",
			Namespace: ns.Name,
		},
		Spec: v1.ServiceSpec{
			Selector: map[string]string{"deployment": instance.Name + "-deployment"},
			Ports: []v1.ServicePort{
				v1.ServicePort{
					Protocol:   v1.ProtocolTCP,
					Port:       int32(servicePort),
					TargetPort: intstr.FromInt(5678),
				},
			},
		},
	}
	if err := controllerutil.SetControllerReference(instance, service, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	{
		// TODO(user): Change this for the object type created by your controller
		// Check if the Service already exists
		// isFound := true
		found := &v1.Service{}
		err = r.Get(context.TODO(), types.NamespacedName{Name: service.Name, Namespace: service.Namespace}, found)
		if err != nil && errors.IsNotFound(err) {
			log.Info("Creating Service", "namespace", service.Namespace, "name", service.Name)
			// isFound = false
			err = r.Create(context.TODO(), service)
		}
		if err != nil {
			return reconcile.Result{}, err
		}

		// // TODO(user): Change this for the object type created by your controller
		// // Update the found object and write the result back if there are any changes
		// if isFound && !reflect.DeepEqual(found.Spec.Ports, service.Spec.Ports) {
		// 	found.Spec.Ports = service.Spec.Ports
		// 	log.Info("Updating Service", "namespace", service.Namespace, "name", service.Name)
		// 	err = r.Update(context.TODO(), found)
		// 	if err != nil {
		// 		return reconcile.Result{}, err
		// 	}
		// }
	}

	// Define the desired Ingress object
	ingressHost := mir_dns_prefix + ".trafficmanager.net"
	ingress := &extv1beta1.Ingress{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Ingress",
			APIVersion: "extensions/v1beta1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.Name,
			Namespace: ns.Name,
		},
		Spec: extv1beta1.IngressSpec{
			Rules: []extv1beta1.IngressRule{
				extv1beta1.IngressRule{
					Host: ingressHost,
					IngressRuleValue: extv1beta1.IngressRuleValue{
						HTTP: &extv1beta1.HTTPIngressRuleValue{
							Paths: []extv1beta1.HTTPIngressPath{
								extv1beta1.HTTPIngressPath{
									Path: "/",
									Backend: extv1beta1.IngressBackend{
										ServiceName: service.Name,
										ServicePort: intstr.FromInt(servicePort),
									},
								},
							},
						},
					},
				},
			},
		},
	}
	if err := controllerutil.SetControllerReference(instance, ingress, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	{
		// TODO(user): Change this for the object type created by your controller
		// Check if the Ingress already exists
		// isFound := true
		found := &extv1beta1.Ingress{}
		err = r.Get(context.TODO(), types.NamespacedName{Name: ingress.Name, Namespace: ingress.Namespace}, found)
		if err != nil && errors.IsNotFound(err) {
			log.Info("Creating Ingress", "namespace", ingress.Namespace, "name", ingress.Name)
			// isFound = false
			err = r.Create(context.TODO(), ingress)
		}
		if err != nil {
			return reconcile.Result{}, err
		}

		// // TODO(user): Change this for the object type created by your controller
		// // Update the found object and write the result back if there are any changes
		// if isFound && !reflect.DeepEqual(ingress.Spec, ingress.Spec) {
		// 	found.Spec = ingress.Spec
		// 	log.Info("Updating Ingress", "namespace", ingress.Namespace, "name", ingress.Name)
		// 	err = r.Update(context.TODO(), found)
		// 	if err != nil {
		// 		return reconcile.Result{}, err
		// 	}
		// }
	}

	// Define the desired Deployment object
	deploy := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.Name + "-deployment",
			Namespace: ns.Name,
			Labels:    map[string]string{"dnsPrefix": mir_dns_prefix},
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
							Name:  "echo",
							Image: "hashicorp/http-echo",
							Args:  []string{"-text=hello"},
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
	// isFound := true
	found := &appsv1.Deployment{}
	err = r.Get(context.TODO(), types.NamespacedName{Name: deploy.Name, Namespace: deploy.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		log.Info("Creating Deployment", "namespace", deploy.Namespace, "name", deploy.Name)
		// isFound = false
		err = r.Create(context.TODO(), deploy)
	}
	if err != nil {
		return reconcile.Result{}, err
	}

	// TODO(user): Change this for the object type created by your controller
	// Update the found object and write the result back if there are any changes
	// if isFound && !reflect.DeepEqual(deploy.Spec, found.Spec) {
	// 	found.Spec = deploy.Spec
	// 	log.Info("Updating Deployment", "namespace", deploy.Namespace, "name", deploy.Name)
	// 	err = r.Update(context.TODO(), found)
	// 	if err != nil {
	// 		return reconcile.Result{}, err
	// 	}
	// }

	return reconcile.Result{}, nil
}

func constructJob(name string, ns string, cmd []string) *batchv1.Job {
	az_vars := use_az_secret()
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: ns,
			Labels:    make(map[string]string),
		},
		Spec: batchv1.JobSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: make(map[string]string),
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:   name,
					Labels: make(map[string]string),
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:    "azcmd",
							Image:   "kcorer/azcmd",
							Command: cmd,
							Env:     az_vars,
						},
					},
					RestartPolicy: corev1.RestartPolicyOnFailure,
				},
			},
		},
	}

	return job
}

//
// Helper functions to check and remove string from a slice of strings.
//
func containsString(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

func removeString(slice []string, s string) (result []string) {
	for _, item := range slice {
		if item == s {
			continue
		}
		result = append(result, item)
	}
	return
}
