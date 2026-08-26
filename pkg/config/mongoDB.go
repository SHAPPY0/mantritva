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
	ReactionCollection      *mongo.Collection
	GrievanceCollection     *mongo.Collection
	GrievanceCosignCollection *mongo.Collection
	ScoreCollection         *mongo.Collection
	ScoreSnapshotCollection *mongo.Collection
	EventCollection         *mongo.Collection
	DevelopmentWelfare      *mongo.Collection
	AnalyticStatsCollection *mongo.Collection
	Filter                  bson.M
}

var (
	connErr   error
	mongoOnce sync.Once
)

func SetupMongoDB() (*MongoDB, error) {
	var (
		username = GetConfig().GetString("MongoDB_USER")
		password = GetConfig().GetString("MongoDB_PASSWORD")
		host     = GetConfig().GetString("MongoDB_HOST")
		port     = GetConfig().GetInt("MongoDB_PORT")
		dbName   = GetConfig().GetString("MongoDB_DB_NAME")
	)
	mongoDB := &MongoDB{}
	mongoOnce.Do(func() {
		uri := fmt.Sprintf("mongodb://%s:%s@%s:%d/%s?authSource=admin", username, url.QueryEscape(password), host, port, dbName)
		if GetConfig().GetString("ENV") == "PROD" {
			// uri = fmt.Sprintf("mongodb+srv://%s:%s@%s/?appName=mantritva&authSource=admin", username, url.QueryEscape(password), host)
			uri = fmt.Sprintf("mongodb://%s:%s@%s:%d/%s?authSource=admin", username, url.QueryEscape(password), host, port, dbName)
		}

		log.Print("MongoDB Connect URI: " + uri)
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

		if err := mongoDB.SetupCollections(dbName); err != nil {
			connErr = err
			return
		}
	})
	return mongoDB, connErr
}

func (m *MongoDB) SetupCollections(dbName string) error {
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

	if err := createCollection(db, "reactions"); err != nil {
		return err
	} else {
		m.ReactionCollection = db.Collection("reactions")
	}

	if err := createCollection(db, "grievances"); err != nil {
		return err
	} else {
		m.GrievanceCollection = db.Collection("grievances")
	}

	if err := createCollection(db, "grievance_cosigns"); err != nil {
		return err
	} else {
		m.GrievanceCosignCollection = db.Collection("grievance_cosigns")
	}

	if err := createCollection(db, "scores"); err != nil {
		return err
	} else {
		m.ScoreCollection = db.Collection("scores")
	}

	if err := createCollection(db, "score_snapshots"); err != nil {
		return err
	} else {
		m.ScoreSnapshotCollection = db.Collection("score_snapshots")
	}

	if err := createCollection(db, "events"); err != nil {
		return err
	} else {
		m.EventCollection = db.Collection("events")
	}

	if err := createCollection(db, "development_welfares"); err != nil {
		return err
	} else {
		m.DevelopmentWelfare = db.Collection("development_welfares")
	}

	if err := createCollection(db, "analytic_stats"); err != nil {
		return err
	} else {
		m.AnalyticStatsCollection = db.Collection("analytic_stats")
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

func (m *MongoDB) Upsert(collection *mongo.Collection, data any) (*mongo.InsertOneResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Try to marshal the incoming data to BSON map to detect an _id field.
	var doc bson.M
	raw, err := bson.Marshal(data)
	if err != nil {
		return nil, err
	}
	if err := bson.Unmarshal(raw, &doc); err != nil {
		return nil, err
	}

	// If document contains an _id, perform a ReplaceOne with upsert=true so existing
	// documents are replaced and missing ones are inserted.
	if id, ok := doc["_id"]; ok && id != nil {
		filter := bson.M{"_id": id}
		opts := options.Replace().SetUpsert(true)
		updRes, err := collection.ReplaceOne(ctx, filter, data, opts)
		if err != nil {
			return nil, err
		}
		// If an upsert happened, return an InsertOneResult-like value with the UpsertedID.
		if updRes.UpsertedCount > 0 {
			return &mongo.InsertOneResult{InsertedID: updRes.UpsertedID}, nil
		}
		// No insert (existing document replaced) — return a nil InsertedID to indicate update.
		return &mongo.InsertOneResult{InsertedID: nil}, nil
	}

	// No _id present — fall back to InsertOne.
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

func (m *MongoDB) FindOne(collection *mongo.Collection, filter bson.M, opts *options.FindOneOptions, row interface{}) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := collection.FindOne(ctx, filter, opts).Decode(row); err != nil {
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

func (m *MongoDB) DeleteOne(collection *mongo.Collection, filter bson.M) (*mongo.DeleteResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return nil, err
	}
	return result, nil
}
