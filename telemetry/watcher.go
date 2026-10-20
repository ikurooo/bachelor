package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func startPodWatcher(ctx context.Context, client *http.Client, baseURL string, notifier *PodNotifier) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		watchURL := fmt.Sprintf("%s/api/v1/namespaces/default/pods?watch=true", baseURL)
		req, err := http.NewRequestWithContext(ctx, "GET", watchURL, nil)
		if err != nil {
			time.Sleep(2 * time.Second)
			continue
		}

		resp, err := client.Do(req)
		if err != nil {
			log.Printf("Watch connection failed, reconnecting in 2s: %v", err)
			time.Sleep(2 * time.Second)
			continue
		}

		// Kubernetes watch streams are chunked JSON lines
		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			line := scanner.Bytes()

			var event struct {
				Type   string   `json:"type"`
				Object PodItems `json:"object"`
			}

			if err := json.Unmarshal(line, &event); err != nil {
				continue
			}

			// We only care about modified/added pods that are running and ready
			pod := event.Object
			if pod.Status.Phase == "Running" && pod.Status.PodIP != "" {
				isReady := false
				for _, cond := range pod.Status.Conditions {
					if cond.Type == "Ready" && cond.Status == "True" {
						isReady = true
						break
					}
				}

				if isReady {
					notifier.Publish(pod.Metadata.Name, pod.Status.PodIP)
				}
			}
		}
		resp.Body.Close()
		time.Sleep(1 * time.Second)
	}
}
