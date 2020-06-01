# AWS Route53 DDNS Docker Image

## Overview

A simple docker image with configuration file to run as a DDNS service.

## Configuration

Application settings are configured with environment variables. Domain configuration is managed by a yml formatted configuration file:
```yaml
---
domains:
  - zone_id: <aws_zone_id> 
    zone_name: <domain_fqdn> 
    values:
      - record: <subdomain> 
      # ... 
      - record: <subdomain>
```

## Usage

### Environment Variables

#### AWS_ACCESS_KEY (required)
Access key with permissions to edit the zone in AWS IAM

#### AWS_SECRET_KEY (required) 
Secret Key to match access key above

### Docker Compose

```bash
--
version: "3.8"
services:
  ddns:
    image: samson4649/aws-ddns:latest
    environment:
      AWS_ACCESS_KEY: <access_key>
      AWS_SECRET_KEY: <secret_key> 
    volumes:
      - type: bind
        source: /path/t/config.yml
        target: /etc/aws-ddns/aws-ddns.yml
```

### Docker Swarm 

```bash
--
version: "3.8"
services:
  ddns:
    image: samson4649/aws-ddns:latest
    environment:
      AWS_ACCESS_KEY: <access_key>
      AWS_SECRET_KEY: <secret_key> 
    command: 
      - "--config=/run/ddns.yml"
    configs:
      - source: ddns-config
        target: /run/ddns.yml
        mode: 0444

configs:
  ddns-config:
    file: /path/to/config.yml
```

### Single usage

Run as a sinle one off update of a record 

	docker run -it --rm \
		-e AWS_ACCESS_KEY=${AWS_ACCESS_KEY_ID} \
		-e AWS_SECRET__KEY=${AWS_SECRET_ACCESS_KEY} \
		samson4649/aws-ddns \
		--config=myconfig.yml

