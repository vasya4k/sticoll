package main

import (
	"sync"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc/stats"
)

type gRPCStats struct {
	sync.Mutex               // guarding following stats
	startTime                time.Time
	totalIn                  uint64
	totalKV                  uint64
	totalInPayloadLength     uint64
	totalInPayloadWireLength uint64
	totalInHeaderWireLength  uint64
	totalLatency             uint64
	totalLatencyPkt          uint64
	totalDdrops              uint64
}

type statsHandler struct {
	cfg *device
}

func (h *statsHandler) TagConn(ctx context.Context, info *stats.ConnTagInfo) context.Context {
	return ctx
}

func (h *statsHandler) TagRPC(ctx context.Context, info *stats.RPCTagInfo) context.Context {
	return ctx
}

func (h *statsHandler) HandleConn(ctx context.Context, s stats.ConnStats) {
	switch s.(type) {
	case *stats.ConnBegin:
	case *stats.ConnEnd:
	default:
	}
}

func (h *statsHandler) HandleRPC(ctx context.Context, s stats.RPCStats) {
	h.cfg.Stats.Lock()
	defer h.cfg.Stats.Unlock()

	switch s.(type) {
	case *stats.InHeader:
		h.cfg.Stats.totalInHeaderWireLength += uint64(s.(*stats.InHeader).WireLength)
	case *stats.OutHeader:
	case *stats.OutPayload:
	case *stats.InPayload:
		h.cfg.Stats.totalInPayloadLength += uint64(s.(*stats.InPayload).Length)
		h.cfg.Stats.totalInPayloadWireLength += uint64(s.(*stats.InPayload).WireLength)

	case *stats.InTrailer:
	case *stats.End:
	default:
	}
}
