deploy:
	gcloud functions deploy Translate --region=asia-northeast1 --runtime go113 --trigger-http

clean:
	gcloud functions delete Translate --region=asia-northeast1

test:
	go test -cover github.com/mactynow/wordlistmaker