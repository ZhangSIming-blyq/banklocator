IMG ?= localimage/bankmap:2.0.0
build: 
	cd cmd && CGO_ENABLE=0 GOOS=linux GOARCH=amd64 go build -a -o ./bankmap  
docker-build:
	docker build -f deploy/Dockerfile . -t ${IMG}
docker-save:
	docker save ${IMG} > deploy/imagefile
docker-push:
	docker push ${IMG}
