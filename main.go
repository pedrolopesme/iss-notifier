package main

import (
	"os"
	"errors"
	log "github.com/sirupsen/logrus"
	"github.com/pedrolopesme/iss-notifier/sqs"
	distanceCalculator "github.com/pedrolopesme/distance-calculator"
	"github.com/pedrolopesme/iss-tracker/iss"
	"strconv"
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

// Checks it ISS pass over a coordinate
// allowing a tolerance of 10,000 Km
func ChecksIntersection(position iss.IssPosition) bool {
	tolerance := float64(10000)
	brazilCoordinate := Coordinate{ Latitude: -13.6578431, Longitude: -69.7095687}

	issLat, err := strconv.ParseFloat(position.Latitude, 64)
	if err != nil {
		log.Error("It was impossible to cast %s to float64", issLat)
		return false
	}

	issLong, err := strconv.ParseFloat(position.Longitude, 64)
	if err != nil {
		log.Error("It was impossible to cast %s to float64", issLong)
		return false
	}

	distance := distanceCalculator.CalcKilometers(brazilCoordinate.Latitude, brazilCoordinate.Longitude, issLat, issLong)
	return distance < tolerance
}

func main(){

	// Retrieving ISS_SQS_URL from the ENV vars
	queueURL := os.Getenv(SQSEnvVar)
	if queueURL == "" {
		err := ErrSQSVarNotFound
		log.Error(err)
		return
	}

	log.Info("Checking out SQS Queue")
	position, err := sqs.Consume(queueURL)
	if err != nil {
		log.Error("It was impossible to retrieve ISS position", err)
		return
	}

	log.Info("Position found %s", position)
	if ChecksIntersection(position) {
		log.Info("ISS is passing over")
	} else {
		log.Info("ISS is NOT passing over")
	}
}

