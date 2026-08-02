package main

import (
	"fmt"

	"github.com/pulumi/pulumi-gcp/sdk/v8/go/gcp/projects"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type enabledAPIs struct {
	SecretManager *projects.Service
	Run           *projects.Service
	ArtifactReg   *projects.Service
	IAM           *projects.Service
	IAMCreds      *projects.Service
	CloudResource *projects.Service
}

func enableAPIs(ctx *pulumi.Context, cfg stackConfig) (*enabledAPIs, error) {
	enable := func(name, service string, opts ...pulumi.ResourceOption) (*projects.Service, error) {
		return projects.NewService(ctx, name, &projects.ServiceArgs{
			Project:                  pulumi.String(cfg.Project),
			Service:                  pulumi.String(service),
			DisableOnDestroy:         pulumi.Bool(false),
			DisableDependentServices: pulumi.Bool(false),
		}, opts...)
	}

	// Secret Manager API stays enabled even when the rest of the stack is destroyed
	// (and is protected by default so destroy --exclude-protected will skip it).
	secretMgr, err := enable(
		"enable-secretmanager",
		"secretmanager.googleapis.com",
		cfg.secretResourceOpts()...,
	)
	if err != nil {
		return nil, fmt.Errorf("enable secretmanager: %w", err)
	}

	runAPI, err := enable("enable-run", "run.googleapis.com")
	if err != nil {
		return nil, fmt.Errorf("enable run: %w", err)
	}

	arAPI, err := enable("enable-artifactregistry", "artifactregistry.googleapis.com")
	if err != nil {
		return nil, fmt.Errorf("enable artifactregistry: %w", err)
	}

	iamAPI, err := enable("enable-iam", "iam.googleapis.com")
	if err != nil {
		return nil, fmt.Errorf("enable iam: %w", err)
	}

	iamCreds, err := enable("enable-iamcredentials", "iamcredentials.googleapis.com")
	if err != nil {
		return nil, fmt.Errorf("enable iamcredentials: %w", err)
	}

	crm, err := enable("enable-cloudresourcemanager", "cloudresourcemanager.googleapis.com")
	if err != nil {
		return nil, fmt.Errorf("enable cloudresourcemanager: %w", err)
	}

	return &enabledAPIs{
		SecretManager: secretMgr,
		Run:           runAPI,
		ArtifactReg:   arAPI,
		IAM:           iamAPI,
		IAMCreds:      iamCreds,
		CloudResource: crm,
	}, nil
}
