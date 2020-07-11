package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math/rand"
	"time"

	"github.com/testwithgrpc/justmessagepb"

	"google.golang.org/grpc"
)

//Messages struct
type Messages struct {
	messageInterface justmessagepb.JustMessageServiceClient
}

//NewClient for initialization
func NewClient() *Messages {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Can't Connect to the Server: %v", err)
	}

	//instantiate object
	c := justmessagepb.NewJustMessageServiceClient(conn)

	return &Messages{
		messageInterface: c,
	}
}

func (m *Messages) messageRequest(theMessage string) (res *justmessagepb.JustMessage, err error) {
	id := rand.Intn(100)
	req := &justmessagepb.MessageRequest{
		Yourmessage: &justmessagepb.JustMessage{
			Id:          int32(id),
			Justmessage: theMessage,
		},
	}
	response, err := m.messageInterface.JustMessage(context.Background(), req)
	if err != nil {
		err = fmt.Errorf("[ERROR] calling JustMessage RPC: %v", err)
		return
	}

	//This is Response from server JustMessage
	return response.Result, nil
}

func (m *Messages) getAllMessages() error {
	streamRes, err := m.messageInterface.GetAllMessages(context.Background(), &justmessagepb.GetAllMessageRequest{})
	if err != nil {
		return err
	}

	for {
		res, err := streamRes.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		fmt.Println(res.GetResult())
	}

	return nil
}

//LetsTalk for talking immediately
func (m *Messages) LetsTalk(theMessages []string) error {
	streamRes, err := m.messageInterface.GetCommunicationMessages(context.Background())
	if err != nil {
		return err
	}

	var messageReq []*justmessagepb.LongLiveMessageConnectionRequest
	for idx, val := range theMessages {
		messageReq = append(messageReq, &justmessagepb.LongLiveMessageConnectionRequest{
			Yourmessage: &justmessagepb.JustMessage{
				Id:          int32(idx),
				Justmessage: val,
			},
		})
	}

	waitc := make(chan struct{})
	// we send a bunch of messages to the client (go routine)
	go func() {
		// function to send a bunch of messages
		for _, req := range messageReq {
			fmt.Printf("Sending message: %v\n", req)
			streamRes.Send(req)
			time.Sleep(1000 * time.Millisecond)
		}
		streamRes.CloseSend()
	}()
	// we receive a bunch of messages from the client (go routine)
	go func() {
		// function to receive a bunch of messages
		for {
			res, err := streamRes.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("[Error] Receiving: %v", err)
				break
			}
			fmt.Printf("Received: %v\n", res.GetResult())
		}
		close(waitc)
	}()

	// block until everything is done
	<-waitc

	return nil
}

func main() {
	fmt.Println("This is a Client ")
	fmt.Println("Create Connection to the Server")

	//Client Initialization
	client := NewClient()

	fmt.Println("-------------------------")
	fmt.Println("Task 1")
	fmt.Println("-------------------------")

	//Request Response - Task 1
	// response, err := client.messageRequest("Message 1")
	// response1, err := client.messageRequest("Message 2")
	// response2, err := client.messageRequest("Message 3")
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// fmt.Println(response)
	// fmt.Println(response1)
	// fmt.Println(response2)

	// fmt.Println("-------------------------")
	// fmt.Println("Task 2")
	// fmt.Println("-------------------------")

	// //Get All Messages with server streaming - Task 2
	// errGetAll := client.getAllMessages()
	// if errGetAll != nil {
	// 	fmt.Println(errGetAll)
	// }

	fmt.Println("-------------------------")
	fmt.Println("Task 3")
	fmt.Println("-------------------------")

	//Bi-directional streaming gRPC for long live connection
	messageRequest := []string{"Message1", "Message2", "Message3", "Message4", "Message5"}
	errTalking := client.LetsTalk(messageRequest)
	if errTalking != nil {
		fmt.Println(errTalking)
	}
}
