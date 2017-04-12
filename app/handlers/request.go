package request

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/gin-gonic/gin"
	"github.com/robertkrimen/otto"
	"github.com/supernova106/ec2_info/app/config"
	"github.com/supernova106/ec2_info/app/models"
	"io/ioutil"
	"log"
	"net/http"
	"runtime"
	"strings"
	"sync"
	"time"
)

func DescribeEC2(c *gin.Context) {
	awsRegion := c.DefaultQuery("region", "us-east-1")
	instanceStateName := c.DefaultQuery("instance-state-name", "running")
	reservedFlag := c.DefaultQuery("reserved", "0")
	instanceIdsString := c.DefaultQuery("instanceIds", "")

	if "1" == reservedFlag {
		state := c.DefaultQuery("state", "active")
		resp := getDescribeReservedEC2(awsRegion, state)
		c.JSON(200, resp)
	} else {
		var instanceIds []*string
		if strings.Contains(instanceIdsString, ",") {
			s := strings.Split(instanceIdsString, ",")
			for _, value := range s {
				instanceIds = append(instanceIds, aws.String(value))
			}
		} else if strings.Contains(instanceIdsString, "i-") {
			instanceIds = append(instanceIds, aws.String(instanceIdsString))
		}

		var params *ec2.DescribeInstancesInput
		if len(instanceIds) > 0 {
			params = &ec2.DescribeInstancesInput{
				Filters: []*ec2.Filter{
					&ec2.Filter{
						Name: aws.String("instance-state-name"),
						Values: []*string{
							aws.String(strings.Join([]string{"*", instanceStateName, "*"}, "")),
						},
					},
				},
				InstanceIds: instanceIds,
			}
		} else {
			params = &ec2.DescribeInstancesInput{
				Filters: []*ec2.Filter{
					&ec2.Filter{
						Name: aws.String("instance-state-name"),
						Values: []*string{
							aws.String(strings.Join([]string{"*", instanceStateName, "*"}, "")),
						},
					},
				},
			}
		}

		resp := getDescribeEC2(awsRegion, params)
		c.JSON(200, resp.Reservations)
	}

	return
}

func getDescribeReservedEC2(awsRegion string, state string) *ec2.DescribeReservedInstancesOutput {
	sess, err := session.NewSession()
	if err != nil {
		panic(err)
	}

	// Create an EC2 service object in the "us-west-2" region
	// Note that you can also configure your region globally by
	// exporting the AWS_REGION environment variable
	svc := ec2.New(sess, &aws.Config{Region: aws.String(awsRegion)})
	params := &ec2.DescribeReservedInstancesInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name: aws.String("state"),
				Values: []*string{
					aws.String(strings.Join([]string{"*", state, "*"}, "")),
				},
			},
		},
	}

	resp, err := svc.DescribeReservedInstances(params)
	if err != nil {
		fmt.Println("there was an error listing instances in", awsRegion, err.Error())
		log.Fatal(err.Error())
	}

	return resp
}

func getDescribeEC2(awsRegion string, params *ec2.DescribeInstancesInput) *ec2.DescribeInstancesOutput {
	sess, err := session.NewSession()
	if err != nil {
		panic(err)
	}
	svc := ec2.New(sess, &aws.Config{Region: aws.String(awsRegion)})
	resp, err := svc.DescribeInstances(params)
	if err != nil {
		fmt.Println("there was an error listing instances in", awsRegion, err.Error())
		log.Fatal(err.Error())
	}

	return resp
}

