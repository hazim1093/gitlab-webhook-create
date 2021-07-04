package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/xanzy/go-gitlab"
)

type Webhook struct {
	ProjectName            string `json:"projectName"`
	WebhookURL             string `json:"webhookURL"`
	PushEvents             bool   `json:"pushEvents"`
	PushEventsBranchFilter string `json:"pushEventsBranchFilter"`
	MergeRequestEvents     bool   `json:"mergeRequestEvents"`
}

func hookExists(client *gitlab.Client, pid int, webhookURL string) bool {
	projectHooks, _, err := client.Projects.ListProjectHooks(pid, nil)
	if err != nil {
		log.Fatalf("Failed to read projects hooks: %v", err)
	}

	for _, projectHook := range projectHooks {
		if projectHook.URL == webhookURL {
			return true
		}
	}

	return false
}

func main() {

	webhooksFilePath := os.Getenv("WEBHOOKS_FILE_PATH")
	gitlabToken := os.Getenv("GITLAB_TOKEN")
	secretToken := os.Getenv("WEBHOOK_SECRET_TOKEN")
	enableSslVerification := true

	webhooksFile, err := os.Open(webhooksFilePath)
	if err != nil {
		log.Fatalf("Failed to read projects file: %v", err)
	}

	byteValue, _ := ioutil.ReadAll(webhooksFile)

	var webhooks []Webhook

	json.Unmarshal(byteValue, &webhooks)

	for _, webhook := range webhooks {

		client, err := gitlab.NewClient(gitlabToken)
		if err != nil {
			log.Fatalf("Failed to create client: %v", err)
		}

		project, _, err := client.Projects.GetProject(webhook.ProjectName, nil)

		if err != nil {
			log.Fatalf("Failed to get Project: %v", err)
		}

		if hookExists(client, project.ID, webhook.WebhookURL) {
			log.Printf("Webhook already exists for project '%s' with ID '%d'", webhook.ProjectName, project.ID)
			continue
		}

		hook, _, err := client.Projects.AddProjectHook(project.ID, &gitlab.AddProjectHookOptions{
			URL:                    &webhook.WebhookURL,
			PushEvents:             &webhook.PushEvents,
			PushEventsBranchFilter: &webhook.PushEventsBranchFilter,
			MergeRequestsEvents:    &webhook.MergeRequestEvents,
			Token:                  &secretToken,
			EnableSSLVerification:  &enableSslVerification,
		})

		if err != nil {
			log.Fatalf("Failed to add Project hook: %v", err)
		}

		log.Printf("Created webhook: '%d' for project '%s' with ID '%d'", hook.ID, webhook.ProjectName, project.ID)
	}
}
