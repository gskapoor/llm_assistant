# GO Backend

This documentation is currently under construction

## Installing Dependencies

```sh
go get
```

## Running the service

1. Put the requisite values in the ".env"
2. Run the Python service
3. Run the following command in this directory:

```sh
go run main.go
```

## Examples (with curl)

```sh
curl -X POST -H "Content-Type: multipart/form-data" -F "audio=@./test.mp3" http://localhost:8080/voice
```

For httpie

```sh
http --multipart localhost:8080/voice audio@test.mp3
```


