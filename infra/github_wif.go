package main

import (
	"fmt"

	"github.com/pulumi/pulumi-gcp/sdk/v8/go/gcp/iam"
	"github.com/pulumi/pulumi-gcp/sdk/v8/go/gcp/serviceaccount"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type workloadIdentity struct {
	Pool     *iam.WorkloadIdentityPool
	Provider *iam.WorkloadIdentityPoolProvider
}

// setupGitHubWorkloadIdentity lets GitHub Actions impersonate the deploy SA
// without storing a long-lived JSON key. Requires vibegopher:githubOwner.
func setupGitHubWorkloadIdentity(ctx *pulumi.Context, cfg stackConfig, apis *enabledAPIs, sas *serviceAccounts) (*workloadIdentity, error) {
	if cfg.GitHubOwner == "" {
		ctx.Log.Info("vibegopher:githubOwner not set — skipping GitHub Workload Identity Federation (set it before using deploy.yml)", nil)
		return nil, nil
	}

	// Use a vibegopher-specific pool id so we do not collide with other apps
	// in a shared GCP project that already own a generic "github-actions" pool.
	pool, err := iam.NewWorkloadIdentityPool(ctx, "github-pool", &iam.WorkloadIdentityPoolArgs{
		Project:                pulumi.String(cfg.Project),
		WorkloadIdentityPoolId: pulumi.String("vibegopher-github"),
		DisplayName:            pulumi.String("VibeGopher GitHub Actions"),
		Description:            pulumi.String("OIDC identity pool for vibegopher GitHub Actions deployments"),
	}, pulumi.DependsOn([]pulumi.Resource{apis.IAM, apis.IAMCreds}))
	if err != nil {
		return nil, fmt.Errorf("create WIF pool: %w", err)
	}

	attributeCondition := fmt.Sprintf(
		`assertion.repository_owner == '%s' && assertion.repository == '%s/%s'`,
		cfg.GitHubOwner, cfg.GitHubOwner, cfg.GitHubRepo,
	)

	provider, err := iam.NewWorkloadIdentityPoolProvider(ctx, "github-provider", &iam.WorkloadIdentityPoolProviderArgs{
		Project:                        pulumi.String(cfg.Project),
		WorkloadIdentityPoolId:         pool.WorkloadIdentityPoolId,
		WorkloadIdentityPoolProviderId: pulumi.String("github"),
		DisplayName:                    pulumi.String("GitHub OIDC"),
		AttributeCondition:             pulumi.String(attributeCondition),
		AttributeMapping: pulumi.StringMap{
			"google.subject":       pulumi.String("assertion.sub"),
			"attribute.actor":      pulumi.String("assertion.actor"),
			"attribute.repository": pulumi.String("assertion.repository"),
			"attribute.ref":        pulumi.String("assertion.ref"),
		},
		Oidc: &iam.WorkloadIdentityPoolProviderOidcArgs{
			IssuerUri: pulumi.String("https://token.actions.githubusercontent.com"),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("create WIF provider: %w", err)
	}

	member := pulumi.All(pool.Name, pulumi.String(cfg.GitHubOwner), pulumi.String(cfg.GitHubRepo)).ApplyT(
		func(args []interface{}) string {
			poolName := args[0].(string)
			owner := args[1].(string)
			repo := args[2].(string)
			return fmt.Sprintf(
				"principalSet://iam.googleapis.com/%s/attribute.repository/%s/%s",
				poolName, owner, repo,
			)
		},
	).(pulumi.StringOutput)

	_, err = serviceaccount.NewIAMMember(ctx, "deploy-wif-user", &serviceaccount.IAMMemberArgs{
		ServiceAccountId: sas.Deploy.Name,
		Role:             pulumi.String("roles/iam.workloadIdentityUser"),
		Member:           member,
	})
	if err != nil {
		return nil, fmt.Errorf("bind WIF to deploy SA: %w", err)
	}

	return &workloadIdentity{Pool: pool, Provider: provider}, nil
}
