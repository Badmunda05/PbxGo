package database

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var sudoUsers sync.Map

type sudoEntry struct {
	UserID int64 `bson:"user_id"`
}

func upsertOpt() *options.UpdateOneOptionsBuilder {
	return options.UpdateOne().SetUpsert(true)
}

func LoadSudoUsers() {
	if !IsConnected() {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cursor, err := SudoColl.Find(ctx, bson.D{})
	if err != nil {
		return
	}
	defer cursor.Close(ctx)
	count := 0
	for cursor.Next(ctx) {
		var e sudoEntry
		if err := cursor.Decode(&e); err == nil {
			sudoUsers.Store(e.UserID, true)
			count++
		}
	}
	slog.Info("Sudo users loaded", "count", count)
}

func AddSudo(userID int64) {
	sudoUsers.Store(userID, true)
	if IsConnected() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		filter := bson.D{{Key: "user_id", Value: userID}}
		update := bson.D{{Key: "$setOnInsert", Value: bson.D{{Key: "user_id", Value: userID}}}}
		_, _ = SudoColl.UpdateOne(ctx, filter, update, upsertOpt())
	}
}

func RemoveSudo(userID int64) {
	sudoUsers.Delete(userID)
	if IsConnected() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_, _ = SudoColl.DeleteOne(ctx, bson.D{{Key: "user_id", Value: userID}})
	}
}

func IsSudo(userID int64) bool {
	_, ok := sudoUsers.Load(userID)
	return ok
}

func FetchSudoList() []int64 {
	var list []int64
	sudoUsers.Range(func(k, _ any) bool {
		list = append(list, k.(int64))
		return true
	})
	return list
}
