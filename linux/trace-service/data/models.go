package data

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Models struct {
	mongo  *mongo.Client
	logger *log.Logger
}

type TraceEntry struct {
	ID        string    `bson:"_id,omitempty" json:"id,omitempty"`
	Src      string     `bson:"src" json:"src"`
	Via      string     `bson:"via" json:"via"`
	Data      string    `bson:"data" json:"data"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	// UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

// New Models instance
func New(mongo *mongo.Client, logger *log.Logger) Models {
	return Models{
		mongo:  mongo,
		logger: logger,
	}
}

func (m *Models) Insert(e TraceEntry) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := m.mongo.Database("trace").Collection("trace")

	res, err := collection.InsertOne(ctx, TraceEntry{
		Src:      	e.Src,
		Via:      	e.Via,
		Data:      	e.Data,
		CreatedAt: 	time.Now(),
		// UpdatedAt: 	time.Now(),
	})

	if err != nil {
		log.Panicln("Failed to Insert TraceEntry", err)
		return err
	}

	log.Println("Inserted TraceEntry Id=", res.InsertedID)

	return nil
}

func (m *Models) All() ([]*TraceEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := m.mongo.Database("trace").Collection("trace")

	opts := options.Find()
	opts.SetSort(bson.D{{"created_at", -1}})

	cursor, err := collection.Find(context.TODO(), bson.D{}, opts)

	if err != nil {
		m.logger.Println("Failed to find all entries", err)
		return nil, err
	}

	var all []*TraceEntry

	for cursor.Next(ctx) {
		var one TraceEntry

		err := cursor.Decode(&one)

		if err != nil {
			m.logger.Println("Failed to decode entry to slice", err)
			return nil, err
		}

		all = append(all, &one)
	}

	return all, nil
}

func (m *Models) GetOne(id string) (*TraceEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := m.mongo.Database("trace").Collection("trace")

	docId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		m.logger.Println("Failed to ObjectIDFromHex", err)
		return nil, err
	}

	var one TraceEntry

	err = collection.FindOne(ctx, bson.M{"_id": docId}).Decode(&one)

	if err != nil {
		m.logger.Println("Failed to ObjectIDFromHex", err)
		return nil, err
	}

	return &one, nil
}

func (m *Models) DropAll() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := m.mongo.Database("trace").Collection("trace")

	if err := collection.Drop(ctx); err != nil {
		return err
	}

	return nil
}

func (m *Models) Update(one TraceEntry) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := m.mongo.Database("trace").Collection("trace")

	docID, err := primitive.ObjectIDFromHex(one.ID)

	if err != nil {
		m.logger.Println("Failed to ObjectIDFromHex during update", err)
		return nil, err
	}

	result, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": docID},
		bson.D{
			{"$set", bson.D{
				{"src", one.Src},
				{"via", one.Via},
				{"data", one.Data},
				// {"updated_at", time.Now()},
			}},
		},
	)

	if err != nil {
		m.logger.Println("Failed to UpdateOne", err)
		return nil, err
	}

	return result, nil
}
