package main

import (
	"context"
	"fmt"
	"log"
	"logger/data"
	"net"
	"net/rpc"
	"time"
)

type RPCServer struct{}

type RPCPayload struct {
	Name, Data string
}

/*
This function, to be used in rpc call, must be EXPORTABLE i.e Start with CAPITAL LETTER
*/
func (r *RPCServer) LogInfo(load RPCPayload, response *string) error {
	collection := client.Database("logs").Collection("logs")

	_, err := collection.InsertOne(context.Background(), data.LogEntry{
		ID:        data.GetRandomSequence(11),
		Name:      load.Name,
		Data:      load.Data,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	if err != nil {
		log.Println("Unable to insert into mongo for payload with name ", load.Name)
		return err
	}

	*response = "Processed payload via RPC for payload with name " + load.Name
	return nil
}

func (app *Config) rpcListen() error {
	log.Println("Starting RPC server on port", RPC_PORT)
	listen, err := net.Listen("tcp", fmt.Sprintf(":%s", RPC_PORT)) // 1. Listening on network
	if err != nil {
		return err
	}
	defer listen.Close()

	for {
		rpcConn, err := listen.Accept() // rpcConn -> net.Conn // 2. accepting connection on tht
		if err != nil {
			continue
		}
		go rpc.ServeConn(rpcConn) // Starts handling incoming `rpc` request on the connection
	}
}
