serverless-aws-api-metrics
----

[![GoDoc][1]][2] [![License: MIT][3]][4] [![Release][5]][6] [![Build Status][7]][8] [![Codecov Coverage][11]][12] [![Go Report Card][13]][14] [![Code Climate][19]][20] [![BCH compliance][21]][22]

[1]: https://godoc.org/github.com/evalphobia/serverless-aws-api-metrics?status.svg
[2]: https://godoc.org/github.com/evalphobia/serverless-aws-api-metrics
[3]: https://img.shields.io/badge/License-MIT-blue.svg
[4]: LICENSE.md
[5]: https://img.shields.io/github/release/evalphobia/serverless-aws-api-metrics.svg
[6]: https://github.com/evalphobia/serverless-aws-api-metrics/releases/latest
[7]: https://github.com/evalphobia/serverless-aws-api-metrics/workflows/test/badge.svg
[8]: https://github.com/evalphobia/serverless-aws-api-metrics/actions?query=workflow%3Atest
[9]: https://coveralls.io/repos/evalphobia/serverless-aws-api-metrics/badge.svg?branch=master&service=github
[10]: https://coveralls.io/github/evalphobia/serverless-aws-api-metrics?branch=master
[11]: https://codecov.io/github/evalphobia/serverless-aws-api-metrics/coverage.svg?branch=master
[12]: https://codecov.io/github/evalphobia/serverless-aws-api-metrics?branch=master
[13]: https://goreportcard.com/badge/github.com/evalphobia/serverless-aws-api-metrics
[14]: https://goreportcard.com/report/github.com/evalphobia/serverless-aws-api-metrics
[15]: https://img.shields.io/github/downloads/evalphobia/serverless-aws-api-metrics/total.svg?maxAge=1800
[16]: https://github.com/evalphobia/serverless-aws-api-metrics/releases
[17]: https://img.shields.io/github/stars/evalphobia/serverless-aws-api-metrics.svg
[18]: https://github.com/evalphobia/serverless-aws-api-metrics/stargazers
[19]: https://codeclimate.com/github/evalphobia/serverless-aws-api-metrics/badges/gpa.svg
[20]: https://codeclimate.com/github/evalphobia/serverless-aws-api-metrics
[21]: https://bettercodehub.com/edge/badge/evalphobia/serverless-aws-api-metrics?branch=master
[22]: https://bettercodehub.com/

`serverless-aws-api-metrics` fetchs AWS API calls from CloudTrail events log and save the events count in sec by sampling into CloudWatch Custom Metrics, powered by AWS Lambda.

# Download

Download serverless-aws-api-metrics by command below.

```bash
$ git clone https://github.com/evalphobia/serverless-aws-api-metrics
$ cd serverless-aws-api-metrics
$ make init
```

# Config

## serverless.yml

Change environment variables below,

```bash
$ vim serverless.yml

------------

provider:
  name: aws
  region: ap-northeast-1  # <- Change to your target region.


...

functions:
  sns:
    handler: bin/serverless
    # tags:
    #   env: prod
    memorySize: 128
    timeout: 290
    environment:
      # Change to your target AWS services.
      METRIC_EVENT_GROUP: 'SNS'
      # Change to your target CloudTrail event name.
      METRIC_EVENT_NAME: 'GetEndpointAttributes,SetEndpointAttributes'
    events:
      - schedule: cron(*/5 * * * ? *)  # exec every 5min
```

## Environment variables

|Name|Description|Default|
|:--|:--|:--|
| `METRIC_EVENT_GROUP` | Metrics namespace and prefix used in CloudWatch Metrics. This should be service name. (e.g. `ses`, `sns`, `dynamodb`) | - |
| `METRIC_EVENT_NAME` | Event name on CloudTrail. Put multiple names with comma.	 (e.g. `GetEndpointAttributes,SetEndpointAttributes`, `CreateNetworkInterface`, `PutBucketPolicy`) | - |
| `METRIC_TARGET_SECOND` | The second of time in minute for sampling. | `31` (= `YYYY-MM-DD hh:mm:31`) |


# Deploy

```bash
$ AWS_ACCESS_KEY_ID=<...> AWS_SECRET_ACCESS_KEY=<...> make deploy
```


# Check Log

```bash
$ AWS_ACCESS_KEY_ID=<...> AWS_SECRET_ACCESS_KEY=<...> sls logs -f <function name> -t
```
