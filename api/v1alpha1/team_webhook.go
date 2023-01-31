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
	"errors"
	"fmt"

	authv1 "k8s.io/api/authorization/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
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

func (r *Team) ValidateCreate() error {
	teamlog.Info("validate create", "name", r.Name)
	fmt.Println("hello") // TODO(user): fill in your validation logic upon object creation.
	return nil
}

func (r *Team) ValidateUpdate(old runtime.Object) error {

	teamlog.Info("validate update", "name", r.Name)

	config, err := rest.InClusterConfig()
	if err != nil {
		teamlog.Error(err, "can not get incluster config")
	}

	config.Impersonate = rest.ImpersonationConfig{
		UserName: r.Spec.TeamAdmin,
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		teamlog.Error(err, "can not create clientset")
	}

	for _, ns := range r.Spec.Namespaces {
		//check if namespace does not exist or has been deleted
		teamns, err := clientset.CoreV1().Namespaces().Get(context.TODO(), ns, metav1.GetOptions{})
		if err != nil {
			errorresp := "namespace " + ns + " doese not exist"
			return errors.New(errorresp)
		}
		fmt.Println(teamns.Name)
		//check if namespace already has been added to another team
		if val, ok := teamns.Labels["snappcloud.io/team"]; ok {
			if teamns.Labels["snappcloud.io/team"] != r.Name {
				errresp := "namespace " + ns + " already have team label " + val
				return errors.New(errresp)
			}

		}
		if err != nil {
			teamlog.Error(err, "can not get ns")
		}

		//Check If user has access to this namespace
		action := authv1.ResourceAttributes{
			Namespace: ns,
			Verb:      "create",
			Resource:  "rolebindings",
			Group:     "rbac.authorization.k8s.io",
			Version:   "v1",
		}

		check := authv1.LocalSubjectAccessReview{
			Spec: authv1.SubjectAccessReviewSpec{
				User:               r.Spec.TeamAdmin,
				ResourceAttributes: &action,
			},
		}

		resp, err := clientset.AuthorizationV1().
			LocalSubjectAccessReviews(ns).
			Create(context.TODO(), &check, metav1.CreateOptions{})

		// selfCheck := authv1.SelfSubjectAccessReview{
		// 	Spec: authv1.SelfSubjectAccessReviewSpec{
		// 		ResourceAttributes: &action,
		// 	},
		// }

		// resp, err := clientset.AuthorizationV1().
		// 	SelfSubjectAccessReviews().
		// 	Create(context.TODO(), &selfCheck, metav1.CreateOptions{})

		if err != nil {
			teamlog.Error(err, "Can not create rolebingdin")
		}
		//test
		if resp.Status.Allowed {
			fmt.Println("Allowed")
		} else {
			fmt.Println("Denied")
			return errors.New("you are not allowed to add namespace " + ns)
		}
	}
	return nil
}

func (r *Team) ValidateDelete() error {
	teamlog.Info("validate delete", "name", r.Name)
	// TODO(user): fill in your validation logic upon object deletion.
	return nil
}
