dynamic-route53
===============

[![Build Status](https://travis-ci.org/shiimaxx/dynamic-route53.svg?branch=master)](https://travis-ci.org/shiimaxx/dynamic-route53)

Dynamic DNS using AWS Route53

## Description

dynamic-route53 is Dynamic DNS client using AWS Route53.

## Usage

Set AWS credentials to environment variables.

```
$ export AWS_ACCESS_KEY_ID=<access_key>
$ export AWS_SECRET_ACCESS_KEY=<secret_key>
```

Execute specify domain name and zone id.

```
dynamic-route53 --name <DOMAIN NAME> --zone_id <YOURE ZONE ID>
```

## Author

[shiimaxx](https://github.com/shiimaxx)
