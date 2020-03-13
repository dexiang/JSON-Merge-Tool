# JSON Merge Tool

Provide tools for merging JSON files

## Usage
```shell script
$ go run main.go -input ${file | folder} -output ${file} -csv
```

## Development Note

### Go Modules Init.

```shell script
$ go mod init json-merge-tool
```

### Build
```shell script
# Mac
$ CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o JSON-Merge-Tool
```
