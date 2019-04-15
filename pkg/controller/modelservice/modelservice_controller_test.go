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
	"testing"
	"time"

	"github.com/onsi/gomega"
	mirv1beta1 "github.com/rakelkar/mir/pkg/apis/mir/v1beta1"
	"golang.org/x/net/context"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var c client.Client

var expectedRequest = reconcile.Request{NamespacedName: types.NamespacedName{Name: "foo", Namespace: "source-ns"}}
var depKey = types.NamespacedName{Name: "foo-deployment", Namespace: "model-ns"}

const timeout = time.Second * 5

func TestReconcile(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	instance := &mirv1beta1.ModelService{
		ObjectMeta: metav1.ObjectMeta{Name: "foo", Namespace: "source-ns"},
		Spec: mirv1beta1.ModelServiceSpec{
			Default: mirv1beta1.ModelSpec{
				Custom: &mirv1beta1.CustomSpec{
					Container: v1.Container{
						Name:  "echo",
						Image: "hashicorp/http-echo",
						Args:  []string{"-text=hello"},
						Ports: []v1.ContainerPort{
							v1.ContainerPort{
								ContainerPort: 8888,
								Protocol:      v1.ProtocolTCP,
							},
						},
					},
				},
			},
		},
	}

	// Setup the Manager and Controller.  Wrap the Controller Reconcile function so it writes each request to a
	// channel when it is finished.
	mgr, err := manager.New(cfg, manager.Options{})
	g.Expect(err).NotTo(gomega.HaveOccurred())
	c = mgr.GetClient()

	// setup expected namespaces
	sourceNs := &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: "source-ns",
			Labels: map[string]string{
				"mir":            "someMir",
				"mir-dns-prefix": "some-prefix",
				"model-ns":       "model-ns",
				"modelsource":    "someSource",
			},
		},
		Spec: v1.NamespaceSpec{},
	}
	err = c.Create(context.TODO(), sourceNs)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	modelNs := &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: "model-ns",
			Labels: map[string]string{
				"mir":            "someMir",
				"mir-dns-prefix": "some-prefix",
			},
		},
		Spec: v1.NamespaceSpec{},
	}
	err = c.Create(context.TODO(), modelNs)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	recFn, requests := SetupTestReconcile(newReconciler(mgr))
	g.Expect(add(mgr, recFn)).NotTo(gomega.HaveOccurred())

	stopMgr, mgrStopped := StartTestManager(mgr, g)

	defer func() {
		close(stopMgr)
		mgrStopped.Wait()
	}()

	// Create the ModelService object and expect the Reconcile and Deployment to be created
	err = c.Create(context.TODO(), instance)
	// The instance object may not be a valid object because it might be missing some required fields.
	// Please modify the instance object by adding required fields and then remove the following if statement.
	if apierrors.IsInvalid(err) {
		t.Logf("failed to create object, got an invalid object error: %v", err)
		return
	}
	g.Expect(err).NotTo(gomega.HaveOccurred())
	defer c.Delete(context.TODO(), instance)
	g.Eventually(requests, timeout).Should(gomega.Receive(gomega.Equal(expectedRequest)))

	deploy := &appsv1.Deployment{}
	g.Eventually(func() error { return c.Get(context.TODO(), depKey, deploy) }, timeout).
		Should(gomega.Succeed())

	// Delete the Deployment and expect Reconcile to be called for Deployment deletion
	g.Expect(c.Delete(context.TODO(), deploy)).NotTo(gomega.HaveOccurred())
	//g.Eventually(requests, timeout).Should(gomega.Receive(gomega.Equal(expectedRequest)))
	// g.Eventually(func() error { return c.Get(context.TODO(), depKey, deploy) }, timeout).
	// 	Should(gomega.Succeed())

	// Manually delete Deployment since GC isn't enabled in the test control plane
	g.Eventually(func() error { return c.Delete(context.TODO(), deploy) }, timeout).
		Should(gomega.MatchError("deployments.apps \"foo-deployment\" not found"))
}
