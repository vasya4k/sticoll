package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	na_pb "sticoll/telemetry"

	"github.com/influxdata/influxdb/client/v2"
)

const (
	linecardPhyIf   = "/junos/system/linecard/interface/"
	linecardLogicIF = "/junos/system/linecard/interface/logical/usage/"
	interfaces      = "/interfaces/"
)

type interfaceStats struct {
	pifsMap       map[string]PhyInterfaceStats
	pif           PhyInterfaceStats
	qStats        QueueStats
	prefixFound   bool
	countersFound bool
	ifxPointCh    chan ifxPoint
}

func newinterfaceStats(ifxPointCh chan ifxPoint) *interfaceStats {
	var i interfaceStats
	i.pifsMap = make(map[string]PhyInterfaceStats)
	i.qStats.queues = make([]OutQueue, 8)
	i.ifxPointCh = ifxPointCh
	return &i
}

//QueueStats ss
type QueueStats struct {
	notOnlyBufSize bool
	queues         []OutQueue
	IfName         string
	Host           string
}

//OutQueue ss
type OutQueue struct {
	OutQueue            int64
	Bytes               int64
	PeakBufferOccupancy int64
	AvgBufferOccupancy  int64
	RedDropPkts         int64
	AllocatedBufferSize int64
	Pkts                int64
	RedDropBytes        int64
}

//PhyInterfaceStats sa
type PhyInterfaceStats struct {
	Prefix                   string
	Host                     string
	Name                     string
	InitTime                 int64
	ParentAeName             string
	OperStatus               string
	CarrierTransitions       int64
	LastChange               int64
	HighSpeed                int64
	CountersOutOctets        int64
	CountersOutUnicastPkts   int64
	CountersOutMulticastPkts int64
	CountersOutBroadcastPkts int64
	CountersInOctets         int64
	CountersInUnicastPkts    int64
	CountersInMulticastPkts  int64
	CountersInBroadcastPkts  int64
	CountersInErrors         int64
	StateMtu                 int64
	StateEnabled             bool
	StateType                string
	StateCountersOutOctets   string
	StateAdminStatus         string
	StateLastChange          int64
	StateCountersInOctets    string
	StateCountersLastClear   string
	StateCountersInPkts      string
	StateCountersOutPkts     string
	OperationalState         string
	StateName                string
	StateDescription         string
	StateIfindex             int64
	StateOperStatus          string
	linePhyIf                bool
	lineQueue                bool
	ifState                  bool
	Timestamp                time.Time
}

//AddPoint add data to influx
func (pif *PhyInterfaceStats) AddPoint(inf *influxDB) {
	tags := map[string]string{
		"name":        pif.Name,
		"host":        pif.Host,
		"desc":        pif.StateDescription,
		"ae_name":     pif.ParentAeName,
		"oper_state":  pif.OperStatus,
		"admin_state": pif.StateAdminStatus,
	}
	fields := map[string]interface{}{
		"carrier_transitions":         pif.CarrierTransitions,
		"last_change":                 pif.LastChange,
		"counters_out_octets":         pif.CountersOutOctets,
		"counters_out_unicast_pkts":   pif.CountersOutUnicastPkts,
		"counters_out_multicast_pkts": pif.CountersOutMulticastPkts,
		"counters_out_broadcast_pkts": pif.CountersOutBroadcastPkts,
		"counters_in_octets":          pif.CountersInOctets,
		"counters_in_unicast_pkts":    pif.CountersInUnicastPkts,
		"counters_in_multicast_pkts":  pif.CountersInMulticastPkts,
		"counters_in_broadcast_pkts":  pif.CountersInBroadcastPkts,
		"counters_in_errors":          pif.CountersInErrors,
		"mtu":                         pif.StateMtu,
	}
	pt, err := client.NewPoint(metricPhyIf, tags, fields, time.Now())
	if err != nil {
		logErrEvent(ifxLogTopic, eventBathcPointrErr, err)
		return
	}
	inf.bp.AddPoint(pt)
}

func parseNameFromPrefixVal(prefixVal string) (string, error) {
	var rightFromEq, name string
	rightFromEqSlice := strings.Split(prefixVal, "name='")
	if len(rightFromEqSlice) > 1 {
		rightFromEq = rightFromEqSlice[1]
	} else {
		return "", errors.New("did not find \"name='\" pattern to split str ")
	}
	rightFromBracket := strings.Split(rightFromEq, "']")
	if len(rightFromBracket) > 0 {
		name = rightFromBracket[0]
	} else {
		return "", errors.New("did not find \"']\" pattern to split str ")
	}
	return name, nil
}

// fmt.Printf("Data: %08b \n", data[:4])
func printAsJSON(i interface{}) {
	b, err := json.MarshalIndent(i, "", "  ")
	if err != nil {
		log.Println(err)
	}
	fmt.Println(string(b))
}

