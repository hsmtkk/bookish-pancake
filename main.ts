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
      github: {
        owner: 'hsmtkk',
        name: repository,
        push: {
          branch: 'main',
        },
      },
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
