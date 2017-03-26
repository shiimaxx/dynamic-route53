package main

import (
	"fmt"

	_ "github.com/urfave/cli"

	"github.com/aws/aws-sdk-go/aws"
	_ "github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
)

func main() {
	sess := session.Must(session.NewSession())
	svc := route53.New(sess)

	params := &route53.ChangeResourceRecordSetsInput{
		ChangeBatch: &route53.ChangeBatch{ // Required
			Changes: []*route53.Change{ // Required
				{ // Required
					Action: aws.String("UPSERT"), // Required
					ResourceRecordSet: &route53.ResourceRecordSet{ // Required
						Name: aws.String("demo.fulldrive.jp."), // Required
						Type: aws.String("A"),                  // Required
						TTL:  aws.Int64(600),
						ResourceRecords: []*route53.ResourceRecord{
							{ // Required
								Value: aws.String("127.0.0.1"), // Required
							},
						},
					},
				},
			},
			Comment: aws.String("Changed by dynamic-route53"),
		},
		HostedZoneId: aws.String("Z1UTTJQH9J2GW1"),
	}

	resp, err := svc.ChangeResourceRecordSets(params)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(resp)
}
