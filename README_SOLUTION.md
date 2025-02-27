# Payment Gateway API

This is a simple API that includes the following endpoints:

- **Deposit**: Handles `Content-Type` of JSON and XML.
- **Withdraw**: Handles `Content-Type` of JSON and XML.
- **Callback**: Handles `Content-Type` of JSON and XML.
- **Update Gateway Priority**: Handles `Content-Type` of JSON and allows specifying preferences for a specific gateway provider.

## Architecture

- **RabbitMQ** was chosen instead of Kafka, as the complexity of Kafka is not required for a simple callback message that contains a transaction ID and a status to update.

- **Redis** is used to store failed transaction IDs and their statuses in case the method [`internal/app/services/transaction_publisher_service.go:22`](internal/app/services/transaction_publisher_service.go#L22) fails to publish the transaction using the circuit breaker.

- The **Retry Operation Pattern** is implemented in [`internal/app/services/fault_tolerance.go:34`](internal/app/services/fault_tolerance.go#L34) to attempt creating a transaction multiple times.

- **PostgreSQL** is used for storing data related to transactions, gateways, countries, and users. Since this system does not deal with user balance changes, explicit locking using `FOR UPDATE` was not implemented.

## Code Structure

- The code structure follows a lightweight implementation of **Domain-Driven Design (DDD)**.
- Tests and mocks are written using the **Mockery** library.

## How to Run

- `make dc`: Runs Docker Compose with the app container exposed on port 8080.
- `make test`: Executes the test suite.
- `make run`: Runs the app locally on port 8080 without Docker. Before running this, make sure to start the PostgreSQL, Redis, and RabbitMQ containers using `make dc`.
- `make lint`: Runs the linter.

### Short Project Demo

## Redis Failed Transaction Creation Script
demo url - https://youtu.be/WMAR5n2BCEo

- A script for creating Redis failed transactions is located at [`testing/redis_failed_transaction_testing.go:10`](testing/redis_failed_transaction_testing.go#L10). Use this script for testing purposes if needed.
