package util

import (
	"context"
	"fmt"
	"os"
	"sync"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"google.golang.org/api/option"
)

type FirebaseApp struct {
	Auth      *auth.Client
	Firestore *firestore.Client
}

var (
	firebaseApp  *FirebaseApp
	firebaseOnce sync.Once
)

func getFirebaseApp() *FirebaseApp {
	firebaseOnce.Do(func() {

		config := GetConfig()

		ctx := context.Background()
		opt := option.WithServiceAccountFile(config.Firebase.Credentials)
		app, err := firebase.NewApp(ctx, nil, opt)
		if err != nil {
			fmt.Printf("[BOOTING] error when initiating firebase: %v\n", err)
			os.Exit(1)
			return
		}

		firebaseApp = new(FirebaseApp)

		firebaseApp.Auth, err = app.Auth(ctx)
		if err != nil {
			fmt.Printf("[BOOTING] error when initiating firebase Auth: %v\n", err)
			os.Exit(1)
			return
		}

		firebaseApp.Firestore, err = app.Firestore(ctx)
		if err != nil {
			fmt.Printf("[BOOTING] error when initiating firebase Firestore: %v\n", err)
			os.Exit(1)
			return
		}

		fmt.Println("[BOOTING] successfully connected to Firebase!")
	})

	return firebaseApp
}
