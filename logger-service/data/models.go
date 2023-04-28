package data

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

var client *mongo.Client

type LogEntry struct {
	ID        string    `bson:"_id,omitempty" json:"id,omitempty"`
	Title     string    `bson:"title" json:"title"`
	Data      any       `bson:"data" json:"data"`
	CreatedAt time.Time `bson:"created_at" json:"createdAt"`
	UpdatedAt time.Time `bson:"updated_at" json:"updatedAt"`
}

type Models struct {
	LogEntry *LogEntry
}

func NewModels(c *mongo.Client) *Models {
	client = c
	return &Models{
		&LogEntry{},
	}
}

func getCollection() *mongo.Collection {
	return client.Database("logs").Collection("logs")
}

func (l *LogEntry) Insert(entry LogEntry) error {
	collection := getCollection()
	res, err := collection.InsertOne(context.TODO(), &LogEntry{
		Title:     entry.Title,
		Data:      entry.Data,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	log.Println(fmt.Sprintf("inserted id: %s", res.InsertedID))

	if err != nil {
		log.Println("error inserting log entry", err)
		return err
	}
	return nil
}

func (l *LogEntry) All() ([]*LogEntry, error) {
	collection := getCollection()

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	opt := options.Find()
	opt.SetSort(bson.D{{"created_at", -1}})
	cursor, err := collection.Find(ctx, bson.D{{}}, opt)
	if err != nil {
		log.Println("error retrieving documents", err)
		return nil, err
	}

	var logs []*LogEntry
	for cursor.Next(ctx) {
		var l LogEntry
		err := cursor.Decode(&l)
		if err != nil {
			log.Println("error decoding LogEntry", err)
			return nil, err
		}
		logs = append(logs, &l)
	}
	return logs, nil
}

func (l *LogEntry) GetOne(id string) (*LogEntry, error) {
	collection := getCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	logId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("error decoding id", err)
		return nil, err
	}

	var entry LogEntry

	err = collection.FindOne(ctx, bson.M{"_id": logId}).Decode(&entry)
	if err != nil {
		log.Println("error finding document", err)
		return nil, err
	}
	return &entry, nil
}

func (l *LogEntry) DropCollection() error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	collection := getCollection()

	err := collection.Drop(ctx)
	if err != nil {
		log.Println("error dropping collection", err)
		return err
	}

	log.Println("collection dropped!")
	return nil
}

func (l *LogEntry) Update() (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	collection := getCollection()

	logId, err := primitive.ObjectIDFromHex(l.ID)
	if err != nil {
		log.Println("error decoding id", err)
		return nil, err
	}
	result, err := collection.UpdateOne(ctx, bson.M{"_id": logId}, bson.D{
		{"$set", bson.D{
			{"title", l.Title},
			{"data", l.Data},
			{"updated_at", time.Now()},
		}},
	})
	if err != nil {
		log.Println("error updating log entry", err)
		return nil, err
	}
	return result, nil
}
