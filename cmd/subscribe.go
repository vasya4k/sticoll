package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"time"

	na_pb "sticoll/telemetry"

	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// I wish Juniper would be doing this differently, but they are not.
// A router/switch just streams KV pairs in portions of data contained inside na_pb.OpenConfigData
// KV pairs which belong to a single entity say a physical interface are separated by a key called "__prefix__."
// But then there can be multiple occurrences of different fieldsets which belong
// to the same interface even within one OC data packet.
// Also, data types do not contain all the fields one would reasonably want.
// As a result, I have to collect information about single interface from 3 or 4 data sets.
// You can imagine how complex this can get as I have to use multiple flags
// to track which data types have already been collected.
func (d *device) subSendAndReceive(client na_pb.OpenConfigTelemetry_TelemetrySubscribeClient) {
	ifStats := newinterfaceStats(d.ifxPointCh)
	go func() {
		sigchan := make(chan os.Signal, 10)
		signal.Notify(sigchan, os.Interrupt)
		<-sigchan
		err := client.CloseSend()
		if err != nil {
			logErrEvent(grpcTopic, eventCloseSendErr, err)
		}
		os.Exit(0)
	}()
	logInfoEvent(grpcTopic, "subscribed and waiting for new data", fmt.Sprintf("hostname: %s port: %d", d.cfg.Host, d.cfg.Port))
	for {
		ocData, err := client.Recv()
		// fmt.Println("RRRRRRRRRRRRRRR", ocData)
		if err == io.EOF {
			break
		}
		if err != nil {
			logErrEvent(grpcTopic, eventRecvErr, err)
			time.Sleep(1 * time.Minute)
		}
		if ocData != nil {
			if ocData.SyncResponse {
				logInfoEvent(grpcTopic, eventSyncRespRecv, "")
			}
			var dataType string
			// Get the interesting part of the Path
			splitPath := strings.Split(ocData.Path, ":")
			if len(splitPath) >= 4 {
				dataType = splitPath[2]
			}
			// Now a path turns into a data type so it can be handeled differently
			// for _, keve := range ocData.Kv {
			// 	fmt.Printf("Path: %s key is %s and value is %s\n", dataType, keve.Key, keve.Value)
			// }
			switch dataType {
			case linecardPhyIf:
				ifStats.linecardPhyIfStats(ocData, d.cfg.Host)
			case linecardLogicIF:
			case interfaces:
				ifStats.interfaceState(ocData, d.cfg.Host)
			}
		}
	}
}

func (d *device) subscribe(conn *grpc.ClientConn) {
	if conn == nil {
		logErrEvent(grpcTopic, grpcSendErrEv, errors.New(eventGRPCConnErr))
		return
	}
	var (
		sR    na_pb.SubscriptionRequest
		adCfg na_pb.SubscriptionAdditionalConfig
		ctx   context.Context
	)
	adCfg.NeedEos = d.cfg.EOS
	for _, p := range d.cfg.Paths {
		sR.PathList = append(sR.PathList, &na_pb.Path{
			Path:            p.Path,
			SampleFrequency: uint32(p.Freq),
		})
	}
	sR.AdditionalConfig = &adCfg
	c := na_pb.NewOpenConfigTelemetryClient(conn)
	if d.cfg.Meta {
		md := metadata.New(map[string]string{
			"username": d.cfg.User,
			"password": d.cfg.Password,
		})
		ctx = metadata.NewOutgoingContext(context.Background(), md)
	} else {
		ctx = context.Background()
	}
	subClient, err := c.TelemetrySubscribe(ctx, &sR)
	if err != nil {
		logFatal(grpcTopic, grpcSendErrEv, err)
	}
	hdr, err := subClient.Header()
	if err != nil {
		logFatal(grpcTopic, grpcHeaderErrEv, err)
	}
	var headers string
	for k, v := range hdr {
		headers = headers + fmt.Sprintf("%s: %s", k, v)
	}
	logrus.WithFields(logrus.Fields{
		"topic":   grpcTopic,
		"headers": headers,
	}).Info("headers list")
	d.subSendAndReceive(subClient)
}
