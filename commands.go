package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func getPods(service string) error {
	cmdStr := fmt.Sprintf("kubectl get pods | grep org-%s", service)
	cmd := exec.Command("/bin/sh", "-c", cmdStr)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// kubePF portforwards to service that is provided
func kubePF(ctx context.Context, service string) error {
	// Build the shell command string
	searchCmd := fmt.Sprintf("kubectl get pods -l app=org-%s | tail -n 1", service)

	// Run the command and capture output
	out, err := exec.Command("/bin/sh", "-c", searchCmd).Output()
	if err != nil {
		fmt.Printf("Failed to get pod-name: %v\n", err)
		return err
	}

	// Handle empty results
	line := string(out)
	if line == "" {
		fmt.Println("No pods found.")
		return nil
	}

	// Get the first field
	parts := strings.Fields(line)
	podName := parts[0]

	// Start the port forward
	cmd := exec.CommandContext(ctx, "kubectl", "port-forward", podName, "50052:50051")

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// print
	fmt.Printf("Starting port-forward to %s... (Press Ctrl+C to stop)\n", podName)

	return cmd.Run()
}

// devclonePF portforwards to service that is provided
func devclonePF(ctx context.Context) error {
	searchCmd := "kubectl get pods | grep sampson | head -n 1"
	out, err := exec.Command("/bin/sh", "-c", searchCmd).Output()
	if err != nil {
		fmt.Printf("Failed to get pod-name: %v\n", err)
		return err
	}

	line := string(out)
	if line == "" {
		fmt.Println("No pods found.")
		return nil
	}

	parts := strings.Fields(line)
	podName := parts[0]

	cmd := exec.CommandContext(ctx, "kubectl", "port-forward", podName, "50052:50051")

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Printf("Starting port-forward to devclone: %s... (Press Ctrl+C to stop)\n", podName)

	return cmd.Run()
}

