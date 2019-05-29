linux: 
	env GOOS=linux GOARCH=amd64 go build -o app .
debug-linux:
	env GOOS=linux GOARCH=amd64 go build -gcflags '-N -l' -o app .
sample:
	curl localhost:8888
debug-context:
	gcloud debug source gen-repo-info-file --output-directory .