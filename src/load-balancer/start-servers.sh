go run main.go
go run ../backendserver/main.go --port=8001
go run ../backendserver/main.go --port=8002
curl --location --request POST 'localhost:8000/register' \
--header 'backend: http://localhost:8001'
curl --location --request POST 'localhost:8000/register' \
--header 'backend: http://localhost:8002'