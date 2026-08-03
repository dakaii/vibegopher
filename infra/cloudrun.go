package main

import (
	"fmt"

	"github.com/pulumi/pulumi-gcp/sdk/v8/go/gcp/cloudrunv2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func deployCloudRun(ctx *pulumi.Context, cfg stackConfig, apis *enabledAPIs, sas *serviceAccounts, secrets *appSecrets) (*cloudrunv2.Service, error) {
	envFromSecrets := cloudrunv2.ServiceTemplateContainerEnvArray{}
	for _, name := range managedSecretNames {
		secret := secrets.Secrets[name]
		envFromSecrets = append(envFromSecrets, &cloudrunv2.ServiceTemplateContainerEnvArgs{
			Name: pulumi.String(name),
			ValueSource: &cloudrunv2.ServiceTemplateContainerEnvValueSourceArgs{
				SecretKeyRef: &cloudrunv2.ServiceTemplateContainerEnvValueSourceSecretKeyRefArgs{
					Secret:  secret.SecretId,
					Version: pulumi.String("latest"),
				},
			},
		})
	}

	// Non-secret runtime config stays as plain env (also suitable for Pulumi config).
	// Do not set PORT — Cloud Run injects it; setting it returns API error 400.
	plainEnv := cloudrunv2.ServiceTemplateContainerEnvArray{
		&cloudrunv2.ServiceTemplateContainerEnvArgs{
			Name:  pulumi.String("APP_ENV"),
			Value: pulumi.String("production"),
		},
		&cloudrunv2.ServiceTemplateContainerEnvArgs{
			Name:  pulumi.String("HASH_COST"),
			Value: pulumi.String("14"),
		},
		&cloudrunv2.ServiceTemplateContainerEnvArgs{
			Name:  pulumi.String("CORS_ORIGIN"),
			Value: pulumi.String(cfg.CorsOrigin),
		},
		// Poll bot_jobs in-process (min instances >= 1 for reliable workers).
		&cloudrunv2.ServiceTemplateContainerEnvArgs{
			Name:  pulumi.String("CRITIC_WORKER_ENABLED"),
			Value: pulumi.String("true"),
		},
	}

	service, err := cloudrunv2.NewService(ctx, "api-service", &cloudrunv2.ServiceArgs{
		Project:  pulumi.String(cfg.Project),
		Location: pulumi.String(cfg.Region),
		Name:     pulumi.String(cfg.ServiceName),
		Ingress:  pulumi.String("INGRESS_TRAFFIC_ALL"),
		Template: &cloudrunv2.ServiceTemplateArgs{
			ServiceAccount: sas.Runtime.Email,
			Containers: cloudrunv2.ServiceTemplateContainerArray{
				&cloudrunv2.ServiceTemplateContainerArgs{
					Image: pulumi.String(cfg.imageURL()),
					Envs:  append(plainEnv, envFromSecrets...),
					Ports: &cloudrunv2.ServiceTemplateContainerPortsArgs{
						ContainerPort: pulumi.Int(8080),
					},
					Resources: &cloudrunv2.ServiceTemplateContainerResourcesArgs{
						Limits: pulumi.StringMap{
							"cpu":    pulumi.String("1"),
							"memory": pulumi.String("512Mi"),
						},
					},
				},
			},
			Scaling: &cloudrunv2.ServiceTemplateScalingArgs{
				// Keep one instance warm so the in-process critic worker keeps polling.
				MinInstanceCount: pulumi.Int(1),
				MaxInstanceCount: pulumi.Int(5),
			},
		},
	}, pulumi.DependsOn([]pulumi.Resource{apis.Run, sas.Runtime}))
	if err != nil {
		return nil, fmt.Errorf("create cloud run service: %w", err)
	}

	_, err = cloudrunv2.NewServiceIamMember(ctx, "api-invoker-public", &cloudrunv2.ServiceIamMemberArgs{
		Project:  pulumi.String(cfg.Project),
		Location: pulumi.String(cfg.Region),
		Name:     service.Name,
		Role:     pulumi.String("roles/run.invoker"),
		Member:   pulumi.String("allUsers"),
	})
	if err != nil {
		return nil, fmt.Errorf("allow unauthenticated invoke: %w", err)
	}

	return service, nil
}
