package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
)

// Exit codes are int values that represent an exit code for a particular error.
const (
	ExitCodeOK    int = 0
	ExitCodeError int = 1 + iota
)

// CLI is the command line object
type CLI struct {
	// outStream and errStream are the stdout and stderr
	// to write message from the CLI.
	outStream, errStream io.Writer
}

func checkCurrentRecode(svc *route53.Route53, name, zoneID string) (string, int64, error) {

	params := &route53.ListResourceRecordSetsInput{
		HostedZoneId:    aws.String(zoneID), // Required
		MaxItems:        aws.String("1"),
		StartRecordName: aws.String(name),
		StartRecordType: aws.String("A"),
	}
	resp, err := svc.ListResourceRecordSets(params)
	if err != nil {
		return "", 0, fmt.Errorf("failed to get resource record")
	}

	return *resp.ResourceRecordSets[0].ResourceRecords[0].Value, *resp.ResourceRecordSets[0].TTL, nil
}

func checkCurrentIP() (string, error) {
	resp, err := http.Get("http://checkip.amazonaws.com")
	if err != nil {
		return "", fmt.Errorf("failed to get current ip")
	}
	defer resp.Body.Close()

	byteArray, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read http response")
	}

	return string(byteArray), nil
}

func upsertRecode(svc *route53.Route53, name, currentIP, zoneID string, currentTTL int64) error {
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
		return fmt.Errorf("failed to change recode")
	}
	return nil
}

// Run invokes the CLI with the given arguments.
func (c *CLI) Run(args []string) int {
	var (
		name   string
		zoneID string
	)
	flag.StringVar(&name, "name", "", "domain name")
	flag.StringVar(&zoneID, "zone_id", "", "zone id")
	flag.Parse()

	sess := session.Must(session.NewSession())
	svc := route53.New(sess)

	currentIP, err := checkCurrentIP()
	if err != nil {
		fmt.Fprintln(c.errStream, err)
		return ExitCodeError
	}
	currentIP = strings.TrimRight(currentIP, "\n")

	currentRecode, currentTTL, err := checkCurrentRecode(svc, name, zoneID)
	if err != nil {
		fmt.Fprintln(c.errStream, err)
		return ExitCodeError
	}

	currentRecode = strings.TrimRight(currentRecode, "\n")

	if currentIP == currentRecode {
		return ExitCodeOK
	}

	err = upsertRecode(svc, name, currentIP, zoneID, currentTTL)
	if err != nil {
		fmt.Fprintln(c.errStream, err)
		return ExitCodeError
	}

	return ExitCodeOK
}
