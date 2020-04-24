package services

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"time"

	"github.com/google/uuid"
	"github.com/samedguener/ImageService/errors"
	"github.com/samedguener/ImageService/utils"
	"github.com/sirupsen/logrus"
)

// Images ...
var Images = images{}

type images struct {
}

func (i images) UploadImage(ctx context.Context, image multipart.File, imageType string) (string, error) {
	bucket, err := utils.GetBucket()
	if err != nil {
		return "", err
	}

	fileExtension := ".jpeg"
	switch imageType {
	case "image/png":
		fileExtension = ".png"
		break
	case "image/jpg":
		fileExtension = ".jpg"
		break
	case "image/jpeg":
		fileExtension = ".jpeg"
		break
	default:
		return "", errors.Internal.Newf("unknown image type '%s'", imageType)
	}
	uuid := uuid.New()
	filename := fmt.Sprintf("%s.%s", uuid.String(), fileExtension)

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, image); err != nil {
		return "", errors.UnprocessableEntity.Wrap(err, "unable to read image into buffer")
	}

	timeout, err := time.ParseDuration(utils.ImageUploadToGCPTimeout.Value)
	if err != nil {
		return "", errors.Internal.Wrapf(err, "unable to parse timeout '%s'", utils.ImageUploadToGCPTimeout.Value)
	}
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	fileWriter := bucket.Object(filename).NewWriter(ctx)
	if _, err = io.Copy(fileWriter, buf); err != nil {
		return "", errors.Internal.Wrap(err, "could not upload image")
	}
	if err := fileWriter.Close(); err != nil {
		return "", errors.Internal.Wrap(err, "could not upload image")
	}

	logrus.Infof("image with name '%s' uploaded", filename)
	return filename, nil
}

func (i images) DeleteImage(hash string) error {
	return nil
}
