# CarServiceCenter
An API that accepts http requests for creating, updating, deleting, and reading appointment data.

It supports the following request(s):

```GET /appointment/{id}```

```POST /appointment/```

```PATCH /appointment/{id}```

```GET /appointments/range/```

```DELETE /appointment/{id}```



## Prerequisites

Be sure to have Go installed locally.

Have A local instance of MongoDB installed on your machine.

## Installation

Clone the project with ```git clone https://github.com/marshacb/CarServiceCenter.git``` in your Go workspace. Then go into the project directory and install all dependencies:

```
cd CarServiceCenter

#Install the dependencies with go get
go get -d ./...

#Install test dependencies
go get -t ./...
```

## Running the tests

From the root project directory run

```go test ./...```

Alternatively cd into src/controller/ and run

```go test```

## Running the server

From the root project directory run

```go run main.go```

You can also run

```go build main.go``` followed by ```./main``` in order to start the server.

## Example Create Request

```curl -d '{"Name": "Ultimate Car Appointment", "Description": "even newer engine appointment", "Status": "open", "Date": "2019-08-28T09:00:01+00:00"}' -H "Content-Type: application/json" -X POST http://localhost:8080/appointment/ ```

## Example GetDateWithinRange Request 

```curl -X GET 'http://localhost:8080/appointments/range/?start=2019-07-29T09:00:01+00:00&end=2019-08-29T09:00:01+00:00'```

# Example UpdateAppointmentStatus Request

```curl -d '{"status": "closed"}' -H "Content-Type: application/json" -X PATCH http://localhost:8080/appointment/{id}```