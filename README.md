# TxParser

## Task
- Implement a parser that can act as part of a notification service
which informs a user about incoming and outgoing transactions for 
ethereum addresses that a user chooses to subscribe to.

Required operations:
- Get current block number
- Subscribe to a new address
- Get transactions for a given address

## Design 

    ### cmd
    The main entry point for the application.

    ### internal
    Contains the core logic of the application.

        #### parser
        Contains the public interface for the parser and parsing logic.

        #### server
        Http server implementation to serve api requests from the user to the parser. 

        #### ethclient
        Interact with the ethereum blockchain via Ethereum JSONRPC. 

        #### storer
        In memory store for data. Could be extended in the future to connect to a database.

## Development Practices
    -  Testing via end to end integration tests. In reality I would have added unit testing and
    utilised a mocking framework to mock out the dependencies.
    -  Logging to take place at top level of stacktrace. Each API request is also to be logged. 

## Usage

### Prerequisites
    - Go 1.22 or higher
    - Git

### Running the application

1. Build and run the application
```bash
go build -o txparser ./cmd/main.go

./txparser
```

2. Hit some endpoints

- Subscribe to an address
```bash
curl "http://localhost:8080/subscribe?address=0x742d35Cc6634C0532925a3b844Bc454e4438f44e"
``` 

- Get the current block number
```bash
curl "http://localhost:8080/getcurrentblock"
```

- Get transactions for an address
```bash
curl "http://localhost:8080/gettransactions?address=0x742d35Cc6634C0532925a3b844Bc454e4438f44e"
``` 
