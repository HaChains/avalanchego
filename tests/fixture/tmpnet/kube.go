// Copyright (C) 2019-2024, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package tmpnet

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/onsi/ginkgo/v2"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/transport/spdy"
	"k8s.io/utils/pointer"

	"github.com/ava-labs/avalanchego/api/info"
	"github.com/ava-labs/avalanchego/config"
	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/tests"
	"github.com/ava-labs/avalanchego/utils/constants"
	"github.com/ava-labs/avalanchego/utils/logging"
	"k8s.io/apimachinery/pkg/util/wait"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	restclient "k8s.io/client-go/rest"
)

// defaultNodeFlags defines common flags for avalanchego nodes used by this test
func DefaultPodFlags(networkName string, dataDir string) map[string]string {
	return map[string]string{
		config.DataDirKey:                dataDir,
		config.NetworkNameKey:            networkName,
		config.SybilProtectionEnabledKey: "false",
		config.HealthCheckFreqKey:        "500ms", // Ensure rapid detection of a healthy state
		config.LogDisplayLevelKey:        logging.Debug.String(),
		config.LogLevelKey:               logging.Debug.String(),
		config.HTTPHostKey:               "0.0.0.0", // Need to bind to pod IP to ensure kubelet can access the http port for the readiness check
	}
}

// newNodeStatefulSet returns a statefulset for an avalanchego node.
func NewNodeStatefulSet(
	name string,
	imageName string,
	containerName string,
	volumeName string,
	volumeSize string,
	volumeMountPath string,
	flags map[string]string,
) *appsv1.StatefulSet {
	return &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: name + "-",
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas:    pointer.Int32(1),
			ServiceName: name,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": name,
				},
			},
			VolumeClaimTemplates: []corev1.PersistentVolumeClaim{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: volumeName,
					},
					Spec: corev1.PersistentVolumeClaimSpec{
						AccessModes: []corev1.PersistentVolumeAccessMode{
							corev1.ReadWriteOnce,
						},
						Resources: corev1.VolumeResourceRequirements{
							Requests: corev1.ResourceList{
								corev1.ResourceStorage: resource.MustParse(volumeSize),
							},
						},
					},
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": name,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  containerName,
							Image: imageName,
							Ports: []corev1.ContainerPort{
								{
									Name:          "http",
									ContainerPort: config.DefaultHTTPPort,
								},
								{
									Name:          "staker",
									ContainerPort: config.DefaultStakingPort,
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      volumeName,
									MountPath: volumeMountPath,
								},
							},
							ReadinessProbe: &corev1.Probe{
								ProbeHandler: corev1.ProbeHandler{
									HTTPGet: &corev1.HTTPGetAction{
										Path: "/ext/health/liveness",
										Port: intstr.FromInt(config.DefaultHTTPPort),
									},
								},
								PeriodSeconds:    1,
								SuccessThreshold: 1,
							},
							Env: stringMapToEnvVarSlice(flags),
						},
					},
				},
			},
		},
	}
}

// stringMapToEnvVarSlice converts a string map to a kube EnvVar slice.
func stringMapToEnvVarSlice(mapping map[string]string) []corev1.EnvVar {
	envVars := make([]corev1.EnvVar, len(mapping))
	var i int
	for k, v := range mapping {
		envVars[i] = corev1.EnvVar{
			Name:  config.EnvVarName(config.EnvPrefix, k),
			Value: v,
		}
		i++
	}
	return envVars
}

// waitForPodCondition waits until the specified pod reports the specified condition
func waitForPodCondition(tc tests.TestContext, clientset *kubernetes.Clientset, namespace string, podName string, conditionType corev1.PodConditionType) {
	require.NoError(tc, WaitForPodCondition(tc.DefaultContext(), clientset, namespace, podName, conditionType))
}

// WaitForPodCondition watches the specified pod until the status includes the specified condition.
func WaitForPodCondition(ctx context.Context, clientset *kubernetes.Clientset, namespace string, podName string, conditionType corev1.PodConditionType) error {
	return waitForPodStatus(
		ctx,
		clientset,
		namespace,
		podName,
		func(status *corev1.PodStatus) bool {
			for _, condition := range status.Conditions {
				if condition.Type == conditionType && condition.Status == corev1.ConditionTrue {
					return true
				}
			}
			return false
		},
	)
}

