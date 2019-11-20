
Installation
1. go get github.com/rakyll/statik


Steps to follow
1. cd maticnetwork/heimdall/server
2. Update swagger.yaml file inside swagger-ui directory
3. cd maticnetwork/heimdall/server && statik -src=./swagger-ui
4. cd maticnetwork/heimdall && make build
5. cd maticnetwork/heimdall && make run-server

Visit http://localhost:1317/swagger-ui/ 


Reference
- https://github.com/rakyll/statik