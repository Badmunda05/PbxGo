package database

import (
	"context"
	"log/slog"
	"sync"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// sudoUsers uses sync.Map for safe concurrent access (Go 1.24 best practice)
var sudoUsers sync.Map

type sudoEntry struct {
	UserID int64 `bson:"user_id"`
}

func LoadSudoUsers() {
	if !IsConnected() {
		slog.Warn("MongoDB not connected — skipping sudo load")
		return
	}

	ctx := context.Background()
	cursor, err := SudoCollection.Find(ctx, bson.D{})
	if err != nil {
		slog.Error("Failed to fetch sudo users", "error", err)
		return
	}
	defer cursor.Close(ctx)

	count := 0
	for cursor.Next(ctx) {
		var entry sudoEntry
		if err := cursor.Decode(&entry); err != nil {
			slog.Warn("Failed to decode sudo entry", "error", err)
			continue
		}
		sudoUsers.Store(entry.UserID, true)
		count++
	}

	slog.Info("Sudo users loaded", "count", count)
}

func AddSudo(userID int64) {
	sudoUsers.Store(userID, true)

	if IsConnected() {
		ctx := context.Background()
		// Upsert to avoid duplicate key errors
		filter := bson.D{{Key: "user_id", Value: userID}}
		update := bson.D{{Key: "$setOnInsert", Value: bson.D{{Key: "user_id", Value: userID}}}}
		opts := options.UpdateOne().SetUpsert(true)
		if _, err := SudoCollection.UpdateOne(ctx, filter, update, opts); err != nil {
			slog.Error("Failed to persist sudo add", "user_id", userID, "error", err)
		}
	}
}

func RemoveSudo(userID int64) {
	sudoUsers.Delete(userID)

	if IsConnected() {
		ctx := context.Background()
		if _, err := SudoCollection.DeleteOne(ctx, bson.D{{Key: "user_id", Value: userID}}); err != nil {
			slog.Error("Failed to persist sudo remove", "user_id", userID, "error", err)
		}
	}
}

func FetchSudoList() []int64 {
	var list []int64
	sudoUsers.Range(func(key, _ any) bool {
		list = append(list, key.(int64))
		return true
	})
	return list
}

func IsSudo(userID int64) bool {
	_, ok := sudoUsers.Load(userID)
	return ok
}
