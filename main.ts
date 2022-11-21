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

    new google.secretManagerSecret.SecretManagerSecret(this, 'openweather_api_key', {
      secretId: 'openweather_api_key',
      replication: {
        automatic: true,
      },
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

    const no_auth_policy = new google.dataGoogleIamPolicy.DataGoogleIamPolicy(this, 'no_auth_policy', {
      binding: [{
        role: 'roles/run.invoker',
        members: ['allUsers'],
      }],
    });

    const v1_backgrpc = new google.cloudRunService.CloudRunService(this, 'v1_backgrpc', {
      location: region,
      name: 'v1-backgrpc',
      template: {
        spec: {
          containers: [{
            image: 'asia-northeast1-docker.pkg.dev/bookish-pancake-369300/bookish-pancake/v1/backgrpc:latest',
          }],
        },
      }
    });

    new google.cloudRunServiceIamPolicy.CloudRunServiceIamPolicy(this, 'v1_backgrpc_policy', {
      service: v1_backgrpc.name,
      policyData: no_auth_policy.policyData,
    });

    const v1_frontweb = new google.cloudRunService.CloudRunService(this, 'v1_frontweb', {
      location: region,
      name: 'v1-frontweb',
      template: {
        spec: {
          containers: [{
            env: [{
              name: 'BACK_URL',
              value: v1_backgrpc.status.get(0).url,
            }],
            image: 'asia-northeast1-docker.pkg.dev/bookish-pancake-369300/bookish-pancake/v1/frontweb:latest',
          }],
        },
      }
    });

    new google.cloudRunServiceIamPolicy.CloudRunServiceIamPolicy(this, 'v1_frontweb_policy', {
      service: v1_frontweb.name,
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
