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
[
  {
    "url": "https://api.quotable.io",
    "enabled": true,
    "enable_mock": true,
    "enable_graphql": true,
    "enable_rest": true,
    "generate_path": "/generate-files",
    "graphql_path": "/graphql",
    "output_folder": "quotable"
  }
]
```
- **url**: URL of the destination host's API
- **enabled**: Enables and utilizes the host for requests
- **enable_mock**: Enables the use of mocks using files generated via mock server
- **enable_graphql**: Enables GraphQL request
- **enable_rest**: Enables REST request
- **graphql_path**: Path to use in GraphQL request
- **generate_path**: Path to generate empty request and response files 
- **output_folder**: Folder to store request and response files inside .output folder

**Note**: The mock server supports only one enabled host. If more than one host is enabled, an error is thrown

## Output Files

### REST
- The request is saved in a file at the output folder. *Ex: request_\<url>.json*
- The response is saved in a file at the output folder. *Ex: response_\<url>.json*

### GraphQL
- The query is saved in a file at the output folder. *Ex: request_\<url>_query\<MD5>.json*
- The mutation is saved in a file at the output folder. *Ex: request_\<url>_mutation\<MD5>.json*
- The query response is saved in a file at the output folder. *Ex: response_\<url>_query\<MD5>.json*
- The mutation response is saved in a file at the output folder. *Ex: response_\<url>_mutation\<MD5>.json*

## Mock responses

When modifying the response file (created by mock server), the mock server will return the modified response, including status codes for error (4xx and 5xx) and delay by reponse time

## TODO

- [x] Generate blank file
- [x] GraphQL Support ![Beta](https://img.shields.io/badge/-Beta-orange)
- [x] Docker Support ![Beta](https://img.shields.io/badge/-Beta-orange)
- [x] Request and response files in output folder
- [ ] How it works in README
