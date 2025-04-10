package chat

import (
	"context"
	"log"
	"os"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

var FirestoreClient *firestore.Client

func InitFirebase() {
	ctx := context.Background()

	// Get the credentials file path from an environment variable
	credentialsFile := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if credentialsFile == "" {
		log.Fatalf("GOOGLE_APPLICATION_CREDENTIALS environment variable is not set")
	}

	opt := option.WithCredentialsFile(credentialsFile)

	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalf("error initializing Firestore: %v\n", err)
	}

	FirestoreClient = client
	log.Println("Firebase Firestore connected")
}