func (s *interfaceStats) sendLinecardStats() {
	if s.countersFound {
		extPhyif, ok := s.pifsMap[s.pif.Name]
		if ok {
			if s.qStats.notOnlyBufSize {
				fmt.Println("QSTAAAATS")
				printAsJSON(s.qStats)
				// fmt.Println("QSTAAAATS", s.qStats.IfName, s.qStats.queues[0].Pkts, s.qStats.queues[1].Pkts, s.qStats.queues[2].Pkts, s.qStats.queues[3].Pkts, s.qStats.queues[4].Pkts)
			}
			extPhyif.InitTime = s.pif.InitTime
			extPhyif.ParentAeName = s.pif.ParentAeName
			extPhyif.OperStatus = s.pif.OperStatus
			extPhyif.CarrierTransitions = s.pif.CarrierTransitions
			extPhyif.LastChange = s.pif.LastChange
			extPhyif.HighSpeed = s.pif.HighSpeed
			extPhyif.CountersOutOctets = s.pif.CountersOutOctets
			extPhyif.CountersOutUnicastPkts = s.pif.CountersOutUnicastPkts
			extPhyif.CountersOutMulticastPkts = s.pif.CountersOutMulticastPkts
			extPhyif.CountersOutBroadcastPkts = s.pif.CountersOutBroadcastPkts
			extPhyif.CountersInOctets = s.pif.CountersInOctets
			extPhyif.CountersInUnicastPkts = s.pif.CountersInUnicastPkts
			extPhyif.CountersInMulticastPkts = s.pif.CountersInMulticastPkts
			extPhyif.CountersInBroadcastPkts = s.pif.CountersInBroadcastPkts
			extPhyif.CountersInErrors = s.pif.CountersInErrors
			extPhyif.Timestamp = s.pif.Timestamp
			extPhyif.linePhyIf = true
			// if extPhyif.linePhyIf {
			if extPhyif.linePhyIf && extPhyif.ifState {
				if extPhyif.Name == "ge-0/0/0" {
					printAsJSON(extPhyif)
					// fmt.Println("Got new data for linecard ", extPhyif.Name, extPhyif.CountersInUnicastPkts)
					// s.ifxPointCh <- &extPhyif
				}
			}
			s.pifsMap[s.pif.Name] = extPhyif
		} else {
			s.pif.linePhyIf = true
			s.pifsMap[s.pif.Name] = s.pif
		}
		s.countersFound = false
	}
	s.pif = *new(PhyInterfaceStats)
	s.qStats.queues = make([]OutQueue, 8)
}

func (s *interfaceStats) prefixMet(prefixVal string) {
	if !s.prefixFound {
		name, err := parseNameFromPrefixVal(prefixVal)
		if err != nil {
			logErrEvent(ifsLogTopic, eventParseFromPrxErr, err)
			return
		}
		s.pif.Name = name
	}
	if s.prefixFound {
		s.sendLinecardStats()
		name, err := parseNameFromPrefixVal(prefixVal)
		if err != nil {
			logErrEvent(ifsLogTopic, eventParseFromPrxErr, err)
			return
		}
		s.pif.Name = name
	}
	s.prefixFound = true
}

func (s *interfaceStats) linecardPhyIfStats(ocData *na_pb.OpenConfigData, hostname string) {
	var kvCount int
	kvLen := len(ocData.Kv)
	for _, kv := range ocData.Kv {
		kvCount++
		// fmt.Printf("key is %s and value is %s\n", kv.Key, kv.Value)
		switch kv.Key {
		case "init_time":
			s.pif.InitTime = kv.GetIntValue()
		case "parent_ae_name":
			s.pif.ParentAeName = kv.GetStrValue()
		case "oper-status":
			s.pif.OperStatus = kv.GetStrValue()
		case "carrier-transitions":
			s.pif.CarrierTransitions = kv.GetIntValue()
		case "last-change":
			s.pif.LastChange = kv.GetIntValue()
		case "high-speed":
			s.pif.HighSpeed = kv.GetIntValue()
		case "counters/out-octets":
			s.countersFound = true
			s.pif.CountersOutOctets = kv.GetIntValue()
		case "counters/out-unicast-pkts":
			s.pif.CountersOutUnicastPkts = kv.GetIntValue()
		case "counters/out-multicast-pkts":
			s.pif.CountersOutMulticastPkts = kv.GetIntValue()
		case "counters/out-broadcast-pkts":
			s.pif.CountersOutBroadcastPkts = kv.GetIntValue()
		case "counters/in-octets":
			s.pif.CountersInOctets = kv.GetIntValue()
		case "counters/in-unicast-pkts":
			s.pif.CountersInUnicastPkts = kv.GetIntValue()
		case "counters/in-multicast-pkts":
			s.pif.CountersInMulticastPkts = kv.GetIntValue()
		case "counters/in-broadcast-pkts":
			s.pif.CountersInBroadcastPkts = kv.GetIntValue()
		case "counters/in-errors":
			s.pif.CountersInErrors = kv.GetIntValue()
		case "__prefix__":
			s.pif.Timestamp = time.Unix(0, int64(ocData.Timestamp)*1000000)
			s.prefixMet(kv.GetStrValue())
		default:
			s.missingPrefixChk(ocData, kv.Key)
			if strings.HasPrefix(kv.Key, "out-queue") {
				s.linecardPhyIfQueue(kv)
			}
		}
		// This was once needed for something but should be perhaps removed
		if kvCount == kvLen {
			s.pif.Host = hostname
			s.sendLinecardStats()
		}
	}
}

