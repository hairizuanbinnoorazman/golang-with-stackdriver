linux:
	env GOOS=linux GOARCH=amd64 go build -gcflags=all='-N -l' -o app .
sample:
	curl localhost:8888
