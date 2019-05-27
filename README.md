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

### Cron 
	
A cron script example ```./build_files/cron.sh``` can be used as a script to be run at intervals

