package main

import (
	"fmt"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		cfg, err := loadConfig(ctx)
		if err != nil {
			return err
		}

		apis, err := enableAPIs(ctx, cfg)
		if err != nil {
			return err
		}

		secrets, err := createSecrets(ctx, cfg, apis)
		if err != nil {
			return err
		}

		repo, err := createArtifactRegistry(ctx, cfg, apis)
		if err != nil {
			return err
		}

		sas, err := createServiceAccounts(ctx, cfg, apis)
		if err != nil {
			return err
		}

		if err := grantRuntimeSecretAccess(ctx, cfg, sas, secrets); err != nil {
			return err
		}
		if err := grantDeployPermissions(ctx, cfg, sas, repo, secrets); err != nil {
			return err
		}

		wif, err := setupGitHubWorkloadIdentity(ctx, cfg, apis, sas)
		if err != nil {
			return err
		}

		ctx.Export("project", pulumi.String(cfg.Project))
		ctx.Export("region", pulumi.String(cfg.Region))
		ctx.Export("protectSecrets", pulumi.Bool(cfg.ProtectSecrets))
		ctx.Export("artifactRegistryRepo", repo.Name)
		ctx.Export("imageUrl", pulumi.String(cfg.imageURL()))
		ctx.Export("runtimeServiceAccount", sas.Runtime.Email)
		ctx.Export("deployServiceAccount", sas.Deploy.Email)
		ctx.Export("secretIds", pulumi.ToStringArray(managedSecretNames))

		if cfg.EnableCloudRun {
			service, err := deployCloudRun(ctx, cfg, apis, sas, secrets)
			if err != nil {
				return err
			}
			ctx.Export("cloudRunService", service.Name)
			ctx.Export("cloudRunUri", service.Uri)
		} else {
			ctx.Log.Info("vibegopher:enableCloudRun=false — skipping Cloud Run (useful for first bootstrap before an image exists)", nil)
			ctx.Export("cloudRunService", pulumi.String(""))
			ctx.Export("cloudRunUri", pulumi.String(""))
		}

		if wif != nil {
			ctx.Export("workloadIdentityProvider", wif.Provider.Name.ApplyT(func(name string) string {
				return name
			}).(pulumi.StringOutput))
			ctx.Export("githubActionsAuthHint", pulumi.All(wif.Provider.Name, sas.Deploy.Email).ApplyT(
				func(args []interface{}) string {
					return fmt.Sprintf(
						"workload_identity_provider=%s service_account=%s",
						args[0], args[1],
					)
				},
			).(pulumi.StringOutput))
		}

		return nil
	})
}
