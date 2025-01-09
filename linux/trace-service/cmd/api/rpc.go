package main

import (
	"context"
	"log"
	"time"
	"trace/data"
	
)

type RPCServer struct {}

type RPCPayload struct {
	Src  string
	Via  string
	Data string
}

// func (r RPCServer) InsertTraceEvent(mongo *mongo.Client, logger *log.Logger,payload RPCPayload, resp *string) error {
func (r RPCServer) InsertTraceEvent(payload RPCPayload, resp *string) error {
	collection := mongoClient.Database("trace").Collection("trace")

	_, err := collection.InsertOne(context.TODO(), data.TraceEntry{
        Src: payload.Src,
        Via: payload.Via,
        Data: payload.Data,
        CreatedAt: time.Now(),
    })

    if err != nil {
         log.Println("Failed to insert to Mongo", err)
        return err
    }

    *resp =  "Processed payload via RPC: " + payload.Src + "[" + payload.Via + "]"

    return nil    

}
