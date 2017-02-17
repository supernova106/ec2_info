package request

import (
	"github.com/supernova106/ec2_info/app/config"
	"github.com/supernova106/ec2_info/app/models"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/robertkrimen/otto"
	"io/ioutil"
	"net/http"
)

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
