package database

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
)

var SudoUsers = make(map[int64]bool)

func LoadSudoUsers() {
	if MongoClient == nil {
		log.Println("⚠️  MongoDB not connected, sudo users will be in-memory only")
		return
	}

	ctx := context.Background()
	cursor, err := SudoCollection.Find(ctx, bson.M{})
	if err != nil {
		log.Printf("❌ Failed to load sudo users: %v", err)
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var result struct {
			UserID int64 `bson:"user_id"`
		}
		err := cursor.Decode(&result)
		if err != nil {
			log.Printf("Failed to decode sudo user: %v", err)
			continue
		}
		SudoUsers[result.UserID] = true
	}

	if err := cursor.Err(); err != nil {
		log.Printf("Cursor error: %v", err)
	}

	log.Printf("✅ Loaded %d sudo users from database", len(SudoUsers))
}

func AddSudo(userID int64) {
	SudoUsers[userID] = true

	if MongoClient != nil {
		ctx := context.Background()
		_, err := SudoCollection.InsertOne(ctx, bson.M{"user_id": userID})
		if err != nil {
			log.Printf("Failed to add sudo user to database: %v", err)
		}
	}
}

func RemoveSudo(userID int64) {
	delete(SudoUsers, userID)

	if MongoClient != nil {
		ctx := context.Background()
		_, err := SudoCollection.DeleteOne(ctx, bson.M{"user_id": userID})
		if err != nil {
			log.Printf("Failed to remove sudo user from database: %v", err)
		}
	}
}

func FetchSudoList() []int64 {
	var list []int64
	for userID := range SudoUsers {
		list = append(list, userID)
	}
	return list
}

func IsSudo(userID int64) bool {
	return SudoUsers[userID]
}
