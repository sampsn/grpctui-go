package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/transport/spdy"
)

const (
	ns      = "default"
	devName = "gabriel-sampson"
)

// Standard boilerplate to get a clientset
var (
	kubeconfig   = os.Getenv("KUBECONFIG")
	config, _    = clientcmd.BuildConfigFromFlags("", kubeconfig)
	clientset, _ = kubernetes.NewForConfig(config)
	ports        = []string{"50052:50051"}
)

func PortForwardPod(app string) error {
	transport, upgrader, err := spdy.RoundTripperFor(config)
	if err != nil {
		return err
	}

	stopChan := make(chan struct{}, 1)
	readyChan := make(chan struct{})
	sigChan := make(chan os.Signal, 1)

	signal.Notify(sigChan, os.Interrupt)

	go func() {
		<-sigChan
		close(stopChan)
	}()

	// Get pod name
	var podName string
	if app != "" {
		podName, err = getMainPodName(app)
		if err != nil {
			return err
		}
	} else {
		podName, err = getDevclonePodName()
		if err != nil {
			return err
		}
	}

	path := fmt.Sprintf("/api/v1/namespaces/%s/pods/%s/portforward", ns, podName)

	hostURL, err := url.Parse(config.Host)
	if err != nil {
		return err
	}

	hostURL.Path = path

	dialer := spdy.NewDialer(upgrader, &http.Client{Transport: transport}, "POST", hostURL)

	out := os.Stdout
	errOut := os.Stderr

	forwarder, err := portforward.New(dialer, ports, stopChan, readyChan, out, errOut)
	if err != nil {
		return err
	}

	fmt.Printf("Starting port forward to %s...\n", podName)

	go func() {
		if err := forwarder.ForwardPorts(); err != nil {
			fmt.Println("port-forward error: ", err)
		}
	}()

	// wait until ready
	<-readyChan

	fmt.Println("Port forwarding ready on localhost:50052")

	return nil
}

func getMainPodName(app string) (string, error) {
	appName := fmt.Sprintf("org-%s", app)
	listOptions := metav1.ListOptions{
		LabelSelector: fmt.Sprintf("app=%s", appName),
	}

	pods, _ := clientset.CoreV1().Pods(ns).List(context.TODO(), listOptions)

	for _, pod := range pods.Items {
		// Only pick pods that start exactly with "org-users"
		// This avoids things like "dev-org-users-xyz"
		if strings.HasPrefix(pod.Name, appName) && pod.Status.Phase == "Running" {
			return pod.Name, nil
		}
	}

	return "", fmt.Errorf("no official running pod found")
}

func getDevclonePodName() (string, error) {
	devcloneName := fmt.Sprintf("dev-%s", devName)

	pods, _ := clientset.CoreV1().Pods(ns).List(context.TODO(), metav1.ListOptions{})

	for _, pod := range pods.Items {
		// Only pick pods that start exactly with "org-users"
		// This avoids things like "dev-org-users-xyz"
		if strings.HasPrefix(pod.Name, devcloneName) && pod.Status.Phase == "Running" {
			return pod.Name, nil
		}
	}

	return "", fmt.Errorf("no official running pod found")
}
