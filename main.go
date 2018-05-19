package main

import (
	"fmt"
	"os"
	"errors"
	log "github.com/sirupsen/logrus"
	"github.com/pedrolopesme/iss-notifier/sqs"
)

const (
	// URL to our queue
	SQSEnvVar = "ISS_SQS_URL"
)

var (
	// ErrSQSVarNotFound is returned when a ENV Var with ISS SQS URL wasn't found
	ErrSQSVarNotFound = errors.New("environment variable with ISS SQS URL not set. Please define a env var " + SQSEnvVar)
)

// ISS Position relative to Earth
type Coordinate struct {
	Latitude float64
	Longitude float64
}

func main(){

	// Retrieving ISS_SQS_URL from the ENV vars
	queueURL := os.Getenv(SQSEnvVar)
	if queueURL == "" {
		err := ErrSQSVarNotFound
		log.Error(err)
		return
	}

	fmt.Print("Checking out SQS Queue")
	sqs.Consume(queueURL)
}

