package main

import (
	"github.com/influxdata/influxdb/client/v2"
	"github.com/sirupsen/logrus"
)

type influxDB struct {
	Addr       string
	Pass       string
	User       string
	Precision  string
	DBName     string
	MetricName string
	BatchSize  int
	client     client.Client
	bp         client.BatchPoints
	dataCh     chan ifxPoint
}

type ifxPoint interface {
	AddPoint(inf *influxDB)
}

func (inf *influxDB) sendToInflux() {
	defer inf.client.Close()
	for d := range inf.dataCh {
		d.AddPoint(inf)
		if len(inf.bp.Points()) > inf.BatchSize {
			// Write the batch
			err := inf.client.Write(inf.bp)
			if err != nil {
				logErrEvent(ifxLogTopic, eventWrPointsErr, err)
			} else {
				logrus.WithFields(logrus.Fields{
					"topic":  ifxLogTopic,
					"event":  "wrote",
					"points": len(inf.bp.Points()),
				}).Info("wrote points into influx db")
				// Recreate a new point batch as we do not want to keep writing same points
				inf.bp, err = client.NewBatchPoints(client.BatchPointsConfig{
					Database:  inf.DBName,
					Precision: inf.Precision,
				})
				if err != nil {
					logErrEvent(ifxLogTopic, eventBathcPtrErr, err)
				}
			}
		}
	}
}

func (inf *influxDB) NewClientAndPoints() error {
	var err error
	// Create a new HTTPClient
	inf.Addr = "http://influx:8086"
	inf.User = "rooba"
	inf.Pass = "cArambaBoom"
	inf.Precision = "ms"
	inf.DBName = "ot"
	inf.BatchSize = 10

	inf.client, err = client.NewHTTPClient(client.HTTPConfig{
		Addr:     inf.Addr,
		Username: inf.User,
		Password: inf.Pass,
	})
	if err != nil {
		logErrEvent(ifxLogTopic, eventClientErr, err)
		return err
	}
	// Create a new point batch
	inf.bp, err = client.NewBatchPoints(client.BatchPointsConfig{
		Database:  inf.DBName,
		Precision: inf.Precision,
	})
	if err != nil {
		logErrEvent(ifxLogTopic, eventBathcPtrErr, err)
		return err
	}
	return nil
}
