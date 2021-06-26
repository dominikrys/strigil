package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func TestMain(t *testing.T) {
	// Setup: clear database (ideally this should be in a separate function)
	mongoUri := "mongodb://localhost:27017"
	mongoClient, usingMongo := connectToMongo(mongoUri)
	if !usingMongo {
		t.Fatalf("Couldn't connect to MongoDB")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := mongoClient.Database("crawler").Collection("profiles")

	if err := collection.Drop(ctx); err != nil {
		t.Fatalf("Error clearing the database: %v", err)
	}

	// Arrange
	profileNo := 3
	args := fmt.Sprintf("--day 1 --month 2 --profileNo %d", profileNo)

	os.Args = append(os.Args, strings.Fields(args)...)

	// Act
	main()

	// Assert
	countRes, err := collection.CountDocuments(ctx, bson.D{})

	if err != nil {
		t.Fatalf("Error counting documents in database: %v", err)
	}

	if countRes != int64(profileNo) {
		t.Fatalf("Expected %d profiles in database, but found %d", profileNo, countRes)
	}
}
