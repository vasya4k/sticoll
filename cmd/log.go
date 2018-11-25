package main

import "github.com/sirupsen/logrus"

const (
	ifxLogTopic          = "influx"
	eventWrPointsErr     = "write points failure"
	eventBathcPtrErr     = "new bathc points failure"
	eventClientErr       = "http client creation failure"
	eventNewClientErr    = "new client and points err"
	metricPhyIf          = "phy_interface"
	ifsLogTopic          = "interface_stats"
	eventParseFromPrxErr = "parse name from prefix val err"
	eventBathcPointrErr  = "new bathc point failure"
	eventParseQueuKeyErr = "strconv atoi err"
	eventCloseSendErr    = "close send err"
	eventRecvErr         = "grpc open config telemetry recv err"
	eventSyncRespRecv    = "recved sync resp"
	cfgErrTopic          = "config"
	cfgReadErrEv         = "read config failure"
	grpcTopic            = "grpc"
	grpcConnErrEv        = "could not connect"
	grpcLoginErrEv       = "login client err"
	grpcAuthErrEv        = "Auth failed"
	grpcSendErrEv        = "send RPC failure"
	grpcHeaderErrEv      = "get header failure"
	grpcDialOptsErrEv    = "dial options creations failure"
	tlsLogTopic          = "tls"
	tlsLoadErrEv         = "load key pair failure"
	tlsCAReadEv          = "failure read CA file"
	tlsCertAppendEv      = "failure to append certs"
	eventGRPCConnErr     = "grpc connection is nil"
)

func logErrEvent(topic, event string, err error) {
	logrus.WithFields(logrus.Fields{
		"topic": topic,
		"event": event,
	}).Error(err.Error())
}

func logFatalEvent(topic, event string, err error) {
	logrus.WithFields(logrus.Fields{
		"topic": topic,
		"event": event,
	}).Fatal(err.Error())
}

func logInfoEvent(topic, event string, msg string) {
	logrus.WithFields(logrus.Fields{
		"topic": topic,
		"event": event,
	}).Info(msg)
}
