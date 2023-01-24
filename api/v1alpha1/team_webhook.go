/*
Copyright 2023.

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

package v1alpha1

import (
	"context"
	"fmt"
	"os"

	authv1 "k8s.io/api/authorization/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var teamlog = logf.Log.WithName("team-resource")

func (r *Team) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
//+kubebuilder:webhook:path=/validate-team-snappcloud-io-v1alpha1-team,mutating=false,failurePolicy=fail,sideEffects=None,groups=team.snappcloud.io,resources=teams,verbs=create;update,versions=v1alpha1,name=vteam.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &Team{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *Team) ValidateCreate() error {
	teamlog.Info("validate create", "name", r.Name)
	fmt.Println("hello")

	// TODO(user): fill in your validation logic upon object creation.
	return nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *Team) ValidateUpdate(old runtime.Object) error {

	teamlog.Info("validate update", "name", r.Name)

	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *Team) ValidateDelete() error {
	teamlog.Info("validate delete", "name", r.Name)
	// TODO(user): fill in your validation logic upon object deletion.
	return nil
}

func (r *Team) ValidateUserAccess(ctx context.Context) error {

	kubeconfig := fmt.Sprintf("%s/.kube/config-anis", os.Getenv("HOME"))
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	action := authv1.ResourceAttributes{
		Namespace: "snappcloud-operators-seldon",
		Verb:      "create",
		Resource:  "rolebinding",
	}
	selfCheck := authv1.SelfSubjectAccessReview{
		Spec: authv1.SelfSubjectAccessReviewSpec{
			ResourceAttributes: &action,
		},
	}

	resp, err := clientset.AuthorizationV1().
		SelfSubjectAccessReviews().
		Create(context.TODO(), &selfCheck, metav1.CreateOptions{})

	if err != nil {
		panic(err.Error())
	}

	if resp.Status.Allowed {
		fmt.Println("allowed")
	} else {
		fmt.Println("denied")
	}
	return nil

}
