package main

import (
	"os"
	"context"
	"log"
	"strings"
	"fmt"

	"github.com/alecthomas/kingpin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	ADMIN = "admin"
	DB = "test_db"
	CollCoord = "user_coordinate"
	CollColor = "user_favcolor"
	CollEmail = "user_email"
)

type MTF struct {
	cl	*mongo.Client
	ctx	context.Context
}

var (
	mtfCmd = kingpin.New("mtf", "Traffic generator for MongoDB")

	mURI = mtfCmd.Flag("mongodb-uri", "MongoDB connection string").String()

	setCmd		= mtfCmd.Command("setup", "setup collection")
	setTypeF	= setCmd.Flag("type", "collection type").Default("1").Int64()

	insertCmd	= mtfCmd.Command("insert", "insert data")
	insertTypeF	= insertCmd.Flag("type", "collection type").Default("1").Int64()
)

func main() {
	cmd, err := mtfCmd.DefaultEnvars().Parse(os.Args[1:])
	if err != nil {
		log.Fatalln(err)
	}

	if *mURI == "" {
		log.Println("no mongodb connection URI")
		mtfCmd.Usage(os.Args[1:])
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mtfClient := connClient(ctx, *mURI)

	switch cmd {
	case setCmd.FullCommand():
		mtfClient.SetupCollection(*setTypeF)
		fmt.Println("Setup collection is done")
	case insertCmd.FullCommand():
		fmt.Printf("Insertion is running... [Type: %d]\n", *insertTypeF)
		mtfClient.InsertData(*insertTypeF)
	}
}

func connClient(ctx context.Context, mongoURI string) *MTF {
	uri := "mongodb://" + strings.Replace(mongoURI, "mongodb://", "", 1)

	client := connect(ctx, uri)

	mtf := &MTF {
		cl:		client,
		ctx:	ctx,
	}
	
	return mtf
}

func connect(ctx context.Context, uri string) *mongo.Client {
	client, err := mongo.NewClient(
		options.Client().ApplyURI(uri),
	)
	if err != nil {
		log.Fatalln(err)
	}
	
	err = client.Connect(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalln(err)
	}

	return client
}