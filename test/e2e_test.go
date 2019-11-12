package test

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"

	"github.com/mrogers950/file-integrity-operator/pkg/common"
)

const (
	// How often to poll for conditions
	Poll = 2 * time.Second
	// Default time to wait for operations to complete
	defaultTimeout = 200 * time.Second
)

func TestFileIntegrityOperatorE2ESuite(t *testing.T) {
	tests := map[string]struct {
		thing string
	}{
		"basic": {
			thing: "",
		},
	}

	kubeConfig, err := loadConfig(os.Getenv("KUBECONFIG"), os.Getenv("KUBECONTEXT"))
	if err != nil {
		t.Fatalf("error loading kubeconfig: '%s', ctx: '%s', err: %s", os.Getenv("KUBECONFIG"), os.Getenv("KUBECONTEXT"), err)
	}
	kubeClientSet, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		t.Fatalf("error generating kube client: %s", err)
	}

	for tcName, _ := range tests {
		t.Run(fmt.Sprintf("setting up e2e tests %s", tcName), func(t *testing.T) {
			execCmd("oc", []string{"create", "-f", "../deploy/"}, "")
			execCmd("oc", []string{"create", "-f", "../deploy/crds"}, "")
			defer func() {
				_, err := kubeClientSet.AppsV1().DaemonSets(common.FileIntegrityNamespace).Create(cleanAideDaemonset())
				if err != nil {
					t.Error(err)
				}
				time.Sleep(5 * time.Second)
				//execCmd("oc", []string{"delete", "-f", "../deploy/operator.yaml"}, "")
				//execCmd("oc", []string{"delete", "-f", "../deploy/crds"}, "")
			}()

			err := waitForOperatorPod(kubeClientSet)
			if err != nil {
				t.Error(err)
			}

			t.Log("operator pod OK")
			err = waitForDaemonSet(kubeClientSet)
			if err != nil {
				t.Error(err)
			}

			t.Log("ds OK")

		})
	}
}

func aideDaemonSetIsReady(c kubernetes.Interface) wait.ConditionFunc {
	return func() (bool, error) {
		masterDaemonSet, err := c.AppsV1().DaemonSets(common.FileIntegrityNamespace).Get(common.DaemonSetName, metav1.GetOptions{})
		if err != nil && !kerr.IsNotFound(err) {
			return false, err
		}
		if kerr.IsNotFound(err) {
			return false, nil
		}
		if masterDaemonSet.Status.DesiredNumberScheduled != masterDaemonSet.Status.NumberAvailable {
			return false, nil
		}
		return true, nil
	}
}

func findRunningOperatorPod(c kubernetes.Interface) wait.ConditionFunc {
	return func() (bool, error) {
		var foundPod *corev1.Pod
		podList, err := c.CoreV1().Pods(common.FileIntegrityNamespace).List(metav1.ListOptions{})
		if err != nil {
			return false, err
		}
		for _, pod := range podList.Items {
			if strings.HasPrefix(pod.Name, "file-integrity-operator") {
				foundPod = &pod
				break
			}
		}
		if foundPod == nil {
			return false, nil
		}

		switch foundPod.Status.Phase {
		case corev1.PodRunning:
			return true, nil
		case corev1.PodFailed, corev1.PodSucceeded:
			return false, fmt.Errorf("pod ran to completion")
		}
		return false, nil
	}

}

func podRunning(c kubernetes.Interface, podName, namespace string) wait.ConditionFunc {
	return func() (bool, error) {
		pod, err := c.CoreV1().Pods(namespace).Get(podName, metav1.GetOptions{})
		if err != nil {
			return false, err
		}
		switch pod.Status.Phase {
		case corev1.PodRunning:
			return true, nil
		case corev1.PodFailed, corev1.PodSucceeded:
			return false, fmt.Errorf("pod ran to completion")
		}
		return false, nil
	}
}

func restClientConfig(config, context string) (*api.Config, error) {
	if config == "" {
		return nil, fmt.Errorf("Config file must be specified to load client config")
	}
	c, err := clientcmd.LoadFromFile(config)
	if err != nil {
		return nil, fmt.Errorf("error loading config: %v", err.Error())
	}
	if context != "" {
		c.CurrentContext = context
	}
	return c, nil
}

func loadConfig(config, context string) (*rest.Config, error) {
	c, err := restClientConfig(config, context)
	if err != nil {
		return nil, err
	}
	return clientcmd.NewDefaultClientConfig(*c, &clientcmd.ConfigOverrides{}).ClientConfig()
}

// execCmd executes a command and returns the stdout + error, if any
func execCmd(cmd string, args []string, input string) (string, error) {
	c := exec.Command(cmd, args...)
	stdin, err := c.StdinPipe()
	if err != nil {
		return "", err
	}

	go func() {
		defer stdin.Close()
		if input != "" {
			io.WriteString(stdin, input)
		}
	}()

	out, err := c.CombinedOutput()
	if err != nil {
		fmt.Printf("Command '%s' failed with: %s\n", cmd, err)
		fmt.Printf("Output: %s\n", out)
		return "", err
	}
	return string(out), nil
}

// Waits default amount of time (PodStartTimeout) for the specified pod to become running.
// Returns an error if timeout occurs first, or pod goes in to failed state.
func waitForPodRunningInNamespace(c kubernetes.Interface, pod *corev1.Pod) error {
	if pod.Status.Phase == corev1.PodRunning {
		return nil
	}
	return waitTimeoutForPodRunningInNamespace(c, pod.Name, pod.Namespace, defaultTimeout)
}

func waitTimeoutForPodRunningInNamespace(c kubernetes.Interface, podName, namespace string, timeout time.Duration) error {
	return wait.PollImmediate(Poll, defaultTimeout, podRunning(c, podName, namespace))
}

func waitForDaemonSet(c kubernetes.Interface) error {
	return wait.PollImmediate(Poll, defaultTimeout, aideDaemonSetIsReady(c))
}

func waitForOperatorPod(c kubernetes.Interface) error {
	return wait.PollImmediate(Poll, defaultTimeout, findRunningOperatorPod(c))
}

// This daemonSet runs a command to clear the aide content from the host
func cleanAideDaemonset() *appsv1.DaemonSet {
	priv := true
	runAs := int64(0)

	return &appsv1.DaemonSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "aide-clean",
			Namespace: common.FileIntegrityNamespace,
		},
		Spec: appsv1.DaemonSetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "aide-clean",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "aide-clean",
					},
				},
				Spec: corev1.PodSpec{
					Tolerations: []corev1.Toleration{
						{
							Key:      "node-role.kubernetes.io/master",
							Operator: "Exists",
							Effect:   "NoSchedule",
						},
					},
					ServiceAccountName: common.OperatorServiceAccountName,
					Containers: []corev1.Container{
						{
							SecurityContext: &corev1.SecurityContext{
								Privileged: &priv,
								RunAsUser:  &runAs,
							},
							Name:    "aide-clean",
							Image:   "busybox",
							Command: []string{"rm"},
							Args:    []string{"-f", "/hostroot/etc/kubernetes/aide*"},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "hostroot",
									MountPath: "/hostroot",
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "hostroot",
							VolumeSource: corev1.VolumeSource{
								HostPath: &corev1.HostPathVolumeSource{
									Path: "/",
								},
							},
						},
					},
				},
			},
		},
	}
}
