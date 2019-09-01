#!/bin/bash

if [ -z "${AWS_ROUTE53_HOSTED_ZONE_ID}" ]; then
	echo "No zone-id provided!"
	exit 1
fi


ZONEID="${AWS_ROUTE53_HOSTED_ZONE_ID}"
DOMAIN=""
COMMENT="Updated: $(date)"


POSITIONAL=()
while [[ $# -gt 0 ]]
do
key="$1"

case $key in
    -r|--record)
    DOMAIN="$2"
    shift 
    shift 
    ;;
    -c|--comment)
    COMMENT="$2"
    shift 
    shift 
    ;;
    --ip)
    IPADDR="$2"
    shift 
    shift 
    ;;
    -v|--verbose)
    VERBOSE=0
    shift
    ;; 
    *)
    POSITIONAL+=("$1") 
    shift 
    ;;
esac
done
set -- "${POSITIONAL[@]}" # restore positional parameters


# check if ip set and exit if not
if [ -z "$IPADDR" ]; then
	echo "No IP address provided!"
	exit 1
fi


# check if record set and exit if not
if [ -z "$DOMAIN" ]; then
	echo "No Domain provided!"
	exit 1
fi



# check if root of domain or not
if [[ "${DOMAIN}" == "." ]]; then 
	DOMAIN=""
fi


BUF=$(mktemp)

sed "s/__COMMENT__/$COMMENT/" template.json > "$BUF"
sed -i "s/__DOMAIN__/$DOMAIN/" "$BUF"
sed -i "s/__IPADDR__/$IPADDR/" "$BUF"

if ${VERBOSE}; then cat "$BUF"; fi

aws route53 change-resource-record-sets --hosted-zone-id "${ZONEID}" --change-batch file:///$BUF 

exit $?