func Utilization(c *gin.Context) {
	awsRegion := c.DefaultQuery("region", "us-east-1")
	instanceId := c.Query("InstanceId")
	metricName := c.DefaultQuery("MetricName", "CPUUtilization")
	parallel := runtime.NumCPU() * 2
	if instanceId == "" {
		c.JSON(400, gin.H{"error": "InstanceID is missing!"})
		return
	}

	var instanceList []string
	var dataMetric map[string]*cloudwatch.GetMetricStatisticsOutput
	dataMetric = make(map[string]*cloudwatch.GetMetricStatisticsOutput)

	if instanceId == "all" {
		// get list of running instances
		params := &ec2.DescribeInstancesInput{
			Filters: []*ec2.Filter{
				&ec2.Filter{
					Name: aws.String("instance-state-name"),
					Values: []*string{
						aws.String(strings.Join([]string{"*", "running", "*"}, "")),
					},
				},
			},
		}
		resp := getDescribeEC2(awsRegion, params)
		// resp has all of the response data, pull out instance IDs:
		for idx, _ := range resp.Reservations {
			for _, inst := range resp.Reservations[idx].Instances {
				instanceList = append(instanceList, *inst.InstanceId)
			}
		}

		// on the number of CPUs available.
		chunkSize := len(instanceList) / parallel

		if chunkSize < parallel {
			chunkSize = len(instanceList)
		}

		for start := 0; start < len(instanceList); start += chunkSize {
			end := start + chunkSize
			if end > len(instanceList) {
				end = len(instanceList)
			}

			var complete sync.WaitGroup
			for _, eachEC2 := range instanceList[start:end] {
				complete.Add(1)
				go func(eachEC2 string) {
					defer complete.Done()
					fmt.Println("Processing ", eachEC2)
					dataMetric[eachEC2] = getUtilization(awsRegion, eachEC2, metricName)
				}(eachEC2)
			}
			complete.Wait()
		}
	} else {
		dataMetric[instanceId] = getUtilization(awsRegion, instanceId, metricName)
	}

	c.JSON(200, dataMetric)
	return
}

func getUtilization(awsRegion string, instanceId string, metricName string) *cloudwatch.GetMetricStatisticsOutput {
	sess, err := session.NewSession()
	if err != nil {
		panic(err)
	}
	// Create new cloudwatch client.
	svc := cloudwatch.New(sess, &aws.Config{Region: aws.String(awsRegion)})
	now := time.Now()
	prev := now.Add(time.Duration(168) * time.Hour * -1)

	params := &cloudwatch.GetMetricStatisticsInput{
		MetricName: aws.String(metricName),
		Namespace:  aws.String("AWS/EC2"),
		Dimensions: []*cloudwatch.Dimension{
			&cloudwatch.Dimension{
				Name:  aws.String("InstanceId"),
				Value: aws.String(instanceId),
			},
		},
		Period:     aws.Int64(86400),
		Statistics: []*string{aws.String("Average"), aws.String("Maximum"), aws.String("Minimum")},
		StartTime:  aws.Time(prev),
		EndTime:    aws.Time(now),
	}

	result, err := svc.GetMetricStatistics(params)
	if err != nil {
		fmt.Println("there was an error listing instances in", awsRegion, err.Error())
		log.Fatal(err.Error())
	}

	return result
}

func GetData(c *gin.Context) {
	cfg := c.MustGet("cfg").(*config.Config)
	awsPrice := getAWSPrices(cfg.LinuxOdPriceUrl)
	prevAwsPrice := getAWSPrices(cfg.LinuxOdPricePreviousUrl)
	c.JSON(200, gin.H{"currentGen": awsPrice, "previousGen": prevAwsPrice})
	return
}

func getAWSPrices(url string) *models.AWSPrice {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error connecting to ", url, err.Error())
		return nil
	}

	defer resp.Body.Close()
	jsObjectBytes, _ := ioutil.ReadAll(resp.Body)

	vm := otto.New()
	vm.Set("jsObject", string(jsObjectBytes))
	vm.Run(`
    var callback= function(x) {
        return eval(x);
    };
    var awsRead = function(x) {
        return eval(x);
    };
    var jsObject = awsRead(jsObject);
    var jsonString = JSON.stringify(jsObject);
    // The value of def is 11
`)

	value, err := vm.Get("jsonString")
	if err != nil {
		fmt.Println("Unable to get the JSON String from the JS VM")
		return nil
	}
	awsPrice := &models.AWSPrice{}
	err = json.Unmarshal([]byte(value.String()), awsPrice)
	if err != nil {
		fmt.Println("Unable to parse the JS JSON String to Go Struct")
		return nil
	}

	return awsPrice
}

func Check(c *gin.Context) {
	c.String(200, "Hello! It's running!")
}
