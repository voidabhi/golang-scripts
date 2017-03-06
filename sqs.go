package main

import (
	"fmt"
	"strconv"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

func main() {
	sqssvc := sqs.New(session.New())
  
  // List queues
	params := &sqs.ListQueuesInput{
		QueueNamePrefix: aws.String("prod-"),
	}
	sqs_resp, err := sqssvc.ListQueues(params)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	for _, url := range sqs_resp.QueueUrls {
		fmt.Println(*url)
	}
  
  // List queue attributes
	params := &sqs.GetQueueAttributesInput{
		QueueUrl: aws.String("https://sqs.<region>.amazonaws.com/<acc_number>/<queue_name>"),
		AttributeNames: []*string{
			aws.String("ApproximateNumberOfMessages"),
			aws.String("ApproximateNumberOfMessagesDelayed"),
			aws.String("ApproximateNumberOfMessagesNotVisible"),
		},
	}
	resp, _ := sqssvc.GetQueueAttributes(params)
	for attrib, _ := range resp.Attributes {
		prop := resp.Attributes[attrib]
		i, _ := strconv.Atoi(*prop)
		fmt.Println(attrib, i)
	}
}
