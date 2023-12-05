
# REST Sitemap Parser Service

This service provides functionality to fetch and parse URLs from sitemap XML files based on a given domain or sitemap URL.

## Requirements

- Go 1.18 or later installed.

## Building and Running

### Build the Docker Image

Run the following command to build the Docker image:

```bash
docker build -t sitemap-parser .
```

### Run the Docker Container

After building the Docker image, start the container using:

```bash
docker run -p 8080:8080 sitemap-parser
```

The service will start on port 8080 within the container.

## Endpoints

### 1. `/sitemap`

- **Method**: POST
- **Payload**: `{"sitemap":"<URL to sitemap>"}`

This endpoint fetches and parses the sitemap provided in the payload.

### 2. `/domain`

- **Method**: POST
- **Payload**: `{"domain":"<Domain URL>"}`

This endpoint fetches the sitemap for the given domain and then parses it.

### 3. `/ping`

- **Method**: GET

A simple endpoint to check if the service is running. Returns "Pong!" as a response.

### Root Endpoint `/`

- **Method**: GET

Provides basic information on how to request URLs by POSTing the link to `/sitemap`.

## Example Usage

### Fetch and Parse Sitemap

```bash
curl -X POST -H "Content-Type: application/json" -d '{"sitemap":"https://stackovercode.com/sitemap.xml"}' http://localhost:8080/sitemap
```

### Fetch Sitemap for a Domain and Parse

```bash
curl -X POST -H "Content-Type: application/json" -d '{"domain":"stackovercode.com"}' http://localhost:8080/domain
```

## Notes

- Replace `http://localhost:8080` with the appropriate host and port where the service is running.

For more information, refer to the [Go source code](./main.go).


## Contributing

Contributions are welcome! If you'd like to enhance this API or fix issues, please fork the repository and create a pull request. Any contributions you make are greatly appreciated.

## License

This project is licensed under the MIT License.