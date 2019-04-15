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

package modelservice

import (
	"context"
	"fmt"
	"reflect"

	mirv1beta1 "github.com/rakelkar/mir/pkg/apis/mir/v1beta1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	extv1beta1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
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

// Add creates a new ModelService Controller and adds it to the Manager with default RBAC. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileModelService{Client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("modelservice-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to ModelService
	err = c.Watch(&source.Kind{Type: &mirv1beta1.ModelService{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create
	// Uncomment watch a Deployment created by ModelService - change this for objects you create
	err = c.Watch(&source.Kind{Type: &appsv1.Deployment{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &mirv1beta1.ModelService{},
	})
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileModelService{}

// ReconcileModelService reconciles a ModelService object
type ReconcileModelService struct {
	client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a ModelService object and makes changes based on the state read
// and what is in the ModelService.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  The scaffolding writes
// a Deployment as an example
// Automatically generate RBAC rules to allow the Controller to read and write Deployments
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=deployments/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=mir.k8s.io,resources=modelservices,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=mir.k8s.io,resources=modelservices/status,verbs=get;update;patch
func (r *ReconcileModelService) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	// Fetch the ModelService instance
	instance := &mirv1beta1.ModelService{}
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

	// TODO: stuff that needs to come labels (set by mutating admission controller)
	mir_dns_prefix := "haha"
	namespace := instance.Namespace

	// Define the desired Service object
	if instance.Spec.Default.Custom == nil ||
		instance.Spec.Default.Custom.Container.Ports == nil {
		log.Info("Invalid spec", "name", instance.Name)
		return reconcile.Result{}, fmt.Errorf("invalid spec, must have custom section")
	}
	// for now assume the spec has a single port that maps to the http port
	if len(instance.Spec.Default.Custom.Container.Ports) != 1 ||
		instance.Spec.Default.Custom.Container.Ports[0].Protocol != v1.ProtocolTCP {
		log.Info("Port missing or wrong protocol", "name", instance.Name)
		return reconcile.Result{}, fmt.Errorf("invalid spec, missing ports")
	}

	containerPort := instance.Spec.Default.Custom.Container.Ports[0]
	servicePort := 80
	servicePorts := []v1.ServicePort{
		v1.ServicePort{
			Protocol:   v1.ProtocolTCP,
			Port:       int32(servicePort),
			TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: containerPort.ContainerPort},
		},
	}

	service := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.Name + "-svc",
			Namespace: namespace,
		},
		Spec: v1.ServiceSpec{
			Selector: map[string]string{"deployment": instance.Name + "-deployment"},
			Ports:    servicePorts,
		},
	}
	if err := controllerutil.SetControllerReference(instance, service, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	{
		// TODO(user): Change this for the object type created by your controller
		// Check if the Service already exists
		isFound := true
		found := &v1.Service{}
		err = r.Get(context.TODO(), types.NamespacedName{Name: service.Name, Namespace: service.Namespace}, found)
		if err != nil && errors.IsNotFound(err) {
			log.Info("Creating Service", "namespace", service.Namespace, "name", service.Name)
			isFound = false
			err = r.Create(context.TODO(), service)
		}
		if err != nil {
			return reconcile.Result{}, err
		}

		// TODO(user): Change this for the object type created by your controller
		// Update the found object and write the result back if there are any changes
		if isFound && !reflect.DeepEqual(found.Spec.Ports, service.Spec.Ports) {
			found.Spec.Ports = service.Spec.Ports
			log.Info("Updating Service", "namespace", service.Namespace, "name", service.Name)
			err = r.Update(context.TODO(), found)
			if err != nil {
				return reconcile.Result{}, err
			}
		}
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
			Namespace: namespace,
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
		isFound := true
		found := &extv1beta1.Ingress{}
		err = r.Get(context.TODO(), types.NamespacedName{Name: ingress.Name, Namespace: ingress.Namespace}, found)
		if err != nil && errors.IsNotFound(err) {
			log.Info("Creating Ingress", "namespace", ingress.Namespace, "name", ingress.Name)
			isFound = false
			err = r.Create(context.TODO(), ingress)
		}
		if err != nil {
			return reconcile.Result{}, err
		}

		// TODO(user): Change this for the object type created by your controller
		// Update the found object and write the result back if there are any changes
		if isFound && !reflect.DeepEqual(ingress.Spec, ingress.Spec) {
			found.Spec = ingress.Spec
			log.Info("Updating Ingress", "namespace", ingress.Namespace, "name", ingress.Name)
			err = r.Update(context.TODO(), found)
			if err != nil {
				return reconcile.Result{}, err
			}
		}
	}

	// TODO(user): Change this to be the object type created by your controller
	// Define the desired Deployment object
	containerSpec := instance.Spec.Default.Custom.Container.DeepCopy()

	deploy := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.Name + "-deployment",
			Namespace: namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"deployment": instance.Name + "-deployment"},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{"deployment": instance.Name + "-deployment"}},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						*containerSpec,
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

	// // TODO(user): Change this for the object type created by your controller
	// // Update the found object and write the result back if there are any changes
	// if !reflect.DeepEqual(deploy.Spec, found.Spec) {
	// 	found.Spec = deploy.Spec
	// 	log.Info("Updating Deployment", "namespace", deploy.Namespace, "name", deploy.Name)
	// 	err = r.Update(context.TODO(), found)
	// 	if err != nil {
	// 		return reconcile.Result{}, err
	// 	}
	// }

	return reconcile.Result{}, nil
}
