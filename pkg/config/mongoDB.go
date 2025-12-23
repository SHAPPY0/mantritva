package config

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"reflect"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	Client                  *mongo.Client
	FeedCollection          *mongo.Collection
	NewsPressPostCollection *mongo.Collection
	PostReactionCollection  *mongo.Collection
	GrievanceCollection     *mongo.Collection
	Filter                  bson.M
}

var (
	connErr   error
	mongoOnce sync.Once
)

const (
	username = "mantritva024"
	password = "Mantritva#399!"
	host     = "localhost"
	port     = 27017
	dbName   = "mantritva"
)

func SetupMongoDB() (*MongoDB, error) {
	mongoDB := &MongoDB{}
	mongoOnce.Do(func() {
		uri := fmt.Sprintf("mongodb://%s:%s@%s:%d/%s?authSource=%s", username, url.QueryEscape(password), host, port, dbName, dbName)
		clientOptions := options.Client().ApplyURI(uri)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		client, err := mongo.Connect(ctx, clientOptions)
		if err != nil {
			connErr = err
			return
		}
		err = client.Ping(ctx, nil)
		if err != nil {
			connErr = err
			return
		}
		mongoDB.Client = client

		if err := mongoDB.SetupCollections(); err != nil {
			connErr = err
			return
		}
	})
	return mongoDB, connErr
}

func (m *MongoDB) SetupCollections() error {
	db := m.Client.Database(dbName)
	if err := createCollection(db, "feeds"); err != nil {
		return err
	} else {
		m.FeedCollection = db.Collection("feeds")
	}
	if err := createCollection(db, "news_press_posts"); err != nil {
		return err
	} else {
		m.NewsPressPostCollection = db.Collection("news_press_posts")
	}
	if err := createCollection(db, "post_reactions"); err != nil {
		return err
	} else {
		m.PostReactionCollection = db.Collection("post_reactions")
	}
	if err := createCollection(db, "grievances"); err != nil {
		return err
	} else {
		m.GrievanceCollection = db.Collection("grievances")
	}
	return nil
}

func createCollection(db *mongo.Database, collectionName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Check if collection already exists
	collections, err := db.ListCollectionNames(ctx, map[string]interface{}{})
	if err != nil {
		return err
	}

	for _, c := range collections {
		if c == collectionName {
			log.Println("ℹ️ Collection already exists:", collectionName)
			return nil
		}
	}

	// Create collection if not exists
	err = db.CreateCollection(ctx, collectionName)
	if err != nil {
		return err
	}

	log.Println("✅ Created collection:", collectionName)
	return nil
}

func (m *MongoDB) Insert(collection *mongo.Collection, data any) (*mongo.InsertOneResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := collection.InsertOne(ctx, data)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (m *MongoDB) InsertMany(collection *mongo.Collection, data []any) (*mongo.InsertManyResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := collection.InsertMany(ctx, data)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (m *MongoDB) FindOne(collection *mongo.Collection, filter bson.M, row interface{}) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := collection.FindOne(ctx, filter).Decode(row); err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (m *MongoDB) FindMany(collection *mongo.Collection, filter bson.M, opts *options.FindOptions, rows interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	sliceValue := reflect.ValueOf(rows).Elem()
	elemType := sliceValue.Type().Elem()
	for cursor.Next(ctx) {
		elem := reflect.New(elemType).Interface()
		if err := cursor.Decode(elem); err != nil {
			return err
		}
		sliceValue.Set(reflect.Append(sliceValue, reflect.ValueOf(elem).Elem()))
	}
	return nil
}

func (m *MongoDB) UpdateOne(collection *mongo.Collection, filter bson.M, update bson.M) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}
	return result, nil
}
