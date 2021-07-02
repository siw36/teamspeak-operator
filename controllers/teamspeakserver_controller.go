/*
Copyright 2021.

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

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	voipv1alpha1 "github.com/siw36/teamspeak-operator/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TeamSpeakServerReconciler reconciles a TeamSpeakServer object
type TeamSpeakServerReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=voip.replicas.io,resources=teamspeakservers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=voip.replicas.io,resources=teamspeakservers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=voip.replicas.io,resources=teamspeakservers/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the TeamSpeakServer object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.8.3/pkg/reconcile
func (r *TeamSpeakServerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	// Lookup the TeamSpeakServer instance for this reconcile request
	tsCR := &voipv1alpha1.TeamSpeakServer{}
	err := r.Get(ctx, req.NamespacedName, tsCR)
	return ctrl.Result{Requeue: true}, err
}

// SetupWithManager sets up the controller with the Manager.
func (r *TeamSpeakServerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&voipv1alpha1.TeamSpeakServer{}).
		Owns(&appsv1.Deployment{}).
		Complete(r)
}

func (r *TeamSpeakServerReconciler) teamspeakserverDeployment(tsCR *voipv1alpha1.TeamSpeakServer, imageURL string, imageTag string) *appsv1.Deployment {

	databasePassword := err := r.Get(ctx, req.NamespacedName, tsCR)

	labels := map[string]string{
		"app": tsCR.Name,
	}
	matchlabels := map[string]string{
		"app": tsCR.Name,
	}

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "teamspeak-server",
			Namespace: tsCR.Namespace,
			Labels:    labels,
		},

		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: matchlabels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: matchlabels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image: tsCR.Spec.TeamSpeak.Image,
						Name:  "teamspeak",

						Env: []corev1.EnvVar{{
							Name:  "TS3SERVER_DB_PLUGIN",
							Value: "ts3db_mariadb",
						},
							{
								Name:  "TS3SERVER_DB_SQLCREATEPATH",
								Value: "create_mariadb",
							},
							{
								Name:  "TS3SERVER_DB_HOST",
								Value: tsCR.Spec.Database.ConfigUnmanaged.Host,
							},
							{
								Name:  "TS3SERVER_DB_USER",
								Value: tsCR.Spec.Database.ConfigUnmanaged.User,
							},
							{
								Name:  "TS3SERVER_DB_PASSWORD",
								Value: tsCR.Spec.Database.ConfigUnmanaged.User,
							},
							{
								Name:  "TS3SERVER_DB_NAME",
								Value: tsCR.Spec.Database.ConfigUnmanaged.Database,
							},
							{
								Name:  "TS3SERVER_DB_WAITUNTILREADY",
								Value: "5",
							},
							{
								Name:  "TS3SERVER_LICENSE",
								Value: "accept",
							},
						},

						Ports: []corev1.ContainerPort{{
							ContainerPort: 80,
							Name:          "wordpress-port",
						}},
						VolumeMounts: []corev1.VolumeMount{
							{
								Name:      "wordpress-persistent-storage",
								MountPath: "/var/www/html",
							},
						},
					},
					},

					Volumes: []corev1.Volume{

						{
							Name: "wordpress-persistent-storage",
							VolumeSource: corev1.VolumeSource{

								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
									ClaimName: "wp-pv-claim",
								},
							},
						},
					},
				},
			},
		},
	}

	controllerutil.SetControllerReference(tsCR, deployment, r.scheme)
	return dep

}
