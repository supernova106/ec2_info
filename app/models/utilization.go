package models

import "github.com/aws/aws-sdk-go/service/cloudwatch"

type Utilization struct {
	InstanceId  string                                `json:"InstanceId"`
	Utilization *cloudwatch.GetMetricStatisticsOutput `json:"Utilization"`
}
