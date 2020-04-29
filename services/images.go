package services

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/storage"
	"github.com/google/uuid"
	"github.com/samedguener/ImageService/errors"
	"github.com/samedguener/ImageService/utils"
	"github.com/sirupsen/logrus"
)

// Images ...
var Images = images{}

type images struct {
}

func (i images) UploadImage(ctx context.Context, image []byte, imageType string) (string, error) {
	bucket, err := utils.GetBucket()
	if err != nil {
		return "", err
	}

	var fileExtension string
	switch imageType {
	case "image/png":
		fileExtension = "png"
		break
	case "image/jpg":
		fileExtension = "jpg"
		break
	case "image/jpeg":
		fileExtension = "jpeg"
		break
	default:
		return "", errors.Internal.Newf("unknown image type '%s'", imageType)
	}
	uuid := uuid.New()
	filename := fmt.Sprintf("%s.%s", uuid.String(), fileExtension)

	timeout, err := time.ParseDuration(utils.ImageUploadToGCPTimeout.Value)
	if err != nil {
		return "", errors.Internal.Wrapf(err, "unable to parse timeout '%s'", utils.ImageUploadToGCPTimeout.Value)
	}
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	object := bucket.Object(filename)
	fileWriter := object.NewWriter(ctx)

	if _, err = fileWriter.Write(image); err != nil {
		return "", errors.Internal.Wrap(err, "could not upload image")
	}
	if err := fileWriter.Close(); err != nil {
		return "", errors.Internal.Wrap(err, "could not upload image")
	}

	acl := object.ACL()
	acl.Set(ctx, storage.AllUsers, storage.RoleReader)

	logrus.Infof("image with name '%s' uploaded", filename)
	return filename, nil
}

func (i images) DeleteImage(ctx context.Context, id string) error {
	bucket, err := utils.GetBucket()
	if err != nil {
		return err
	}

	object := bucket.Object(id)
	_, err = object.Attrs(ctx)
	if err != nil {
		switch err {
		case storage.ErrObjectNotExist:
			return errors.NotFound.Newf("image with id '%s' not found", id)
		default:
			return errors.Internal.Wrapf(err, "internal server error")
		}
	}
	object.Delete(context.Background())
	if err != nil {
		return errors.Internal.Wrapf(err, "internal server error")
	}
	logrus.Infof("image with id '%s' successfully deleted", id)
	return nil
}
