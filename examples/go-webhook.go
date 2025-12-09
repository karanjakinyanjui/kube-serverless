package main

import (
	"fmt"
	"time"
)

func Handler(event map[string]interface{}) (interface{}, error) {
	fmt.Printf("Webhook received at %s\n", time.Now().Format(time.RFC3339))

	// Process webhook payload
	body := event["body"]

	response := map[string]interface{}{
		"statusCode": 200,
		"body": map[string]interface{}{
			"message":    "Webhook processed successfully",
			"received":   body,
			"processedAt": time.Now().Format(time.RFC3339),
		},
	}

	return response, nil
}
