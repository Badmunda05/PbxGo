package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type SessionEntry struct {
	UserID    int64  `bson:"user_id"`
	Session   string `bson:"session"`
	FirstName string `bson:"first_name"`
	AddedAt   time.Time `bson:"added_at"`
}

func AddSession(userID int64, session, firstName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filter := bson.D{{Key: "user_id", Value: userID}}
	update := bson.D{{Key: "$set", Value: SessionEntry{
		UserID:    userID,
		Session:   session,
		FirstName: firstName,
		AddedAt:   time.Now(),
	}}}
	_, err := Sessions.UpdateOne(ctx, filter, update, upsertOpt())
	return err
}

func RemoveSession(userID int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := Sessions.DeleteOne(ctx, bson.D{{Key: "user_id", Value: userID}})
	return err
}

func GetAllSessions() ([]SessionEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cursor, err := Sessions.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var result []SessionEntry
	for cursor.Next(ctx) {
		var e SessionEntry
		if err := cursor.Decode(&e); err == nil {
			result = append(result, e)
		}
	}
	return result, nil
}

func SessionExists(userID int64) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	count, _ := Sessions.CountDocuments(ctx, bson.D{{Key: "user_id", Value: userID}})
	return count > 0
}
