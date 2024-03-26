## Prerequisites

- If you want build mock server, make sure you have Go installed on your system. You can download it from [golang.org](https://golang.org/dl/).

## Generating Builds

#### For Linux:

```bash
GOOS=linux GOARCH=amd64 go build -o mock-server
chmod a+x mock-server
```

#### For macOS

```bash
GOOS=darwin GOARCH=amd64 go build -o mock-server
```

#### For Windows

```bash
GOOS=windows GOARCH=amd64 go build -o mock-server.exe
```

## Config Files
At the root of the project, there should be two configuration files:

*server.config.json*

```json
{
  "port": 8080,
  "path": "/"
}
```
- **port**: Port that the mock server will use
- **path**: Path that will be the root of the requests

*hosts.config.json*
```json
{
  "host": "https://api.quotable.io",
  "enabled": true,
  "use_mock": true
}
```
- **host**: URL of the destination host's API
- **enabled**: Enables and utilizes the host for requests
- **use_mock**: Enables the use of mocks using files generated via mock server

**Note**: The mock server supports only one enabled host. If more than one host is enabled, an error is thrown

## Output Files

When a request is made:
- The request is saved in a file at the root. *Ex: request_\<url>.json*
- The response is saved in a file at the root. *Ex: response_\<url>.json*

## Mock responses

When modifying the response file (created by mock server), the mock server will return the modified response, including status codes for error (4xx and 5xx).