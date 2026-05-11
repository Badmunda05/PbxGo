package database

import (
	"context"
	"log/slog"
	"pbxgo/config"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var (
	mClient    *mongo.Client
	db         *mongo.Database
	Sessions   *mongo.Collection
	SudoColl   *mongo.Collection
)

func Init() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var err error
	mClient, err = mongo.Connect(options.Client().ApplyURI(config.App.MongoURL))
	if err != nil {
		slog.Error("MongoDB connect failed", "error", err)
		return
	}
	if err = mClient.Ping(ctx, nil); err != nil {
		slog.Error("MongoDB ping failed", "error", err)
		mClient = nil
		return
	}
	db = mClient.Database("pbxgo")
	Sessions = db.Collection("sessions")
	SudoColl = db.Collection("sudo_users")
	slog.Info("✅ MongoDB connected")
}

func IsConnected() bool { return mClient != nil }

func Disconnect() {
	if mClient != nil {
		_ = mClient.Disconnect(context.Background())
	}
}
