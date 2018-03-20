# Các bước thực hành gRPC
## Chuẩn bị
1. Cài đặt golang. ```brew install go```
2. Cài đặt protobuf tool ```brew install protobuf```
3. Kiểm tra biến ```echo $GOPATH```
  - Nếu biến GOPATH chưa tồn tại hãy tạo một thư mục sẽ lưu các dự án Golang của bạn và các thư viện Golang tải về từ Github, Gitlab....
  - Lệnh tạo biến đường dẫn $GOPATH nếu bạn đang chạy Fish Shell thì như sau:
    ```
    set -gx GOPATH /PathToGoProjects/
    set -gx PATH /PathToGoProjects/bin $PATH
    ```
  - Tham khảo thêm [Lập trình Golang trên MacOSX](https://techmaster.vn/posts/34567/lap-trinh-golang-iris-framework-tren-macosx)
4. Kiểm tra môi trường lập trình Golang
  ```
  $ go version
  go version go1.10 darwin/amd64
  $ protoc --version
  libprotoc 3.5.1
  ```
5. Tạo thư mục dự án
```
$ mkdir $GOPATH\src\github.com\TechMaster\LearnGRPC
$ cd $GOPATH\src\github.com\TechMaster\LearnGRPC
$ mkdir api
```
## Tạo protobuf file
Trong thư mục $GOPATH\src\github.com\TechMaster\LearnGRPC\api tạo file
```proto
syntax = "proto3";
package api;
message PingMessage {
  string greeting = 1;
}
service Ping {
  rpc SayHello(PingMessage) returns (PingMessage) {}
}
```
Sau đó chạy lệnh để biên dịch file định nghĩa protobuf sang mã nguồn Golang
```
protoc -I api/ \
          -I {$GOPATH}/src \
          --go_out=plugins=grpc:api \
          api/api.proto
```
Kết quả đầu ra sẽ là file api\api.pb.go
## Tạo gRPC hander
Trong thư mục api, tạo handler.go
```go
package api

import (
	"log"

	"golang.org/x/net/context"
)

// Server represents the gRPC server
type Server struct {
}

// SayHello generates response to a Ping request
func (s *Server) SayHello(ctx context.Context, in *PingMessage) (*PingMessage, error) {
	log.Printf("Receive message %s", in.Greeting)
	return &PingMessage{Greeting: "bar"}, nil
}
```
Giải thích handler.go
- Server struct sẽ đại diện cho gRPC server. Nó chứa các hàm ứng với các hàm đã định nghĩa trong file proto, bản chất là nó phải implement interface này
  ```go
  // Phần này ở file api.pb.go được sinh ra sau khi biên dịch api.proto
  type PingServer interface {
	  SayHello(context.Context, *PingMessage) (*PingMessage, error)
  }
  
  // Đây là hàm trong api/handler.go implement interface
  func (s *Server) SayHello(ctx context.Context, in *PingMessage) (*PingMessage, error) {
	  log.Printf("Receive message %s", in.Greeting)
	  return &PingMessage{Greeting: "bar"}, nil
  }
  ```
## Tạo server main.go
Sau đó biên dịch. Chú ý nếu để file đúng thư mục thì go build biên dịch mà không cần pull code từ github.com. Nếu chưa có sẵn code thì go build sẽ pull code từ github.
```
go build -i -v -o bin/server github.com/TechMaster/LearnGRPC/server
```
## Tạo client main.go
```go
package main

import (
	"log"

	"github.com/TechMaster/LearnGRPC/api"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func main() {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":7777", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()
	c := api.NewPingClient(conn)
	response, err := c.SayHello(context.Background(), &api.PingMessage{Greeting: "foo"})
	if err != nil {
		log.Fatalf("Error when calling SayHello: %s", err)
	}
	log.Printf("Response from server: %s", response.Greeting)
}
```
Biên dịch
```
go build -i -v -o bin/client github.com/TechMaster/LearnGRPC/client
```
# Tham khảo
1. [How we use gRPC to build a client/server system in Go](https://medium.com/pantomath/how-we-use-grpc-to-build-a-client-server-system-in-go-dd20045fa1c2)