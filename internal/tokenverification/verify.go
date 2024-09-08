package tokenverification

import (
	"context"
	"fmt"
	"log"
	"os"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

func ValidateFirebaseToken(token string) (string, string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "extraction of userid failed", "validate failed", fmt.Errorf("error getting current working directory: %v", err)
	}
	fmt.Println("Current working directory:", wd)
	serviceAccountJSON, err := os.ReadFile("firebase.json")
	if err != nil {
		return "extraction of userid failed", "validate failed", fmt.Errorf("error reading service account file: %v", err)
	}
	app, err := firebase.NewApp(context.Background(), nil, option.WithCredentialsJSON(serviceAccountJSON))
	if err != nil {
		return "extraction of userid failed", "validate failed", fmt.Errorf("error initializing Firebase app: %v", err)
	}
	client, err := app.Auth(context.Background())
	if err != nil {
		return "extraction of userid failed", "validate failed", fmt.Errorf("error getting Auth client: %v", err)
	}
	tokenInfo, err := client.VerifyIDToken(context.Background(), token)
	if err != nil {
		return "extraction of userid failed", "validate failed", fmt.Errorf("error verifying ID token: %v", err)
	}

	log.Printf("Verified ID token: %v\n", tokenInfo)
	userID := tokenInfo.UID

	log.Printf("User ID: %s\n", userID)
	return userID, "validate success", nil
}
