package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/samedguener/ImageService/utils"
)

// VerifyToken ...
func VerifyToken(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var ctx context.Context

		if utils.AuthenticationMethod.Value == "firebase" {
			reqToken := r.Header.Get("Authorization")

			splitToken := strings.Split(reqToken, "Bearer")
			if len(splitToken) != 2 {
				err := fmt.Errorf("unauthorized or malformed token")
				logrus.Error(err)
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			reqToken = strings.Trim(splitToken[1], " ")

			auth, err := utils.InitFirebaseAuth()
			if err != nil {
				logrus.Error(err)
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			authToken, err := auth.VerifyIDToken(r.Context(), reqToken)
			if err != nil {
				err := fmt.Errorf("unauthorized or malformed token")
				logrus.Printf(err.Error())
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			role, ok := authToken.Claims["role"]
			if !ok {
				err := fmt.Errorf("registration not completed for %s", authToken.UID)
				logrus.Printf(err.Error())
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			if role != "api" {
				err := fmt.Errorf("unauthorized access of %s", authToken.UID)
				logrus.Printf(err.Error())
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			logrus.Printf("api authenticated %s with role %s", authToken.UID, role)

		} else if utils.AuthenticationMethod.Value == "env" {
			logrus.Info("authentication method is set to 'env'")
		} else {
			logrus.Errorf("no authentication method in env found")
			http.Error(w, "no authentication method in env found", http.StatusInternalServerError)
			return
		}

		r = r.WithContext(ctx)

		h.ServeHTTP(w, r)
	})
}
