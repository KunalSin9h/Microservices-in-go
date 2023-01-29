package data

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func New(m *mongo.Client) Models {
	client = m
	return Models{
		LogEntry: LogEntry{},
	}
}

/*
Models -> all the tables / data models
*/
type Models struct {
	LogEntry LogEntry
}

/*
LogEntry is data model for each log entry
*/
type LogEntry struct {
	ID        string    `bson:"_id" json:"id"`
	Name      string    `bson:"name" json:"name"`
	Data      string    `bson:"data" json:"data"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

/*
Insert data into mongodb
*/
func (l *LogEntry) Insert(entry LogEntry) error {

	collection := client.Database("logs").Collection("logs")

	_, err := collection.InsertOne(context.TODO(), LogEntry{
		Name:      entry.Name,
		Data:      entry.Data,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	if err != nil {
		return err
	}

	return nil

}

/*
All get all the log entries
*/
func (l *LogEntry) All() ([]*LogEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("logs").Collection("logs")

	opts := options.Find()
	opts.SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := collection.Find(context.TODO(), bson.D{}, opts)

	if err != nil {
		return nil, errors.New("Error getting log entries from mongodb: " + err.Error())
	}

	var logs []*LogEntry

	for cursor.Next(ctx) {
		var entry LogEntry
		err := cursor.Decode(&entry)
		if err != nil {
			return nil, errors.New("Error Decoding log entry from mongodb collection: " + err.Error())
		}
		logs = append(logs, &entry)
	}

	return logs, nil
}

/*
GetLogEntryById yes this is what it is
*/
func (l *LogEntry) One(id string) (*LogEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("logs").Collection("logs")

	docId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return nil, errors.New("Error getting id form hex: " + err.Error())
	}

	var entry LogEntry

	err = collection.FindOne(ctx, bson.M{"_id": docId}).Decode(&entry)

	if err != nil {
		return nil, errors.New("Error getting log entry form id: " + err.Error())
	}

	return &entry, nil
}

/*
ClearLogEntries clear all logs up till now
Alt Name : DropCollection
*/
func (l *LogEntry) Clear() error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("logs").Collection("logs")

	if err := collection.Drop(ctx); err != nil {
		return errors.New("Error dropping logs collection: " + err.Error())
	}

	return nil
}

/*
UpdateLogEntry yes
*/

func (l *LogEntry) Update() (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("logs").Collection("logs")

	docId, err := primitive.ObjectIDFromHex(l.ID)

	if err != nil {
		return nil, errors.New("Error while getting id from hex: " + err.Error())
	}

	res, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": docId},
		bson.D{
			{Key: "$set", Value: bson.D{
				{Key: "name", Value: l.Name},
				{Key: "data", Value: l.Data},
				{Key: "updated_at", Value: time.Now()},
			}},
		},
	)

	if err != nil {
		return nil, errors.New("Error updating single log entry: " + err.Error())
	}

	return res, nil
}
