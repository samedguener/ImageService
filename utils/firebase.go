package utils

import (
	"context"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/sirupsen/logrus"
)

// InitFirebaseAuth ...
func InitFirebaseAuth() (*auth.Client, error) {

	// TODO check for a global app instance
	app, err := firebase.NewApp(context.Background(), nil)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	auth, err := app.Auth(context.Background())
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return auth, nil
}
