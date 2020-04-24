package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/samedguener/ImageService/dtos"
	"github.com/samedguener/ImageService/errors"
	"github.com/samedguener/ImageService/services"
	"github.com/samedguener/ImageService/utils"
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
	json.NewEncoder(w).Encode(imageAccessEndpoint)
}

// Get ...
func (i images) Post(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	r.ParseMultipartForm(32 << 20)

	file, _, err := r.FormFile("image")
	if err != nil {
		err = errors.BadRequestHTTP.New("could not read image from request")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileHeader := make([]byte, 512)
	if _, err := file.Read(fileHeader); err != nil {
		err = errors.UnprocessableEntityHTTP.New("could not read image from request")
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	imageType := http.DetectContentType(fileHeader)
	if imageType != "image/png" && imageType != "image/jpg" && imageType != "image/jpeg" {
		err = errors.UnprocessableEntityHTTP.New("unsupported image type")
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	imageID, err := services.Images.UploadImage(ctx, file, imageType)
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

	json.NewEncoder(w).Encode(response)
}

// Delete ...
func (i images) Delete(w http.ResponseWriter, r *http.Request) error {

	return nil
}
