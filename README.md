# MySQL Tuning API

This Go program implements a simple REST API that calculates MySQL tuning parameters based on the provided flavor details.
It listens on port 8080.

## API Endpoint
### Get Tuning Parameters

- Endpoint: /tune
- Method: POST
- Request Body:
```
{
  "memory_gb": 4,
  "vcpus": 2,
  "db_type": "mysql",
  "db_version": "8.0"
}
```

- Response:
```
{
  "max_connections": 100,
  "innodb_buffer_pool_size": "2G",
  "innodb_dedicated_server": "ON"
  // Add other tuning parameters here
}
```

### Usage

1. Install Go: https://golang.org/doc/install
2. Clone this repository.
3. Run the program: `go run main.go`
4. Send a POST request to http://localhost:8080/tune with the flavor details in the request body. 
You can use tools like curl or Postman for testing.

>Note: The program calculates tuning parameters based on the provided flavor, but you can customize the calculation logic and add more tuning parameters as needed.

- Example Test:
```
curl -X POST -H "Content-Type: application/json" -d '{
  "flavor": "some_flavor",
  "memory": "4 GB",
  "vcpus": 2,
  "dbType": "mysql",
  "dbVersion": "8.0"
}' http://localhost:8080/tune
```