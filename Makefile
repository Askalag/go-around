check_swagger:
	command -v swagger || go get -u github.com/go-swagger/go-swagger/cmd/swagger

swagger_prepare:
	$(shell mkdir -p "api-docs")

swagger: check_swagger swagger_prepare
	GO111MODULE=off swagger generate spec -o ./api-docs/swagger.yaml --scan-models