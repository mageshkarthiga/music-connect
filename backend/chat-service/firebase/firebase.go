package firebase

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
)

var FirestoreClient *firestore.Client

func InitFirebase() {
	ctx := context.Background()

	// Automatically uses ADC: 
	// - Locally: from `gcloud auth application-default login`
	// - Cloud Run: from default service account
	app, err := firebase.NewApp(ctx, nil)
	if err != nil {
		log.Fatalf("error initializing Firebase app: %v\n", err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalf("error initializing Firestore: %v\n", err)
	}

	FirestoreClient = client
	log.Println("Firestore connected using Application Default Credentials (ADC)")
}
