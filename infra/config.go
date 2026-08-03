package main

import (
	"fmt"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

// stackConfig holds non-secret infrastructure settings from Pulumi config.
// Secret payloads do not live here — see SECRETS.md.
type stackConfig struct {
	Project                     string
	Region                      string
	ServiceName                 string
	ArtifactRepoID              string
	ImageName                   string
	ImageTag                    string
	GitHubOwner                 string
	GitHubRepo                  string
	CorsOrigin                  string
	ProtectSecrets              bool
	CreatePlaceholderSecretVers bool
	EnableCloudRun              bool
	CloudRunServiceAccountID    string
	DeployServiceAccountID      string
}

func loadConfig(ctx *pulumi.Context) (stackConfig, error) {
	cfg := config.New(ctx, "vibegopher")
	gcpCfg := config.New(ctx, "gcp")

	project := gcpCfg.Get("project")
	if project == "" {
		return stackConfig{}, fmt.Errorf("gcp:project is required (set with: pulumi config set gcp:project PROJECT_ID)")
	}

	// Defaults match Pulumi.yaml; only override when the key is present.
	protectSecrets := true
	if cfg.Get("protectSecrets") != "" {
		protectSecrets = cfg.GetBool("protectSecrets")
	}

	createPlaceholders := true
	if cfg.Get("createPlaceholderSecretVersions") != "" {
		createPlaceholders = cfg.GetBool("createPlaceholderSecretVersions")
	}

	enableCloudRun := true
	if cfg.Get("enableCloudRun") != "" {
		enableCloudRun = cfg.GetBool("enableCloudRun")
	}

	region := cfg.Get("region")
	if region == "" {
		region = "us-central1"
	}
	serviceName := cfg.Get("serviceName")
	if serviceName == "" {
		serviceName = "vibegopher-api"
	}
	artifactRepoID := cfg.Get("artifactRepoId")
	if artifactRepoID == "" {
		artifactRepoID = "vibegopher"
	}
	imageName := cfg.Get("imageName")
	if imageName == "" {
		imageName = "api"
	}
	imageTag := cfg.Get("imageTag")
	if imageTag == "" {
		imageTag = "latest"
	}
	githubRepo := cfg.Get("githubRepo")
	if githubRepo == "" {
		githubRepo = "vibegopher"
	}

	corsOrigin := cfg.Get("corsOrigin")
	if enableCloudRun && corsOrigin == "" {
		return stackConfig{}, fmt.Errorf("vibegopher:corsOrigin is required when Cloud Run is enabled (frontend origin, e.g. https://app.example.com or http://localhost:5173)")
	}

	return stackConfig{
		Project:                     project,
		Region:                      region,
		ServiceName:                 serviceName,
		ArtifactRepoID:              artifactRepoID,
		ImageName:                   imageName,
		ImageTag:                    imageTag,
		GitHubOwner:                 cfg.Get("githubOwner"),
		GitHubRepo:                  githubRepo,
		CorsOrigin:                  corsOrigin,
		ProtectSecrets:              protectSecrets,
		CreatePlaceholderSecretVers: createPlaceholders,
		EnableCloudRun:              enableCloudRun,
		CloudRunServiceAccountID:    "vibegopher-run",
		DeployServiceAccountID:      "vibegopher-deploy",
	}, nil
}

func (c stackConfig) imageURL() string {
	return fmt.Sprintf("%s-docker.pkg.dev/%s/%s/%s:%s", c.Region, c.Project, c.ArtifactRepoID, c.ImageName, c.ImageTag)
}

// secretResourceOpts returns Pulumi options that keep Secret Manager resources
// from being destroyed when protectSecrets is enabled (default).
func (c stackConfig) secretResourceOpts() []pulumi.ResourceOption {
	if !c.ProtectSecrets {
		return nil
	}
	return []pulumi.ResourceOption{
		pulumi.Protect(true),
		pulumi.RetainOnDelete(true),
	}
}
