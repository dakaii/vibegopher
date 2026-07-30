package main

import (
	"fmt"

	"github.com/pulumi/pulumi-gcp/sdk/v8/go/gcp/artifactregistry"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func createArtifactRegistry(ctx *pulumi.Context, cfg stackConfig, apis *enabledAPIs) (*artifactregistry.Repository, error) {
	repo, err := artifactregistry.NewRepository(ctx, "artifact-repo", &artifactregistry.RepositoryArgs{
		Project:      pulumi.String(cfg.Project),
		Location:     pulumi.String(cfg.Region),
		RepositoryId: pulumi.String(cfg.ArtifactRepoID),
		Format:       pulumi.String("DOCKER"),
		Description:  pulumi.String("VibeGopher container images"),
		Labels: pulumi.StringMap{
			"app":        pulumi.String("vibegopher"),
			"managed-by": pulumi.String("pulumi"),
		},
	}, pulumi.DependsOn([]pulumi.Resource{apis.ArtifactReg}))
	if err != nil {
		return nil, fmt.Errorf("create artifact registry: %w", err)
	}
	return repo, nil
}
