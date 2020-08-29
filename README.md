![master](https://github.com/samedguener/ImageService/workflows/master/badge.svg)
# ImageService
This simple image service is written in Golang and can be currently only used Google Cloud Platform. You can easily deploy this application to Google App Engine. It will consume Google Cloud Bucket to manage uploaded images.

## Functionality

- token-based (JWT) authentication with Google Firebase
- upload images
- retrieve images from link
- delete images

## Configuration
Before starting a few environment variables needs to be set

|Name| Description | Default |
|---|---|---|
|`ENVIRONMENT`| _`REQUIRED`_ The name of environment. Setting to `dev` will set log level to verbose (debug enabled), otherwise JSON formated log.  | _no default_ |
|`BUCKET_NAME`| _`REQUIRED`_ The name of the Google Cloud Bucket to store images to.  | _no default_ |
|`PROJECT_ID`| _`REQUIRED`_ The name Google Project ID where service is deployed to.  | _no default_ |
|`IMAGE_ACCESS_ENDPOINT`| _`REQUIRED`_ The endpoint where image can be accessed. This is pointing to your Google Cloud Storage. This requires that your bucket is publicly accessible for read. | _no default_ |
|`AUTHENTICATION_METHOD`|  The name of the authentication method. Options: `env` for no authentication or `firebase` for Google Firebase authentication  | _firebase_ |
|`TIMEOUT_IMAGE_UPLOAD_GCP`| The uploading timeout of the service to Google .  | _no default_ |
|`GOOGLE_APPLICATION_CREDENTIALS`| _`REQUIRED`_ The name of environment. Setting to `dev` will set log level to verbose (debug enabled), otherwise JSON formated log.  | _no default_ |
