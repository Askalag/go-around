check_swagger:
	command -v swagger || go get -u github.com/go-swagger/go-swagger/cmd/swagger

swagger: check_swagger
	GO111MODULE=off swagger generate spec -o ./swagger.yaml --scan-models