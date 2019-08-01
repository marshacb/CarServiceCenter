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

Clone the project with ```git clone https://github.com/marshacb/StreamsAPIChallenge.git``` in your Go workspace. Then go into the project directory and install all dependencies:

```
cd StreamsAPIChallenge 

#Install the dependencies with go get
go get -d ./...

#Install test dependencies
go get -t ./...
```

Run the following in the command line from the root project directory to import the exported data into a test database:

```mongoimport --db test --collection streams --file streams.json --jsonArray```

## Running the tests

From the root project directory run

```ginkgo -r```

## Running the server

From the root project directory run

```go run main.go```

You can also run

```go build main.go``` followed by ```./main``` in order to start the server.

Navigate to ```http://localhost:8080/v1/streams/{id}``` with a streams id to use the api.

## Example

```http://localhost:8080/v1/streams/5938b99cb6906eb1fbaf1f1e```

Additional stream ids can be found in the ```streams.json``` file in the root directory.