func (s *interfaceStats) linecardPhyIfQueue(kv *na_pb.KeyValue) {
	s.qStats.IfName = s.pif.Name
	qNum, metric, err := parseQueuKey(kv.Key)
	if err == nil {
		switch metric {
		case "allocated-buffer-size":
			s.qStats.queues[qNum].AllocatedBufferSize = kv.GetIntValue()
		case "pkts":
			s.qStats.notOnlyBufSize = true
			s.qStats.queues[qNum].Pkts = kv.GetIntValue()
		case "bytes":
			s.qStats.notOnlyBufSize = true
			s.qStats.queues[qNum].Bytes = kv.GetIntValue()
		case "avg-buffer-occupancy":
			s.qStats.notOnlyBufSize = true
			s.qStats.queues[qNum].AvgBufferOccupancy = kv.GetIntValue()
		case "peak-buffer-occupancy":
			s.qStats.notOnlyBufSize = true
			s.qStats.queues[qNum].PeakBufferOccupancy = kv.GetIntValue()
		case "red-drop-pkts":
			s.qStats.notOnlyBufSize = true
			s.qStats.queues[qNum].RedDropPkts = kv.GetIntValue()
		case "red-drop-bytes":
			s.qStats.notOnlyBufSize = true
			s.qStats.queues[qNum].RedDropPkts = kv.GetIntValue()
		default:
			log.Println("Unknown metric", metric)
		}
	} else {
		logErrEvent(ifsLogTopic, eventParseQueuKeyErr, err)
	}
}

func parseQueuKey(key string) (int, string, error) {
	var (
		queueNuber int
		metric     string
		err        error
	)
	// out-queue [queue-number=0]/pkts &{121494907}
	//remove everything before =
	woEq := strings.Split(key, "=")
	if len(woEq) > 1 {
		queueNuber, err = strconv.Atoi((strings.Split(woEq[1], "]")[0]))
		if err != nil {
			return 0, "", err
		}
		metric = strings.Split(woEq[1], "/")[1]
	}
	return queueNuber, metric, nil
}

func (s *interfaceStats) missingPrefixChk(ocData *na_pb.OpenConfigData, key string) {
	if !strings.HasPrefix(key, "__") {
		if !s.prefixFound && !strings.HasPrefix(key, "/") {
			fmt.Printf("Missing prefix for sensor: %s\n", ocData.Path)
		}
	}
}

func (s *interfaceStats) interfaceState(ocData *na_pb.OpenConfigData, hostname string) {
	for _, kv := range ocData.Kv {
		switch kv.Key {
		case "name":
			s.pif.Name = kv.GetStrValue()
		case "state/type":
			s.pif.StateType = kv.GetStrValue()
		case "state/mtu":
			s.pif.StateMtu = kv.GetIntValue()
		case "state/name":
			s.pif.StateName = kv.GetStrValue()
		case "state/description":
			s.pif.StateDescription = kv.GetStrValue()
		case "state/enabled":
			s.pif.StateEnabled = kv.GetBoolValue()
		case "state/ifindex":
			s.pif.StateIfindex = kv.GetIntValue()
		case "state/admin-status":
			s.pif.StateAdminStatus = kv.GetStrValue()
		case "state/oper-status":
			s.pif.OperStatus = kv.GetStrValue()
		case "state/last-change":
			s.pif.StateLastChange = kv.GetIntValue()
		case "__prefix__":
			extPhyif, ok := s.pifsMap[s.pif.Name]
			if ok {
				extPhyif.Name = s.pif.Name
				extPhyif.StateType = s.pif.StateType
				extPhyif.StateMtu = s.pif.StateMtu
				extPhyif.StateName = s.pif.StateName
				extPhyif.StateDescription = s.pif.StateDescription
				extPhyif.StateEnabled = s.pif.StateEnabled
				extPhyif.StateIfindex = s.pif.StateIfindex
				extPhyif.StateAdminStatus = s.pif.StateAdminStatus
				extPhyif.OperStatus = s.pif.OperStatus
				extPhyif.StateLastChange = s.pif.StateLastChange
				extPhyif.Host = hostname
				// extPhyif.Timestamp = time.Unix(0, int64(ocData.Timestamp)*1000000)
				extPhyif.ifState = true
				s.pifsMap[s.pif.Name] = extPhyif
				if extPhyif.linePhyIf && extPhyif.ifState {
					if extPhyif.Name == "ge-0/0/0" {
						fmt.Println("State")
						printAsJSON(extPhyif)
						// fmt.Println("Got new state data ", extPhyif.Name, extPhyif.CountersInUnicastPkts)
						// s.ifxPointCh <- &extPhyif

					}
				}
			} else {
				extPhyif.ifState = true
				s.pifsMap[s.pif.Name] = s.pif
			}
			s.pif = *new(PhyInterfaceStats)
			s.prefixFound = true
		default:

		}
	}
}
