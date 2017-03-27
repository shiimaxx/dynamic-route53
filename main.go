package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	_ "github.com/urfave/cli"

	"github.com/aws/aws-sdk-go/aws"
	_ "github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
)

func checkCurrentIP() string {
	resp, _ := http.Get("http://checkip.amazonaws.com")
	defer resp.Body.Close()

	byteArray, _ := ioutil.ReadAll(resp.Body)
	return string(byteArray)
}

func upsertRecode(name, ipAddr, zoneID string) {
	sess := session.Must(session.NewSession())
	svc := route53.New(sess)

	params := &route53.ChangeResourceRecordSetsInput{
		ChangeBatch: &route53.ChangeBatch{ // Required
			Changes: []*route53.Change{ // Required
				{ // Required
					Action: aws.String("UPSERT"), // Required
					ResourceRecordSet: &route53.ResourceRecordSet{ // Required
						Name: aws.String(name), // Required
						Type: aws.String("A"),  // Required
						TTL:  aws.Int64(600),
						ResourceRecords: []*route53.ResourceRecord{
							{ // Required
								Value: aws.String(ipAddr), // Required
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
		ipAddr string
		zoneID string
	)
	flag.StringVar(&name, "name", "", "domain name")
	flag.StringVar(&ipAddr, "ip", "", "ip address")
	flag.StringVar(&zoneID, "zone_id", "", "zone id")
	flag.Parse()

	currentIP := checkCurrentIP()

	if currentIP == ipAddr {
		os.Exit(1)
	}

	upsertRecode(name, ipAddr, zoneID)
}
