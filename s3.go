package main
 
/* https://github.com/aws/aws-sdk-go/wiki */
 
import (
    "fmt"
 
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/s3"
)
 
func main() {
    /*
       export AWS_REGION=""
       export AWS_ACCESS_KEY_ID=""
       export AWS_SECRET_ACCESS_KEY=""*/
 
    s3svc := s3.New(session.New())
    result, err := s3svc.ListBuckets(&s3.ListBucketsInput{})
    if err != nil {
        fmt.Println("Failed to list buckets", err)
        return
    }
 
    fmt.Println("Buckets:")
    for _, bucket := range result.Buckets {
        fmt.Printf("%s : %s\n", aws.StringValue(bucket.Name), bucket.CreationDate)
    }
}
