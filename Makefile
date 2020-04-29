deploy:
	gcloud functions deploy Translate --region=asia-northeast1 --runtime go113 --trigger-http --allow-unauthenticated

clean:
	gcloud functions delete Translate --region=asia-northeast1

test:
	go test github.com/mactynow/wordlistmaker