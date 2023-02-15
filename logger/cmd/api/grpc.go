package main

import (
	"context"
	"fmt"
	"log"
	"logger/data"
	"logger/logs"
	"net"
	"time"

	"google.golang.org/grpc"
)

/*
UnimplementedLogServiceServer is a placeholder type that is generated by the gRPC
code generator when a service definition file includes a service but no implementation for it.
It is used as a temporary implementation for the server-side of a gRPC service until a proper
implementation is provided. When a gRPC client attempts to invoke a method on the service,
the server will respond with an "unimplemented" error.
*/
type LogServer struct {
	logs.UnimplementedLogServiceServer
	Models data.Models
}

// First Parameter should always be of type context.Context
func (l *LogServer) WriteLog(ctx context.Context, req *logs.LogRequest) (*logs.LogResponse, error) {
	input := req.GetLogEntry()

	logEntry := data.LogEntry{
		ID:        data.GetRandomSequence(11),
		Name:      input.Name,
		Data:      input.Data,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := l.Models.LogEntry.Insert(logEntry)

	if err != nil {
		res := &logs.LogResponse{
			Result: "failed",
		}
		return res, err
	}

	res := &logs.LogResponse{
		Result: "logged!",
	}

	return res, nil
}

func (app *Config) gRPCListen() {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", GRPC_PORT))

	if err != nil {
		log.Fatalf("Failed to listen for gRPC: %v", err)
	}

	s := grpc.NewServer()

	logs.RegisterLogServiceServer(s, &LogServer{Models: app.Models}) // registering object LogServer so that
	//                                                       we can call all the methods of that object

	log.Printf("gRPC server started at port %s", GRPC_PORT)

	if err := s.Serve(listener); err != nil {
		log.Fatalf("Unable to server: %v", err)
	}
}