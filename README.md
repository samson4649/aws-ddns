# AWS Route53 DDNS Docker Image

## Overview

A docker image running the aws-cli to allow zone record updates in a route53 zone. It uses environment variable to populate required field and run the update on an *EXISING ZONE RECORD*.

## Usage

### Building 

	docker build -t aws-ddns .

### Single usage

Run as a sinle one off update of a record 

	docker run -it --rm \
		-e AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID} \
		-e AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY} \
		-e AWS_ROUTE53_HOSTED_ZONE_ID=${AWS_ROUTE53_HOSTED_ZONE_ID} \
		aws-ddns \
		-r "${domain}.${TLD}" \
		--ip "${IPADDR}"

### Cron example
	
Configure this as a cron script to be run at intervals

	#!/bin/bash

	# Local Config
	TLD=<top_level_domain>

	# get current machine public IP
	IPADDR=$(dig @resolver1.opendns.com +short myip.opendns.com)
	
	# AWS Route53 Config
	AWS_ACCESS_KEY_ID=<AWS_ACCESS_KEY_ID>
	AWS_SECRET_ACCESS_KEY=<AWS_SECRET_ACCESS_KEY>
	AWS_ROUTE53_HOSTED_ZONE_ID=<AWS_ROUTE53_ZONE_ID>

	# domains to update
	DDNS_DOMAINS="domaina domainb sub.domainc"
	
	
	for domain in $DDNS_DOMAINS; do 
		if EXADDR=$(dig +short ${domain}.${TLD} | egrep '^([0-9]{1,3}\.){3}[0-9]{1,3}$'); then
			if [ "${EXADDR}" != "${IPADDR}" ]; then
				echo "Change detected: ${domain}.${TLD}"
				docker run -it --rm \
					-e AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID} \
					-e AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY} \
					-e AWS_ROUTE53_HOSTED_ZONE_ID=${AWS_ROUTE53_HOSTED_ZONE_ID} \
					awscli \
					-r "${domain}.${TLD}" \
					--ip "${IPADDR}"
			else
				echo "No change: ${domain}.${TLD}"
			fi
		else
			echo "Failed to get IP address for domain: ${domain}.${TLD}"
		fi
	done
