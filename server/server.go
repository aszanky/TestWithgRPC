package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/testwithgrpc/justmessagepb"
	"google.golang.org/grpc"
)

//Because I don't use DB, just save message on slice
var response []*justmessagepb.JustMessage

type server struct{}

//JustMessage only Request/Response then save message on response slice
func (s *server) JustMessage(ctx context.Context, req *justmessagepb.MessageRequest) (res *justmessagepb.MessageResponse, err error) {
	fmt.Printf("This function was invoked with %v\n", req)
	yourMessage := req.GetYourmessage()

	//append new message to slice
	response = append(response, yourMessage)

	res = &justmessagepb.MessageResponse{
		Result: yourMessage,
	}

	return
}

//GetAllMessages for getting all messages that was saved on responce slice
func (s *server) GetAllMessages(gmr *justmessagepb.GetAllMessageRequest, stream justmessagepb.JustMessageService_GetAllMessagesServer) error {
	fmt.Println("Getting all messages")

	for _, message := range response {
		res := &justmessagepb.GetAllMessagesResponse{
			Result: &justmessagepb.JustMessage{
				Id:          message.Id,
				Justmessage: message.Justmessage,
			},
		}
		stream.Send(res)
	}

	return nil
}

func (s server) GetCommunicationMessages(stream justmessagepb.JustMessageService_GetCommunicationMessagesServer) error {
	fmt.Println("Let's start communication")
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		yourMessage := req.GetYourmessage()

		//append new message to slice
		response = append(response, yourMessage)

		res := &justmessagepb.LongLiveMessageConnectionResponse{
			Result: yourMessage,
		}

		sendingError := stream.Send(res)
		if sendingError != nil {
			sendingError = fmt.Errorf("Error sending streaming message: %v", sendingError)
			return sendingError
		}

	}
}

func main() {
	fmt.Println("Server Running")

	//Create Listener
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	justmessagepb.RegisterJustMessageServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)

	}
}
