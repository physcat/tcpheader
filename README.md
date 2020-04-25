# Example of how to handle TCP headers in Go

To read a message with a two byte unsigned int lenght indicator:
```go

    import "github.com/physcat/tcpheader"

    length, err := tcpheader.ReadLen(r, tcpheader.TwoByteUnsigned)

    buf := make([]byte, lenght)
    err = tcpheader.ReadMessage(r, buf)
```

To write:
```go
    err := tcpheader.WriteMessage(r, []byte("Message to send"), header)
```

There is a sample program attached:
```bash
$ go get github.com/physcat/tcpheader/cmd/broker
$ broker server
$ broker client
```
