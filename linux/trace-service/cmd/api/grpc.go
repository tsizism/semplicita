package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"trace/data"
	"trace/trace"

	"google.golang.org/grpc"
)

type TraceServer struct {
	trace.UnimplementedTraceServiceServer
	Models data.Models
	logger *log.Logger
}

func (t *TraceServer) TraceEvent(ctx context.Context, req *trace.TraceRequest) (*trace.TraceResponse, error) {
	t.logger.Printf("TraceServer.TraceEvent req=%+v", req)

	input := req.GetTraceEntry()

	traceEntry := data.TraceEntry{
		Src:  input.Src,
		Via:  "gRPC",
		Data: input.Data,
	}

	err := t.Models.Insert(traceEntry)

	if err != nil {
		return nil, err
	}

	res := &trace.TraceResponse{Result: "gRPC logged"}
	return res, nil
}

func (appCtx applicationContext) gRPCListen() {
	listenr, err := net.Listen("tcp", fmt.Sprintf(":%d", appCtx.cfg.grpcPort))

	if err != nil {
		appCtx.logger.Fatalln("Failed to start listen err=", err)
		return
	}

	s := grpc.NewServer()

	trace.RegisterTraceServiceServer(s, &TraceServer{Models: appCtx.models, logger: appCtx.logger})

	appCtx.logger.Printf("gRPC Server started on Port=%d", appCtx.cfg.grpcPort)

	if err = s.Serve(listenr); err != nil {
		appCtx.logger.Fatalf("Failed listen for gRPC=%v", err)
	}
}
