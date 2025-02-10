# Post-to-Get Middleware

This project is a simple HTTP middleware that converts POST requests to GET requests and proxies them to a legacy server.

## Requirements

- Go 1.23.5 or later

## Installation

1. Clone the repository:
    ```sh
    git clone https://github.com/yourusername/post-to-get-middleware.git
    cd post-to-get-middleware
    ```

2. Build the project:
    ```sh
    go build -o post-to-get-middleware main.go
    ```

## Usage

Run the middleware with the following command:
```sh
./post-to-get-middleware --legacy-server=http://localhost:8081 --port=8080
```

### Options

- `--legacy-server`: The legacy server to proxy requests to (required).
- `--port`: The port to run the proxy server on (default: 8080).
- `--cert-file`: The path to the certificate file (optional, for HTTPS).
- `--key-file`: The path to the key file (optional, for HTTPS).

## Example

To run the middleware and proxy requests to a legacy server running on `http://localhost:8081`:
```sh
./post-to-get-middleware --legacy-server=http://localhost:8081
```

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
