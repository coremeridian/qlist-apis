// Package mongodb instantiates and communicates with database server
package mongodb

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"coremeridian.xyz/app/qlist/cmd/api/services/tests/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var db *mongo.Database

var (
	ErrNoConnection = errors.New("no connection")
	ErrNoInsertion  = errors.New("data not inserted")
)

func GetDB(database string) (db *mongo.Database, err error) {
	if db != nil {
		return
	}

	vault, err := utils.GetAPIVault()
	if err != nil {
		return
	}

	credentials := options.Credential{
		Username: vault["username"].(string),
		Password: vault["password"].(string),
	}

	uri := os.Getenv("DATABASE_URL")
	clientOpts := options.Client().ApplyURI(uri).SetAuth(credentials)
	client, err := mongo.NewClient(clientOpts)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return
	}

	fmt.Println("Connected to MongoDB")
	db = client.Database(database)
	return
}
