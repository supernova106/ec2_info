## Description

[![Build Status](https://travis-ci.org/supernova106/ec2_info.svg)](https://travis-ci.org/supernova106/ec2_info)
[![Join the chat at https://gitter.im/ec2_info/Lobby](https://badges.gitter.im/ec2_info/Lobby.svg)](https://gitter.im/ec2_info/Lobby?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)
- EC2 info API
- dependencies:

## Usage
### Describe Instances
- `source .env`
- *region* default value `us-east-1`
- *instance-state-name* default value `running`
- *state* default value `active`

```
# describe all running instances in a region
localhost/describe
localhost/describe?region=us-west-2&instance-state-name=running
# check for active reserved instances
localhost/describe?reserved=1
localhost/describe?reserved=1&state=retired
```

### Get EC2 metrics
```
localhost/utilization?InstanceID=xxxxxx
localhost/utilization?InstanceID=xxxxxx&MetricName=CPUUtilization
```

### Get Ondemand EC2 prices
```
localhost/price
```

## Todo
- rate limiting
- add caching layer
- update UI components
- export to JSON format

## References
- https://github.com/aws/aws-sdk-go/tree/f8f7a96133a04892b935a2ae9bebdf80bf7c6397/example 
- https://docs.aws.amazon.com/sdk-for-go/api/service/ec2/#EC2.DescribeInstances
- https://godoc.org/github.com/datacratic/aws-sdk-go/service/cloudwatch#example-CloudWatch-GetMetricStatistics 
- https://github.com/aws/aws-sdk-go 
- https://blog.golang.org/go-slices-usage-and-internals
- https://golang.org/doc/effective_go.html

## Contact
- Binh Nguyen
