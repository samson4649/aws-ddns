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

