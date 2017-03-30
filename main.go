package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
)

func checkCurrentRecode(svc *route53.Route53, name, zoneID string) (string, int64) {

	params := &route53.ListResourceRecordSetsInput{
		HostedZoneId:    aws.String(zoneID), // Required
		MaxItems:        aws.String("1"),
		StartRecordName: aws.String(name),
		StartRecordType: aws.String("A"),
	}
	resp, err := svc.ListResourceRecordSets(params)

	if err != nil {
		fmt.Println(err.Error())
		return "", 0
	}

	return *resp.ResourceRecordSets[0].ResourceRecords[0].Value, *resp.ResourceRecordSets[0].TTL
}

func checkCurrentIP() string {
	resp, _ := http.Get("http://checkip.amazonaws.com")
	defer resp.Body.Close()

	byteArray, _ := ioutil.ReadAll(resp.Body)
	return string(byteArray)
}

func upsertRecode(svc *route53.Route53, name, currentIP, zoneID string, currentTTL int64) {
	params := &route53.ChangeResourceRecordSetsInput{
		ChangeBatch: &route53.ChangeBatch{ // Required
			Changes: []*route53.Change{ // Required
				{ // Required
					Action: aws.String("UPSERT"), // Required
					ResourceRecordSet: &route53.ResourceRecordSet{ // Required
						Name: aws.String(name), // Required
						Type: aws.String("A"),  // Required
						TTL:  aws.Int64(currentTTL),
						ResourceRecords: []*route53.ResourceRecord{
							{ // Required
								Value: aws.String(currentIP), // Required
							},
						},
					},
				},
			},
			Comment: aws.String("Changed by dynamic-route53"),
		},
		HostedZoneId: aws.String(zoneID),
	}

	_, err := svc.ChangeResourceRecordSets(params)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	return
}

func main() {
	var (
		name   string
		zoneID string
	)
	flag.StringVar(&name, "name", "", "domain name")
	flag.StringVar(&zoneID, "zone_id", "", "zone id")
	flag.Parse()

	sess := session.Must(session.NewSession())
	svc := route53.New(sess)

	currentIP := checkCurrentIP()
	currentIP = strings.TrimRight(currentIP, "\n")

	currentRecode, currentTTL := checkCurrentRecode(svc, name, zoneID)
	currentRecode = strings.TrimRight(currentRecode, "\n")

	if currentIP == currentRecode {
		os.Exit(0)
	}

	upsertRecode(svc, name, currentIP, zoneID, currentTTL)
}
