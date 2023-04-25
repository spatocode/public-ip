# public-ip

Get your public IP address

## Install

```sh
go get github.com/spatocode/public-ip
```

## Usage
```go
ip, _ := publicip.V4()
fmt.Println(ip) // 101.122.42.103
```