/*

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

package controllers

import (
	"context"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"strconv"
	"strings"

	bindingv1alpha1 "github.com/trisberg/binding/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// SecretMappingReconciler reconciles a SecretMapping object
type SecretMappingReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	Log    logr.Logger
}

// +kubebuilder:rbac:groups=binding.projectriff.io,resources=secretmappings,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=binding.projectriff.io,resources=secretmappings/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=core,resources=secrets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=secrets/status,verbs=get;update;patch

func (r *SecretMappingReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	_ = context.Background()
	log := r.Log.WithValues("secretmapping", req.NamespacedName)

	// your logic here
	instance := &bindingv1alpha1.SecretMapping{}
	err := r.Get(context.TODO(), req.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Object not found, return.  Created objects are automatically garbage collected.
			// For additional cleanup logic use finalizers.
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}
	log.Info("SecretMapping found", "name", instance.Name)

	bindingSecretName := instance.Name + "-binding"
	if instance.Spec.BindingSecret != "" {
		bindingSecretName = instance.Spec.BindingSecret
	}
	bindingPrefix := ""
	bindingKey := "config.yaml"
	if instance.Spec.BindingPrefix != "" {
		bindingPrefix = instance.Spec.BindingPrefix
	}
	log.Info("Creating binding secret", "name", bindingSecretName, "prefix", bindingPrefix)
	bindingSecret := &corev1.Secret{
		Type: corev1.SecretTypeOpaque,
		ObjectMeta: metav1.ObjectMeta{
			Name:      bindingSecretName,
			Namespace: instance.Namespace,
		},
		StringData: map[string]string{
			"source": instance.Name,
		},
	}

	config := ""
	indent := ""
	if len(bindingPrefix) > 0 {
		parts := strings.Split(bindingPrefix, ".")
		for _, p := range parts {
			config += indent + p + ":\n"
			indent += "  "
		}
	}

	if instance.Spec.URI != "" || instance.Spec.URIKey != "" {
		uri := instance.Spec.URI
		if instance.Spec.URIKey != "" {
			value := getSecretValue(r.Client, context.TODO(), instance.Namespace, instance.Spec.SecretRef, instance.Spec.URIKey)
			if len(value) > 0 {
				uri = string(value)
			}
		}
		if bindingPrefix == "spring.datasource" {
			config += indent + "url: " + "jdbc:" + uri + "\n"
		} else {
			config += indent + "uri: " + uri + "\n"
		}
	}
	if instance.Spec.Host != "" || instance.Spec.HostKey != "" {
		host := instance.Spec.Host
		if instance.Spec.HostKey != "" {
			value := getSecretValue(r.Client, context.TODO(), instance.Namespace, instance.Spec.SecretRef, instance.Spec.HostKey)
			if len(value) > 0 {
				host = string(value)
			}
		}
		config += indent + "host: " + host + "\n"
	}
	if instance.Spec.Port > 0 || instance.Spec.PortKey != "" {
		port := instance.Spec.Port
		if instance.Spec.PortKey != "" {
			value := getSecretValue(r.Client, context.TODO(), instance.Namespace, instance.Spec.SecretRef, instance.Spec.PortKey)
			if len(value) > 0 {
				port, _ = strconv.Atoi(string(value))
			}
		}
		config += indent + "port: " + strconv.Itoa(port) + "\n"
	}
	if instance.Spec.Username != "" || instance.Spec.UsernameKey != "" {
		username := instance.Spec.Username
		if instance.Spec.UsernameKey != "" {
			value := getSecretValue(r.Client, context.TODO(), instance.Namespace, instance.Spec.SecretRef, instance.Spec.UsernameKey)
			if len(value) > 0 {
				username = string(value)
			}
		}
		config += indent + "username: " + username + "\n"
	}
	if instance.Spec.SecretRef != "" {
		passwordKey := "password"
		if instance.Spec.PasswordKey != "" {
			passwordKey = instance.Spec.PasswordKey
		}
		password := getSecretValue(r.Client, context.TODO(), instance.Namespace, instance.Spec.SecretRef, passwordKey)
		config += indent + "password: " + string(password) + "\n"
	}
	bindingSecret.StringData[bindingKey] = string(config)

	log.Info("Setting owner for binding secret", "name", bindingSecretName)
	if err := ctrl.SetControllerReference(instance, bindingSecret, r.Scheme); err != nil {
		return reconcile.Result{}, err
	}

	log.Info("Look for existing binding secret", "name", bindingSecretName)

	// TODO(user): Change this for the object type created by your controller
	// Check if the Secret already exists
	found := &corev1.Secret{}
	err = r.Get(context.TODO(), types.NamespacedName{Name: bindingSecret.Name, Namespace: bindingSecret.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		log.Info("Creating Secret", "namespace", bindingSecret.Namespace, "name", bindingSecret.Name)
		err = r.Create(context.TODO(), bindingSecret)
		return reconcile.Result{}, err
	} else if err != nil {
		return reconcile.Result{}, err
	}

	// TODO(user): Change this for the object type created by your controller
	// Update the found object and write the result back if there are any changes
	//if !reflect.DeepEqual(bindingSecret.Data, found.Data) {
	found.StringData = bindingSecret.StringData
	log.Info("Updating Secret", "namespace", bindingSecret.Namespace, "name", bindingSecret.Name)
	err = r.Update(context.TODO(), found)
	if err != nil {
		return reconcile.Result{}, err
	}
	//}
	return ctrl.Result{}, nil
}

func getSecretValue(c client.Client, ctx context.Context, namespace string, name string, key string) []byte {
	secret := &corev1.Secret{}
	err := c.Get(ctx, types.NamespacedName{Name: name, Namespace: namespace}, secret)
	if err != nil {
		return nil
	}
	return secret.Data[key]
}

func (r *SecretMappingReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&bindingv1alpha1.SecretMapping{}).
		Complete(r)
}
