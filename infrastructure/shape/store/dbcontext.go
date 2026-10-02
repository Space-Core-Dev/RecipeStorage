package store

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

var mClient *mongo.Client
var m sync.RWMutex

func InitDb(ctx context.Context) error {
	var err error
	uri := os.Getenv("CONNECTION_STRING")
	if uri == "" {
		return errors.New("wrong connection string")
	}

	m.Lock()
	defer m.Unlock()
	if mClient != nil {
		return errors.New("the connection already exists")
	}
	mClient, err = mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		return fmt.Errorf("db connection error: %w", err)
	}

	if err = mClient.Ping(ctx, readpref.Primary()); err != nil {
		mClient = nil
		return fmt.Errorf("db connection error: %w", err)
	}

	return err
}

func Open(ctx context.Context, dbName, colName string) (*mongo.Collection, error) {
	m.RLock()
	defer m.RUnlock()
	if mClient == nil {
		return nil, errors.New("there is no active db client")
	}

	coll := mClient.Database(dbName).Collection(colName)
	return coll, nil
}

func CloseConnection(ctx context.Context) error {
	m.Lock()
	defer m.Unlock()
	if mClient == nil {
		return errors.New("there is no active db client")
	}

	if err := mClient.Disconnect(ctx); err != nil {
		return err
	}

	mClient = nil

	return nil
}
