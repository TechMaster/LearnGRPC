package main

import (
	"log"
	"io"

	"github.com/TechMaster/LearnGRPC/api"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

//Bổ xung cho phần login của user
// Authentication holds the login/password
type Authentication struct {
	Login    string
	Password string
	Role     string
}

/* Authentication struct adopt interface PerRPCCredentials bằng
cách implement 2 phương thức
1. GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error)
2. RequireTransportSecurity() bool
*/
/* GetRequestMetadata gets the current request metadata
Bổ xung thêm thông tin vào request metdata
*/
func (a *Authentication) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
	return map[string]string{
		"login":    a.Login,
		"password": a.Password,
		"role":     a.Role,
	}, nil
}

/* RequireTransportSecurity indicates whether the credentials requires transport security
trả về true nếu yêu cầu mã hoá
*/
func (a *Authentication) RequireTransportSecurity() bool {
	return true
}

//--------
func main() {
	var conn *grpc.ClientConn
	// Create the client TLS credentials
	creds, err := credentials.NewClientTLSFromFile("cert/server.crt", "")
	if err != nil {
		log.Fatalf("could not load tls cert: %s", err)
	}

	// Setup the login/pass
	auth := Authentication{
		Login:    "john",
		Password: "doe",
		Role:     "operator",
	}

	// Initiate a connection with the server
	conn, err = grpc.Dial("localhost:7777",
		grpc.WithTransportCredentials(creds),
		grpc.WithPerRPCCredentials(&auth))

	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	c := api.NewPingClient(conn)

	response, err := c.SayHello(context.Background(), &api.PingMessage{Greeting: "foo"})

	if err != nil {
		log.Fatalf("error when calling SayHello: %s", err)
	}
	log.Printf("Response from server: %s", response.Greeting)

	response, err = c.SayHello(context.Background(), &api.PingMessage{Greeting: "fart"})
	if err != nil {
		log.Fatalf("error when calling SayHello: %s", err)
	}
	log.Printf("Response from server: %s", response.Greeting)

	stream, err := c.GetStudents(context.Background(), &api.Empty{})
	if err != nil {
		log.Fatalf("Error on get customers: %v", err)
	}
	for {
		// Receiving the stream of data
		student, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("%v.GetStudents(_) = _, %v", c, err)
		}
		log.Printf("Students: %v", student)
	}

}
