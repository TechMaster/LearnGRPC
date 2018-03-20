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
	return &PingMessage{Greeting: in.Greeting}, nil
}

func (s *Server) Add(ctx context.Context, in *TwoNumbers) (*ResultNumber, error) {
	return &ResultNumber{Result: in.A + in.B}, nil
}

//Hàm trả về danh sách Student
func (s *Server) GetStudents(empty *Empty, stream Ping_GetStudentsServer) error {
	students := []Student{
		Student{
			Id:   1,
			Name: "Cuong",
		},
		Student{
			Id:   2,
			Name: "Long",
		},
	}
	for _, student := range students {
		if err := stream.Send(&student); err != nil {
			return err
		}
	}
	return nil //no error

}
