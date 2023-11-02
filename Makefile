dev:
	SET GOOS=windows&&SET GO_ENV=development&&air

swagger:
	swag init --dir ./,./handlers/

build:
	SET GOOS=linux&&SET GOARCH=amd64&&go build