// waitForPodStatus watches the specified pod until the status is deemed acceptable by the provided test function.
func waitForPodStatus(
	ctx context.Context,
	clientset *kubernetes.Clientset,
	namespace string,
	name string,
	acceptable func(*corev1.PodStatus) bool,
) error {
	watch, err := clientset.CoreV1().Pods(namespace).Watch(ctx, metav1.SingleObject(metav1.ObjectMeta{Name: name}))
	if err != nil {
		return fmt.Errorf("failed to initiate watch of pod %s/%s: %w", namespace, name, err)
	}

	for {
		select {
		case event := <-watch.ResultChan():
			pod, ok := event.Object.(*corev1.Pod)
			if !ok {
				continue
			}

			if acceptable(&pod.Status) {
				return nil
			}
		case <-ctx.Done():
			return errors.New("timeout waiting for pod readiness")
		}
	}
}

// waitForNodeHealthy waits for the node running in the specified pod to report healthy.
// TODO(marun) Use a structured logger
func WaitForNodeHealthy(
	ctx context.Context,
	log logging.Logger,
	kubeconfig *restclient.Config,
	namespace string,
	podName string,
	healthCheckInterval time.Duration,
	out io.Writer,
	outErr io.Writer,
) (ids.NodeID, error) {
	// A forwarded connection enables connectivity without exposing the node external to the kube cluster
	localPort, localPortStopChan, err := enableLocalForwardForPod(
		kubeconfig,
		namespace,
		podName,
		config.DefaultHTTPPort,
		out,
		outErr,
	)
	if err != nil {
		return ids.NodeID{}, fmt.Errorf("failed to enable local forward for pod: %w", err)
	}
	defer close(localPortStopChan)
	localNodeURI := fmt.Sprintf("http://127.0.0.1:%d", localPort)

	// TODO(marun) A node started with tmpnet should know the node ID before start
	infoClient := info.NewClient(localNodeURI)
	bootstrapNodeID, _, err := infoClient.GetNodeID(ctx)
	if err != nil {
		return ids.NodeID{}, fmt.Errorf("failed to retrieve node bootstrap ID: %w", err)
	}
	if err := wait.PollImmediateInfinite(healthCheckInterval, func() (bool, error) {
		healthReply, err := CheckNodeHealth(ctx, localNodeURI)
		if err != nil {
			log.Debug("Error checking node health", zap.Error(err))
			return false, nil
		}
		return healthReply.Healthy, nil

	}); err != nil {
		return ids.NodeID{}, fmt.Errorf("failed to wait for node to report healthy: %w", err)
	}

	return bootstrapNodeID, nil
}

// enableLocalForwardForPod enables traffic forwarding from a local port to the specified pod with client-go. The returned
// stop channel should be closed to stop the port forwarding.
func enableLocalForwardForPod(kubeconfig *restclient.Config, namespace string, name string, port int, out, errOut io.Writer) (uint16, chan struct{}, error) {
	transport, upgrader, err := spdy.RoundTripperFor(kubeconfig)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to create round tripper: %w", err)
	}

	dialer := spdy.NewDialer(
		upgrader,
		&http.Client{
			Transport: transport,
		},
		http.MethodPost,
		&url.URL{
			Scheme: "https",
			Path:   fmt.Sprintf("/api/v1/namespaces/%s/pods/%s/portforward", namespace, name),
			Host:   strings.TrimPrefix(kubeconfig.Host, "https://"),
		},
	)
	ports := []string{fmt.Sprintf("0:%d", port)}

	// Need to specify 127.0.0.1 to ensure that forwarding is only attempted for the ipv4
	// address of the pod. By default, kind is deployed with only ipv4, and attempting to
	// connect to a pod with ipv6 will fail.
	// TODO(marun) Is this needed?
	addresses := []string{"127.0.0.1"}

	stopChan, readyChan := make(chan struct{}, 1), make(chan struct{}, 1)
	forwarder, err := portforward.NewOnAddresses(dialer, addresses, ports, stopChan, readyChan, out, errOut)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to create forwarder: %w", err)
	}

	go func() {
		if err := forwarder.ForwardPorts(); err != nil {
			// TODO(marun) Need better error handling here? Or is ok for test-only usage?
			panic(err)
		}
	}()

	<-readyChan // Wait for port forwarding to be ready

	// Retrieve the dynamically allocated local port
	forwardedPorts, err := forwarder.GetPorts()
	if err != nil {
		close(stopChan)
		return 0, nil, fmt.Errorf("failed to get forwarded ports: %w", err)
	}
	if len(forwardedPorts) == 0 {
		close(stopChan)
		return 0, nil, fmt.Errorf("failed to find at least one forwarded port: %w", err)
	}
	return forwardedPorts[0].Local, stopChan, nil
}
