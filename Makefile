GCP_PROJECT ?= hcs-integration-node
GCLOUD_EXEC ?= gcloud
TERRAFORM_EXEC ?= terraform
SSH_KEYGEN_EXEC ?= ssh-keygen
SSH_USER ?= ubuntu

ROOT_DIR=$(shell pwd)
TERRAFORM_ROOT_DIR=${ROOT_DIR}/terraform

TMP_DIR=${ROOT_DIR}/.tmp
TERRAFORM_OUTPUT = ${TERRAFORM_ROOT_DIR}/terraform.output

set_gcp_project:
	set -e

	${GCLOUD_EXEC} config set project ${GCP_PROJECT}

terraform_deploy:
	set -e

	cd ${TERRAFORM_ROOT_DIR} && \
		${SSH_KEYGEN_EXEC} -b 2048 -t rsa -f ${SSH_USER} -q -N "" && \
		${TERRAFORM_EXEC} init -no-color && \
		${TERRAFORM_EXEC} plan -no-color --out ${TERRAFORM_ROOT_DIR}/provision.plan && \
		${TERRAFORM_EXEC} apply -no-color ${TERRAFORM_ROOT_DIR}/provision.plan && \
		${TERRAFORM_EXEC} output -no-color --json > ${TERRAFORM_OUTPUT}

copy_project:
	set -e

	$(eval external_ip := ${shell cat ${TERRAFORM_OUTPUT} | jq -r .external_ip.value})
	echo ${external_ip}
	${SSH_KEYGEN_EXEC} -R ${external_ip}
	scp -i ${TERRAFORM_ROOT_DIR}/${SSH_USER} -o StrictHostKeyChecking=no -q -o LogLevel=QUIET -r ./* ${SSH_USER}@${external_ip}:/home/${SSH_USER}

provision: set_gcp_project terraform_deploy copy_project