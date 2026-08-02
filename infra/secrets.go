package main

import (
	"fmt"

	"github.com/pulumi/pulumi-gcp/sdk/v8/go/gcp/secretmanager"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// managedSecretNames are Secret Manager secret IDs created by Pulumi.
// Pulumi owns the secret *containers* and IAM; payloads are synced from
// GitHub Actions (or gcloud) so values are not stored in Pulumi state.
var managedSecretNames = []string{
	"DATABASE_URL",             // Neon pooled connection string
	"AUTH_SECRET",              // App session/JWT signing secret (if used after Google Auth)
	"GOOGLE_OAUTH_CLIENT_ID",   // Audience for Google ID token verification
	"GOOGLE_OAUTH_CLIENT_SECRET", // Only needed if using auth-code exchange (optional)
	"GEMINI_API_KEY",           // AI bot provider key (omit if using Vertex ADC later)
}

type appSecrets struct {
	Secrets map[string]*secretmanager.Secret
}

func createSecrets(ctx *pulumi.Context, cfg stackConfig, apis *enabledAPIs) (*appSecrets, error) {
	out := &appSecrets{Secrets: make(map[string]*secretmanager.Secret, len(managedSecretNames))}
	opts := append([]pulumi.ResourceOption{
		pulumi.DependsOn([]pulumi.Resource{apis.SecretManager}),
	}, cfg.secretResourceOpts()...)

	for _, name := range managedSecretNames {
		secret, err := secretmanager.NewSecret(ctx, "secret-"+name, &secretmanager.SecretArgs{
			Project:  pulumi.String(cfg.Project),
			SecretId: pulumi.String(name),
			Replication: &secretmanager.SecretReplicationArgs{
				Auto: &secretmanager.SecretReplicationAutoArgs{},
			},
			Labels: pulumi.StringMap{
				"app":        pulumi.String("vibegopher"),
				"managed-by": pulumi.String("pulumi"),
			},
		}, opts...)
		if err != nil {
			return nil, fmt.Errorf("create secret %s: %w", name, err)
		}
		out.Secrets[name] = secret

		if cfg.CreatePlaceholderSecretVers {
			_, err := secretmanager.NewSecretVersion(ctx, "secret-version-placeholder-"+name, &secretmanager.SecretVersionArgs{
				Secret:     secret.ID(),
				SecretData: pulumi.String("REPLACE_ME"),
			}, opts...)
			if err != nil {
				return nil, fmt.Errorf("create placeholder version for %s: %w", name, err)
			}
		}
	}

	return out, nil
}

func grantSecretAccessor(ctx *pulumi.Context, cfg stackConfig, secret *secretmanager.Secret, member pulumi.StringInput, name string) error {
	_, err := secretmanager.NewSecretIamMember(ctx, name, &secretmanager.SecretIamMemberArgs{
		Project:  pulumi.String(cfg.Project),
		SecretId: secret.SecretId,
		Role:     pulumi.String("roles/secretmanager.secretAccessor"),
		Member:   member,
	}, cfg.secretResourceOpts()...)
	if err != nil {
		return fmt.Errorf("grant secret accessor %s: %w", name, err)
	}
	return nil
}

func grantSecretVersionManager(ctx *pulumi.Context, cfg stackConfig, secret *secretmanager.Secret, member pulumi.StringInput, name string) error {
	_, err := secretmanager.NewSecretIamMember(ctx, name, &secretmanager.SecretIamMemberArgs{
		Project:  pulumi.String(cfg.Project),
		SecretId: secret.SecretId,
		Role:     pulumi.String("roles/secretmanager.secretVersionManager"),
		Member:   member,
	}, cfg.secretResourceOpts()...)
	if err != nil {
		return fmt.Errorf("grant secret version manager %s: %w", name, err)
	}
	return nil
}
