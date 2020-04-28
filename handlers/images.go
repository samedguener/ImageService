package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/samedguener/ImageService/dtos"
	"github.com/samedguener/ImageService/errors"
	"github.com/samedguener/ImageService/services"
	"github.com/samedguener/ImageService/utils"
	"github.com/sirupsen/logrus"
)

// Images ..
var Images = images{}

type images struct {
	handlers
}

// GetImageAccessEndpoint ...
func (i images) GetImageAccessEndpoint(w http.ResponseWriter, r *http.Request) {
	var imageAccessEndpoint dtos.ImageAccessEndpointResponse = dtos.ImageAccessEndpointResponse{}
	imageAccessEndpoint.Endpoint = utils.ImageAccessEndpoint.Value

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(imageAccessEndpoint)
}

// Get ...
func (i images) Post(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	r.ParseMultipartForm(32 << 20)

	logrus.Info(r.FormValue("image"))
	file, _, err := r.FormFile("image")
	if err != nil {
		err = errors.BadRequestHTTP.Wrap(err, "could not read image from request")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	fullFile := bytes.NewBuffer(nil)
	if _, err := io.Copy(fullFile, file); err != nil {
		err = errors.UnprocessableEntityHTTP.Wrap(err, "could not read image from request")
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	fileHeader := fullFile.Bytes()[:512]
	imageType := http.DetectContentType(fileHeader)
	if imageType != "image/png" && imageType != "image/jpg" && imageType != "image/jpeg" {
		err = errors.UnprocessableEntityHTTP.New("unsupported image type")
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	imageID, err := services.Images.UploadImage(ctx, fullFile.Bytes(), imageType)
	if err != nil {
		switch errors.GetType(err) {
		case errors.Internal:
			http.Error(w, err.Error(), http.StatusInternalServerError)
			break
		case errors.UnprocessableEntity:
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			break
		}
		return
	}

	var response dtos.ImageUploadResponse
	response.ID = imageID

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Delete ...
func (i images) Delete(w http.ResponseWriter, r *http.Request) error {

	return nil
}
