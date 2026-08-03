package main

import (
	"fmt"

	"github.com/pulumi/pulumi-gcp/sdk/v8/go/gcp/artifactregistry"
	"github.com/pulumi/pulumi-gcp/sdk/v8/go/gcp/organizations"
	"github.com/pulumi/pulumi-gcp/sdk/v8/go/gcp/projects"
	"github.com/pulumi/pulumi-gcp/sdk/v8/go/gcp/serviceaccount"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type serviceAccounts struct {
	Runtime *serviceaccount.Account
	Deploy  *serviceaccount.Account
}

func createServiceAccounts(ctx *pulumi.Context, cfg stackConfig, apis *enabledAPIs) (*serviceAccounts, error) {
	runtime, err := serviceaccount.NewAccount(ctx, "run-sa", &serviceaccount.AccountArgs{
		Project:     pulumi.String(cfg.Project),
		AccountId:   pulumi.String(cfg.CloudRunServiceAccountID),
		DisplayName: pulumi.String("VibeGopher Cloud Run runtime"),
		Description: pulumi.String("Runtime identity for the API; can read app secrets"),
	}, pulumi.DependsOn([]pulumi.Resource{apis.IAM}))
	if err != nil {
		return nil, fmt.Errorf("create runtime SA: %w", err)
	}

	deploy, err := serviceaccount.NewAccount(ctx, "deploy-sa", &serviceaccount.AccountArgs{
		Project:     pulumi.String(cfg.Project),
		AccountId:   pulumi.String(cfg.DeployServiceAccountID),
		DisplayName: pulumi.String("VibeGopher GitHub Actions deploy"),
		Description: pulumi.String("Used by GitHub Actions via Workload Identity Federation"),
	}, pulumi.DependsOn([]pulumi.Resource{apis.IAM}))
	if err != nil {
		return nil, fmt.Errorf("create deploy SA: %w", err)
	}

	return &serviceAccounts{Runtime: runtime, Deploy: deploy}, nil
}

func grantDeployPermissions(ctx *pulumi.Context, cfg stackConfig, sas *serviceAccounts, repo *artifactregistry.Repository, secrets *appSecrets) error {
	member := sas.Deploy.Email.ApplyT(func(email string) string {
		return "serviceAccount:" + email
	}).(pulumi.StringOutput)

	// Broad enough for GitHub Actions to run `pulumi up` / `pulumi destroy` on this stack.
	// projectIamAdmin is required so CI can remove project-level IAMMember bindings on destroy
	// (without it, destroy fails with 403 getIamPolicy after other resources are gone).
	// Secret *payloads* are still written via secretVersionManager + sync-secrets.sh.
	roles := []struct {
		name string
		role string
	}{
		{"deploy-run-admin", "roles/run.admin"},
		{"deploy-sa-user", "roles/iam.serviceAccountUser"},
		{"deploy-sa-admin", "roles/iam.serviceAccountAdmin"},
		{"deploy-ar-admin", "roles/artifactregistry.admin"},
		{"deploy-sm-admin", "roles/secretmanager.admin"},
		{"deploy-wif-admin", "roles/iam.workloadIdentityPoolAdmin"},
		{"deploy-serviceusage", "roles/serviceusage.serviceUsageAdmin"},
		{"deploy-project-iam", "roles/resourcemanager.projectIamAdmin"},
	}

	for _, r := range roles {
		_, err := projects.NewIAMMember(ctx, r.name, &projects.IAMMemberArgs{
			Project: pulumi.String(cfg.Project),
			Role:    pulumi.String(r.role),
			Member:  member,
		})
		if err != nil {
			return fmt.Errorf("grant %s: %w", r.name, err)
		}
	}

	// Narrower secret version manager on each secret (defense in depth alongside project admin).
	for name, secret := range secrets.Secrets {
		if err := grantSecretVersionManager(ctx, cfg, secret, member, "deploy-secret-versions-"+name); err != nil {
			return err
		}
	}

	_ = repo // repository exists for image push; project-level AR admin covers it
	return nil
}

func grantRuntimeSecretAccess(ctx *pulumi.Context, cfg stackConfig, sas *serviceAccounts, secrets *appSecrets) error {
	member := sas.Runtime.Email.ApplyT(func(email string) string {
		return "serviceAccount:" + email
	}).(pulumi.StringOutput)

	for name, secret := range secrets.Secrets {
		if err := grantSecretAccessor(ctx, cfg, secret, member, "run-secret-accessor-"+name); err != nil {
			return err
		}
	}
	return nil
}

// grantArtifactRegistryReaders lets Cloud Run pull private images from the repo.
// Cloud Run uses the Google-managed service agent for pulls; the runtime SA is granted too.
func grantArtifactRegistryReaders(ctx *pulumi.Context, cfg stackConfig, sas *serviceAccounts, repo *artifactregistry.Repository) error {
	proj, err := organizations.LookupProject(ctx, &organizations.LookupProjectArgs{
		ProjectId: &cfg.Project,
	}, nil)
	if err != nil {
		return fmt.Errorf("lookup project number for Cloud Run AR access: %w", err)
	}

	runAgent := fmt.Sprintf(
		"serviceAccount:service-%s@serverless-robot-prod.iam.gserviceaccount.com",
		proj.Number,
	)

	_, err = artifactregistry.NewRepositoryIamMember(ctx, "ar-reader-cloudrun-agent", &artifactregistry.RepositoryIamMemberArgs{
		Project:    pulumi.String(cfg.Project),
		Location:   pulumi.String(cfg.Region),
		Repository: repo.RepositoryId,
		Role:       pulumi.String("roles/artifactregistry.reader"),
		Member:     pulumi.String(runAgent),
	})
	if err != nil {
		return fmt.Errorf("grant AR reader to Cloud Run service agent: %w", err)
	}

	runtimeMember := sas.Runtime.Email.ApplyT(func(email string) string {
		return "serviceAccount:" + email
	}).(pulumi.StringOutput)

	_, err = artifactregistry.NewRepositoryIamMember(ctx, "ar-reader-runtime-sa", &artifactregistry.RepositoryIamMemberArgs{
		Project:    pulumi.String(cfg.Project),
		Location:   pulumi.String(cfg.Region),
		Repository: repo.RepositoryId,
		Role:       pulumi.String("roles/artifactregistry.reader"),
		Member:     runtimeMember,
	})
	if err != nil {
		return fmt.Errorf("grant AR reader to runtime SA: %w", err)
	}
	return nil
}
