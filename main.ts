// Copyright (c) HashiCorp, Inc
// SPDX-License-Identifier: MPL-2.0
import { Construct } from "constructs";
import { App, TerraformStack, CloudBackend, NamedCloudWorkspace } from "cdktf";
import * as google from '@cdktf/provider-google';

const region = 'asia-northeast1';
const project = 'bookish-pancake-369300';
const repository = 'bookish-pancake';

class MyStack extends TerraformStack {
  constructor(scope: Construct, id: string) {
    super(scope, id);

    new google.provider.GoogleProvider(this, 'google', {
      region,
      project,
    });

    new google.artifactRegistryRepository.ArtifactRegistryRepository(this, 'artifact_registry', {
      format: 'docker',
      location: region,
      repositoryId: repository,
    });

    new google.cloudbuildTrigger.CloudbuildTrigger(this, 'build_trigger', {
      filename: 'cloudbuild.yaml',
      github: {
        owner: 'hsmtkk',
        name: repository,
        push: {
          branch: 'main',
        },
      },
    });

    const openweather_api_key = new google.secretManagerSecret.SecretManagerSecret(this, 'openweather_api_key', {
      secretId: 'openweather_api_key',
      replication: {
        automatic: true,
      },
    });

    const cloud_run_service_account = new google.serviceAccount.ServiceAccount(this, 'cloud_run_service_account', {
      accountId: 'cloud-run-service-account',
      displayName: 'service account for Cloud Run',
    });

    new google.secretManagerSecretIamBinding.SecretManagerSecretIamBinding(this, 'from_cloudrun_to_secretmanager', {
      secretId: openweather_api_key.id,
      members: [`serviceAccount:${cloud_run_service_account.email}`],
      role: 'roles/secretmanager.secretAccessor',
    });

    new google.projectIamBinding.ProjectIamBinding(this, 'from_cloudrun_to_cloudtrace', {
      members: [`serviceAccount:${cloud_run_service_account.email}`],
      project,
      role: 'roles/cloudtrace.agent',
    });

    new google.projectIamBinding.ProjectIamBinding(this, 'from_cloudrun_to_cloudprofiler', {
      members: [`serviceAccount:${cloud_run_service_account.email}`],
      project,
      role: 'roles/cloudprofiler.agent',
    });


    const no_auth_policy = new google.dataGoogleIamPolicy.DataGoogleIamPolicy(this, 'cloud_run_no_auth_policy', {
      binding: [{
        role: 'roles/run.invoker',
        members: ['allUsers'],
      }],
    });

    
    const v1_backgrpc = new google.cloudRunService.CloudRunService(this, 'v1_backgrpc', {
      autogenerateRevisionName: true,
      location: region,
      name: 'v1-backgrpc',
      template: {
        spec: {
          containers: [{
            image: 'us-docker.pkg.dev/cloudrun/container/hello', // Placeholder. To be overwritten by Cloud Build.
          }],
          serviceAccountName: cloud_run_service_account.email,
        },
      }
    });

    new google.cloudRunServiceIamPolicy.CloudRunServiceIamPolicy(this, 'v1_backgrpc_policy', {
      location: region,
      service: v1_backgrpc.name,
      policyData: no_auth_policy.policyData,
    });

    const v1_frontweb = new google.cloudRunService.CloudRunService(this, 'v1_frontweb', {
      autogenerateRevisionName: true,
      location: region,
      name: 'v1-frontweb',
      template: {
        spec: {
          containers: [{
            env: [{
              name: 'BACK_URL',
              value: v1_backgrpc.status.get(0).url,
            }],
            image: 'us-docker.pkg.dev/cloudrun/container/hello',
          }],
          serviceAccountName: cloud_run_service_account.email,
        },
      }
    });

    new google.cloudRunServiceIamPolicy.CloudRunServiceIamPolicy(this, 'v1_frontweb_policy', {
      location: region,
      service: v1_frontweb.name,
      policyData: no_auth_policy.policyData,
    });

    const v2_backgrpc = new google.cloudRunService.CloudRunService(this, 'v2_backgrpc', {
      autogenerateRevisionName: true,
      location: region,
      name: 'v2-backgrpc',
      template: {
        spec: {
          containers: [{
            image: 'us-docker.pkg.dev/cloudrun/container/hello',
          }],
          serviceAccountName: cloud_run_service_account.email,
        },
      }
    });

    new google.cloudRunServiceIamPolicy.CloudRunServiceIamPolicy(this, 'v2_backgrpc_policy', {
      location: region,
      service: v2_backgrpc.name,
      policyData: no_auth_policy.policyData,
    });

    const v2_frontweb = new google.cloudRunService.CloudRunService(this, 'v2_frontweb', {
      autogenerateRevisionName: true,
      location: region,
      name: 'v2-frontweb',
      template: {
        spec: {
          containers: [{
            env: [{
              name: 'BACK_URL',
              value: v2_backgrpc.status.get(0).url,
            }],
            image: 'us-docker.pkg.dev/cloudrun/container/hello',
          }],
          serviceAccountName: cloud_run_service_account.email,
        },
      }
    });

    new google.cloudRunServiceIamPolicy.CloudRunServiceIamPolicy(this, 'v2_frontweb_policy', {
      location: region,
      service: v2_frontweb.name,
      policyData: no_auth_policy.policyData,
    });
  }
}

const app = new App();
const stack = new MyStack(app, "bookish-pancake");
new CloudBackend(stack, {
  hostname: "app.terraform.io",
  organization: "hsmtkkdefault",
  workspaces: new NamedCloudWorkspace("bookish-pancake")
});
app.synth();
