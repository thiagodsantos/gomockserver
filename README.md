## Prerequisites

- Make sure you have Go installed on your system. You can download it from [golang.org](https://golang.org/dl/).

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