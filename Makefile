.PHONY: deploy delete

deploy:
	gcloud app deploy app.yaml --format json --no-promote -q

promote:
	gcloud app deploy app.yaml --format json -q

test:
	go test .

integration: # TODO add integration tests

smoke: # TODO add smoke tests
