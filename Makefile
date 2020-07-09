GCP_PROJECT ?= hcs-integration-node
GCLOUD_EXEC ?= gcloud
TERRAFORM_EXEC ?= terraform
SSH_KEYGEN_EXEC ?= ssh-keygen
SSH_USER_PEER1 ?= ubuntu_peer1
SSH_USER_PEER2 ?= ubuntu_peer2

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
		${SSH_KEYGEN_EXEC} -b 2048 -t rsa -f ${SSH_USER_PEER1} <<<y 2>&1 >/dev/null -q -N "" && \
		${SSH_KEYGEN_EXEC} -b 2048 -t rsa -f ${SSH_USER_PEER2} <<<y 2>&1 >/dev/null -q -N "" && \
		${TERRAFORM_EXEC} init -no-color && \
		${TERRAFORM_EXEC} plan -no-color --out ${TERRAFORM_ROOT_DIR}/provision.plan && \
		${TERRAFORM_EXEC} apply -no-color ${TERRAFORM_ROOT_DIR}/provision.plan && \
		${TERRAFORM_EXEC} output -no-color --json > ${TERRAFORM_OUTPUT}

prepare_env:
	set -e

	$(eval internal_ip_peer2 := ${shell cat ${TERRAFORM_OUTPUT} | jq -r .internal_ip_peer2.value})
	echo "" >>./peer2.env
	echo "P2P_IP=${internal_ip_peer2}" >>./peer2.env
	cat peer2.env

copy_project_peer2:
	set -e

	$(eval external_ip_peer2 := ${shell cat ${TERRAFORM_OUTPUT} | jq -r .external_ip_peer2.value})
	${SSH_KEYGEN_EXEC} -R ${external_ip_peer2}
	scp -i ${TERRAFORM_ROOT_DIR}/${SSH_USER_PEER2} -o StrictHostKeyChecking=no -q -o LogLevel=QUIET -r ./* ${SSH_USER_PEER2}@${external_ip_peer2}:/home/${SSH_USER_PEER2}

start_peer2:
	set -e

	ssh -t -i ${TERRAFORM_ROOT_DIR}/${SSH_USER_PEER2} ${SSH_USER_PEER2}@${external_ip_peer2} sudo 'chmod +x /home/${SSH_USER_PEER2}/install_prerequisites.sh && chmod +x /home/${SSH_USER_PEER2}/start_peer2_tmux.sh'
	ssh -i ${TERRAFORM_ROOT_DIR}/${SSH_USER_PEER2} ${SSH_USER_PEER2}@${external_ip_peer2} 'bash /home/${SSH_USER_PEER2}/install_prerequisites.sh'
	ssh -i ${TERRAFORM_ROOT_DIR}/${SSH_USER_PEER2} ${SSH_USER_PEER2}@${external_ip_peer2} 'bash /home/${SSH_USER_PEER2}/start_peer2_tmux.sh'

get_node_log:
	set -e

	$(eval external_ip_peer2 := ${shell cat ${TERRAFORM_OUTPUT} | jq -r .external_ip_peer2.value})
	scp -i ${TERRAFORM_ROOT_DIR}/${SSH_USER_PEER2} -o StrictHostKeyChecking=no -q ${SSH_USER_PEER2}@${external_ip_peer2}:/home/${SSH_USER_PEER2}/hcsnode.log ./

prepare_second_peer_env:
	set -e

	$(eval lib_line := ${shell ls -l | grep -A0 "ip4" ./hcsnode.log})
	$(eval lib_env := ${shell echo ${lib_line} | rev | cut -d' ' -f 2 | rev})
	echo ${lib_env}

	echo "" >>./.env
	echo "PEER_ADDRESS=${lib_env}" >>./.env

	$(eval internal_ip_peer1 := ${shell cat ${TERRAFORM_OUTPUT} | jq -r .internal_ip_peer1.value})
	echo "" >>./.env
	echo "P2P_IP=${internal_ip_peer1}" >>./.env
	cat .env

copy_project_peer1:
	set -e

	$(eval external_ip_peer1 := ${shell cat ${TERRAFORM_OUTPUT} | jq -r .external_ip_peer1.value})
	${SSH_KEYGEN_EXEC} -R ${external_ip_peer1}
	scp -i ${TERRAFORM_ROOT_DIR}/${SSH_USER_PEER1} -o StrictHostKeyChecking=no -q -o LogLevel=QUIET -r ./* ${SSH_USER_PEER1}@${external_ip_peer1}:/home/${SSH_USER_PEER1}
	scp -i ${TERRAFORM_ROOT_DIR}/${SSH_USER_PEER1} -o StrictHostKeyChecking=no -q -o LogLevel=QUIET -r ./.env ${SSH_USER_PEER1}@${external_ip_peer1}:/home/${SSH_USER_PEER1}

start_peer1:
	set -e

	$(eval external_ip_peer1 := ${shell cat ${TERRAFORM_OUTPUT} | jq -r .external_ip_peer1.value})
	ssh -t -i ${TERRAFORM_ROOT_DIR}/${SSH_USER_PEER1} ${SSH_USER_PEER1}@${external_ip_peer1} sudo 'chmod +x /home/${SSH_USER_PEER1}/install_prerequisites.sh && chmod +x /home/${SSH_USER_PEER1}/start_tmux.sh'
	ssh -i ${TERRAFORM_ROOT_DIR}/${SSH_USER_PEER1} ${SSH_USER_PEER1}@${external_ip_peer1} 'bash /home/${SSH_USER_PEER1}/install_prerequisites.sh'
	ssh -i ${TERRAFORM_ROOT_DIR}/${SSH_USER_PEER1} ${SSH_USER_PEER1}@${external_ip_peer1} 'bash /home/${SSH_USER_PEER1}/start_tmux.sh'

provision: set_gcp_project terraform_deploy prepare_env copy_project_peer2 start_peer2 get_node_log prepare_second_peer_env copy_project_peer1 start_peer1

destroy_infrastructure:
	set -e

	cd ${TERRAFORM_ROOT_DIR} && \
		${TERRAFORM_EXEC} destroy

clear_log_file:
	set -e

	rm -f ./hcsnode.log

deprovision: destroy_infrastructure clear_log_file