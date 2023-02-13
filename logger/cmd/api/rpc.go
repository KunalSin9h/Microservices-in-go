package main

import (
	"context"
	"log"
	"logger/data"
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
