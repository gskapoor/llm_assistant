# GO Backend

This documentation is currently under construction

## TESTS

```sh
curl -X POST -H "Content-Type: multipart/form-data" -F "audio=@./test.mp3" http://localhost:8080/voice
```

For httpie

```sh
http --multipart localhost:8080/voice audio@test.mp3
```


