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

func (ifx *influxDB) sendToInflux() {
	defer ifx.client.Close()
	for d := range ifx.dataCh {
		d.AddPoint(ifx)
		if len(ifx.bp.Points()) > ifx.BatchSize {
			// Write the batch
			err := ifx.client.Write(ifx.bp)
			if err != nil {
				logErrEvent(ifxLogTopic, eventWrPointsErr, err)
			} else {
				logrus.WithFields(logrus.Fields{
					"topic":  ifxLogTopic,
					"event":  "wrote",
					"points": len(ifx.bp.Points()),
				}).Info("wrote points into influx db")
				// Recreate a new point batch as we do not want to keep writing same points
				ifx.bp, err = client.NewBatchPoints(client.BatchPointsConfig{
					Database:  ifx.DBName,
					Precision: ifx.Precision,
				})
				if err != nil {
					logErrEvent(ifxLogTopic, eventBathcPtrErr, err)
				}
			}
		}
	}
}

func (ifx *influxDB) NewClientAndPoints() error {
	var err error
	// Create a new HTTPClient
	ifx.Addr = "http://influx:8086"
	ifx.User = "rooba"
	ifx.Pass = "cArambaBoom"
	ifx.Precision = "ms"
	ifx.DBName = "ot"
	ifx.BatchSize = 10

	ifx.client, err = client.NewHTTPClient(client.HTTPConfig{
		Addr:     ifx.Addr,
		Username: ifx.User,
		Password: ifx.Pass,
	})
	if err != nil {
		logErrEvent(ifxLogTopic, eventClientErr, err)
		return err
	}
	// Create a new point batch
	ifx.bp, err = client.NewBatchPoints(client.BatchPointsConfig{
		Database:  ifx.DBName,
		Precision: ifx.Precision,
	})
	if err != nil {
		logErrEvent(ifxLogTopic, eventBathcPtrErr, err)
		return err
	}
	return nil
}
