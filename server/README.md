
Installation
1. go get github.com/rakyll/statik
2. go get -u github.com/go-swagger/go-swagger/cmd/swagger #For downloading the Go Swagger to create the spec using the swagger comments.


Steps to follow
1. Add the Swagger Comments to the API added or updated using documention at https://goswagger.io/use/spec.html.
2. Run GO111MODULE=off swagger generate spec -o ./swagger.yaml --scan-models  from the root directory.
3. cd maticnetwork/heimdall/server
4. Replace the Swagger.yaml file inside swagger-ui directory with the swagger.yaml newly generated in root directly in step 2
5. cd maticnetwork/heimdall/server && statik -src=./swagger-ui
6. cd maticnetwork/heimdall && make build
7. cd maticnetwork/heimdall && make run-server

Visit http://localhost:1317/swagger-ui/ 


Reference
- https://github.com/rakyll/statik