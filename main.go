package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Listen for Ctrl+C (SIGINT) or Kill (SIGTERM)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-stop // Wait for the user to hit Ctrl+C
		fmt.Println("\nShutting down gracefully...")
		cancel()
	}()

	if len(os.Args) > 1 {
		arg := os.Args[1]
		if err := kubePF(ctx, arg); err != nil {
			if ctx.Err() == context.Canceled {
				fmt.Println("Port-forward closed.")
			} else {
				fmt.Printf("Error port-forwarding: %v\n", err)
			}
		}
	} else {
		if err := devclonePF(ctx); err != nil {
			if ctx.Err() == context.Canceled {
				fmt.Println("Port-forward closed.")
			} else {
				fmt.Printf("Error port-forwarding: %v\n", err)
			}
		}
	}
}
