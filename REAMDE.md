# General
## Overview
## Running the node
### Create the crypto keys
- Using a Hedera Operator Account:
    - set follow environment variables - `HCS_OPERATOR_ID` and `HCS_OPERATOR_PRV_KEY`;
    Example: `HCS_OPERATOR_ID=0.0.7506 HCS_OPERATOR_PRV_KEY=... ./cryptogen.sh`
- Using own PEM encoded private key:
    - place a PEM encoded file in `config/key.pem`;

If you want to have another keypair, just rename the file and run cryptogen again.

### Creating HCS Topic with threshold keys using Hedera Operator Account
Example `HCS_OPERATOR_ID=0.0.7506 HCS_OPERATOR_PRV_KEY=... A_PUB_KEY=302a300506032b6570032100086e579c72b037e72bddc3d5c8af5e7b5e5269ab6bb025792a480940ba501b16 B_PUB_KEY=302a300506032b65700321005685758381e67fdaf28b6a992e1e725a707e95280d89eae47fc5132271dc2b1a ./createhcstopic.sh`

### Run your mongo database
`docker run --name hedera-mongo -d -p 27017:27017 -v ~/data:/data/db mongo`

### Create .env file
Create file named `.env` based on the `.env.example`

### Starting the node
Run `./start.sh`

## Running the second peer
1. Backup your `key.pem` and run `./cryptogen`. Rename the new `key.pem` to another name and restore the previous one.
2. Run mongo with different name port and data storage path
3. Create `peer2.env` based on `.env.example`
4. Run `./start-peer2.sh`

### Deployment process - Using Terraform to create HCS nodes in Google Cloud

We will be deploying two Compute Engine VM instances with a started and connected HCS node on each of them.  

**Before we begin, have the following tools locally:**

- gcloud;
- [Terraform](https://learn.hashicorp.com/terraform/getting-started/install.html);
- [jq](https://stedolan.github.io/jq/download/);

**Create a Google Cloud project**

We will start by creating a new project to keep this separate and easy to tear down later. After creating it, copy down the project ID and replace the default project ID with yours in the following places:

`./Makefile`:

```jsx
GCP_PROJECT ?= hcs-integration-node
```

`.terraform/variables.tf`:

```jsx
variable "project_id" {
  default = "hcs-integration-node"
}
```

**Prepare Google Cloud SDK authentication:**

```jsx
$ gcloud auth login
$ gcloud auth application-default login
```

**Switch to newly created project:**

```jsx
$ gcloud config set project PROJECT_ID
```

**Environment variable files:**

Before starting the deployment procedure, make sure you have:

- a Hedera Operator Account variables or a PEM encoded private key placed in `./config/`;
- created HCS topic;
- mongo connection string and db name;
- make sure the API_PORT variables of the both nodes match those specified in `.terraform/variables.tf`;
- activated a log with the following name: `LOG_FILE=hcsnode.log` ;

`P2P_IP` and `PEER_ADDRESS` will be configured during the deployment process.

**Provision the infrastructure:**

The command that will ensure the deployment of the infrastructure is:

```jsx
make provision
```

After successful completion, we will be able to interact with each HCS node through their external IPs.

**Deprovision the infrastructure:**

```jsx
make deprovision
```