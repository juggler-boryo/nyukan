package lib

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/db"
	"google.golang.org/api/option"
)

var (
	DB   *firestore.Client
	RTDB *db.Client
)

func InitializeFirebase(credPath string) error {
	ctx := context.Background()
	opt := option.WithCredentialsFile(credPath)

	conf := &firebase.Config{
		DatabaseURL:   "https://aigrid-23256-default-rtdb.asia-southeast1.firebasedatabase.app",
		ProjectID:     "aigrid-23256",
		StorageBucket: "aigrid-23256.firebasestorage.app",
	}

	app, err := firebase.NewApp(ctx, conf, opt)
	if err != nil {
		return fmt.Errorf("error initializing firebase app: %w", err)
	}

	DB, err = app.Firestore(ctx)
	if err != nil {
		return fmt.Errorf("error initializing firestore: %w", err)
	}

	RTDB, err = app.Database(ctx)
	if err != nil {
		return fmt.Errorf("error initializing realtime database client: %w", err)
	}

	if RTDB == nil {
		return fmt.Errorf("realtime database client is nil after initialization")
	}

	return nil
}
