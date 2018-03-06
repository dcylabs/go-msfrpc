Light Golang lib to interact with MSFRPC

## Start msfrpc

`msfrpcd -U username -P password -f`

## Check connectivity

`go run main.go 127.0.0.1 55553 username password`

## Expected output

```Error  of 'auth.login':      <nil>
Result of 'console.create':  ??id?0?prompt?msf > ?busy?
Error  of 'console.create':  <nil>```
