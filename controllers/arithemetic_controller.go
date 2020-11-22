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
	"fmt"
	"io/ioutil"
	"time"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	mathsv1 "math-controller/api/v1"
)

// ArithemeticReconciler reconciles a Arithemetic object
type ArithemeticReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=maths.stream.com,resources=arithemetics,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=maths.stream.com,resources=arithemetics/status,verbs=get;update;patch

func (r *ArithemeticReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("arithemetic", req.NamespacedName)

	var problem mathsv1.Arithemetic
	if err := r.Get(ctx, req.NamespacedName, &problem); err != nil {
		log.Error(err, "could not get the Arithematic object")
		return ctrl.Result{}, err
	}

	if problem.Status.Answer == "" {
		log.Info(fmt.Sprintf("Reconciling for %s", req.NamespacedName))
		log.Info(fmt.Sprintf("Expression: %s", problem.Spec.Expression))

		pod := corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:      fmt.Sprintf("job-%s", req.Name),
				Namespace: "default",
			},
			Spec: corev1.PodSpec{
				RestartPolicy: "Never",
				Containers: []corev1.Container{
					{
						Name:  "problem-solver",
						Image: "python:latest",
						Args:  []string{"python", "-c", fmt.Sprintf("print(%s)", problem.Spec.Expression)},
					},
				},
			},
		}

		if err := r.Create(ctx, &pod, &client.CreateOptions{}); err != nil {
			log.Error(err, "could not create the container")
			return ctrl.Result{}, err
		}
		log.Info("Created the container")
		time.Sleep(10 * time.Second)

		answer, err := readPodLogs(pod)
		if err != nil {
			log.Error(err, "could not read logs")
			return ctrl.Result{}, err
		}

		log.Info(fmt.Sprintf("Answer is %s", answer))

		problem.Status.Answer = answer
		if err := r.Update(ctx, &problem, &client.UpdateOptions{}); err != nil {
			log.Error(err, "could not update resource")
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

func readPodLogs(pod corev1.Pod) (string, error) {
	config := ctrl.GetConfigOrDie()
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return "", err
	}

	req := clientSet.CoreV1().Pods(pod.Namespace).GetLogs(pod.Name, &corev1.PodLogOptions{})

	reader, err := req.Stream()
	if err != nil {
		return "", err
	}

	defer reader.Close()

	answer, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", err
	}

	return string(answer), nil
}

func (r *ArithemeticReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&mathsv1.Arithemetic{}).
		Complete(r)
}
