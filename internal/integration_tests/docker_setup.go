package integration_tests

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go/modules/compose"
	"github.com/testcontainers/testcontainers-go/wait"
)

var composeInstance compose.ComposeStack

func DockerComposeUp(t *testing.T) {
	ctx := context.Background()

	composeFilePaths := []string{"../../docker-compose.yml"}
	var err error

	composeInstance, err = compose.NewDockerCompose(composeFilePaths...)
	if err != nil {
		t.Fatalf("Failed to create Docker Compose instance: %v", err)
	}

	startDockerCompose(t, ctx)
	setEnvironmentVariables(t)
	waitForServices(t, ctx)
	log.Printf("All environment services are up and running")
}

func DockerComposeDown(t *testing.T) {
	ctx := context.Background()

	if composeInstance != nil {
		if err := composeInstance.Down(ctx); err != nil {
			t.Fatalf("Failed to stop Docker Compose: %v", err)
		}
	}
}

func startDockerCompose(t *testing.T, ctx context.Context) {
	if err := composeInstance.Up(ctx); err != nil {
		t.Fatalf("Failed to start Docker Compose: %v", err)
	}
}

func setEnvironmentVariables(t *testing.T) {
	handleError(t, os.Setenv("DB_HOST", "localhost"))
	handleError(t, os.Setenv("DB_PORT", "5432"))
	handleError(t, os.Setenv("DB_USER", "root"))
	handleError(t, os.Setenv("DB_PASSWORD", "password"))
	handleError(t, os.Setenv("DB_NAME", "currency_notifier"))
	handleError(t, os.Setenv("SMTP_HOST", "localhost"))
	handleError(t, os.Setenv("SMTP_PORT", "1025"))
	handleError(t, os.Setenv("MONOBANK_HOST_URL", "http://localhost:8282"))
}

func waitForServices(t *testing.T, ctx context.Context) {
	_, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()

	// Waiting for services to be healthy
	handleError(t, composeInstance.WaitForService("db", wait.ForHealthCheck()).Up(ctx, compose.Wait(true)))
	log.Printf("Database is up and running")
	handleError(t, composeInstance.WaitForService("wiremock", wait.ForHealthCheck()).Up(ctx, compose.Wait(true)))
	log.Printf("WireMock is up and running")
}

func handleError(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
}
