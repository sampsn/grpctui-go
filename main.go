package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
)

func main() {
	if len(os.Args) > 1 {
		arg := os.Args[1]
		// if err := kubePF(ctx, arg); err != nil {
		// 	if ctx.Err() == context.Canceled {
		// 		fmt.Println("Port-forward closed.")
		// 	} else {
		// 		fmt.Printf("Error port-forwarding: %v\n", err)
		// 	}
		// }

		if err := PortForwardPod(arg); err != nil {
			fmt.Printf("Error port-forwarding: %v\n", err)
		}
	} else {
		// if err := devclonePF(ctx); err != nil {
		// 	if ctx.Err() == context.Canceled {
		// 		fmt.Println("Port-forward closed.")
		// 	} else {
		// 		fmt.Printf("Error port-forwarding: %v\n", err)
		// 	}
		// }
		if err := PortForwardPod(""); err != nil {
			fmt.Printf("Error port-forwarding: %v\n", err)
		}
	}

	conn, stream, err := createReflectionClient()
	if err != nil {
		log.Fatalf("Failed to create reflection client connection: %v", err)
		return
	}
	defer conn.Close()

	getServices(stream)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	fmt.Println("Press Ctrl+c to stop port forwarding...")
	<-sigChan

	fmt.Println("Shutting down...")
}
