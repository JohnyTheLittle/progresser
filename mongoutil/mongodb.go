package mongoutil

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/johnythelittle/goupdateyourself/configs"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var c, _ = configs.LoadConfig("../")
var URI string = c.DB_URI

func DB(cllctn string) *mongo.Collection {
	client, err := mongo.NewClient(options.Client().ApplyURI(URI))
	if err != nil {
		fmt.Println(err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	err = client.Connect(ctx)

	if err != nil {
		log.Panic(err)
	}

	return client.Database("tasker").Collection(cllctn)
}
