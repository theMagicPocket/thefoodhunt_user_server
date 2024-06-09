package tokenverification
import (
	"context"
	"fmt"
	"os"
	"log"
   
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
   )
   
   func ValidateFirebaseToken(token string) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "validate failed", fmt.Errorf("error getting current working directory: %v", err)
	}
	fmt.Println("Current working directory:", wd)
	serviceAccountJSON, err := os.ReadFile("/home/teja/Documents/teja2.0/golang_food_delivery_api/firebase.json")
	if err != nil {
	 return "validate failed", fmt.Errorf("error reading service account file: %v", err)
	}
	app, err := firebase.NewApp(context.Background(), nil, option.WithCredentialsJSON(serviceAccountJSON))
	if err != nil {
	 return "validate failed", fmt.Errorf("error initializing Firebase app: %v", err)
	}
	client, err := app.Auth(context.Background())
	if err != nil {
	 return "validate failed", fmt.Errorf("error getting Auth client: %v", err)
	}
	tokenInfo, err := client.VerifyIDToken(context.Background(), token)
	if err != nil {
	 return "validate failed", fmt.Errorf("error verifying ID token: %v", err) }

	 log.Printf("Verified ID token: %v\n", tokenInfo)
	 return "validate success", nil
	}