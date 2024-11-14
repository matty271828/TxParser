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
-  Testing via automated tests
-  Logging to take place at top level of stacktrace. Each API request is also to be logged. 

## Usage

