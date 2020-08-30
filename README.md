![master](https://github.com/samedguener/ImageService/workflows/master/badge.svg)
# ImageService
This simple image service is written in Golang and can be currently only used Google Cloud Platform. You can easily deploy this application to Google App Engine. It will consume Google Cloud Bucket to manage uploaded images.

## Repo Structure
This project is structured as following:

|Folder | Description |
|---|---|
|dtos| data transfer objects - definitions for transforming models for sending and receiving |
|errors| HTTP and application errors|
|handlers| handler for endpoints |
|middleware| authentication, etc. |
|scripts| CI/CD related scripts |
|services| service implementation details |
|utils| helper functions e.g reading environment variables |

## Functionality

- token-based (JWT) authentication with Google Firebase
- upload images
- retrieve images from link
- delete images

## Available Endpoints

`api/v1/health`: `GET` a  `200` status code if service is able to serve.

`api/v1/images/endpoint`: `GET` the image access endpoint from which images can be retrieved.

`api/v1/images`: `POST` an image as `form-data` with key `image`. Returns a `201` status code and following payload if the image upload was successful:
```json
{
  "id" : "adh12ahajcs/jpg"
}
```
`api/vi/images/<image-id>`: `DELETE` an image with image id `image-id` as parameters. Returns a `204` status code if the image was deleted successfully.


## Configuration
Before starting a few environment variables needs to be set

|Name| Description | Default |
|---|---|---|
|`ENVIRONMENT`| _`REQUIRED`_ The name of environment. Setting to `dev` will set log level to verbose (debug enabled), otherwise JSON formated log.  | _no default_ |
|`BUCKET_NAME`| _`REQUIRED`_ The name of the Google Cloud Bucket to store images to.  | _no default_ |
|`PROJECT_ID`| _`REQUIRED`_ The name Google Project ID where service is deployed to.  | _no default_ |
|`IMAGE_ACCESS_ENDPOINT`| _`REQUIRED`_ The endpoint where image can be accessed. This is pointing to your Google Cloud Storage. This requires that your bucket is publicly accessible for read. | _no default_ |
|`AUTHENTICATION_METHOD`|  The name of the authentication method. Options: `env` for no authentication or `firebase` for Google Firebase authentication  | _firebase_ |
|`TIMEOUT_IMAGE_UPLOAD_GCP`| The uploading timeout of the service to Google Cloud Storage .  | _2s_ |
|`GOOGLE_APPLICATION_CREDENTIALS`| Path to the service account credential file. `REQUIRED` only if environment is set to `dev`  | _no default_ |

## Contribution

Feel free to open an issue or a PR ! Any contributions are welcome!

## Planned Functionality

Due to my limited time, I am unable to provide more features for this service. But still, feel free to extend this service on your own!
