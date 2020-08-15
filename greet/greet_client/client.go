package main

import (
	"context"
	"fmt"
	"grpc/greet/greetpb"
	"io"
	"log"
	"time"

	"google.golang.org/grpc/credentials"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Hello I'm a client")

	tls := true
	opts := grpc.WithInsecure()
	if tls {
		certFile := "ssl/ca.crt" // certificate authority trust certificate
		creds, sslErr := credentials.NewClientTLSFromFile(certFile, "")

		if sslErr != nil {
			log.Fatalf("Failed loading certificate : %v", sslErr)
			return
		}

		opts = grpc.WithTransportCredentials(creds)
	}

	cc, err := grpc.Dial("localhost:50051", opts)
	if err != nil {
		log.Fatalf("could not connect : %v", err)
	}
	defer cc.Close()

	c := greetpb.NewGreetServiceClient(cc)
	doUnary(c)

	// doServerStreaming(c)

	//doClientStreaming(c)

	// doBiDiStreaming(c)

	// doUnaryWithDeadline(c, 5*time.Second)
	// doUnaryWithDeadline(c, 1*time.Second)
}

func doUnaryWithDeadline(c greetpb.GreetServiceClient, timeout time.Duration) {
	req := &greetpb.GreetWithDeadlineRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Ammar",
			LastName:  "Anwari",
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	res, err := c.GreetWithDeadline(ctx, req)
	if err != nil {

		statusErr, ok := status.FromError(err)

		if ok {
			if statusErr.Code() == codes.DeadlineExceeded {
				fmt.Println("Timeout was hit")
			} else {
				fmt.Printf("Unexpected error : %v", statusErr)
			}
		} else {
			log.Fatalf("error while calling Greet RPC : %v", err)
		}

		return
	}

	log.Printf("Response from Greet; %v", res.Result)

}

func doBiDiStreaming(c greetpb.GreetServiceClient) {

	// create a stream by invoking client
	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("Error while creating stream : %v", err)
		return
	}

	requests := []*greetpb.GreetEveryoneRequest{
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Anwar 1",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Anwar 2",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Anwar 3",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Anwar 4",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Anwar 5",
			},
		},
	}

	waitc := make(chan struct{})
	// send a bunch of messages to the client (goroutine)
	go func() {
		// function to send a bunch of messages
		for _, req := range requests {
			fmt.Printf("Sending message : %v \n", req)
			stream.Send(req)
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()

	// receive bunc of messages from the client (go routine)
	go func() {
		// function to receive a bunch of response
		for {
			res, err := stream.Recv()

			if err == io.EOF {
				break
			}

			if err != nil {
				log.Fatalf("Error while receiving : %v \n", err)
				break
			}

			fmt.Printf("Received: %v\n", res.GetResult())
		}
		close(waitc)
	}()

	// block until everything is done
	<-waitc
	fmt.Printf("Done BiDi requesting. \n")
}

func doClientStreaming(c greetpb.GreetServiceClient) {

	requests := []*greetpb.LongGreetRequest{
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Anwar 1",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Anwar 2",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Anwar 3",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Anwar 4",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Anwar 5",
			},
		},
	}

	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("error while calling long greet : %v", err)
	}

	for _, req := range requests {
		fmt.Printf("Sending req : %v \n", req)
		stream.Send(req)
		time.Sleep(100 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while receiving response : %v", err)
	}

	fmt.Printf("Response LongGreet : %v\n", res)

}

func doServerStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to streaming RPC server")

	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Ammar",
			LastName:  "Anwari",
		},
	}

	resStream, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling GreetManyTimes RPC: %v", err)
	}

	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("Error while reading the stream : %v", err)
		}

		log.Printf("Response from greetManyTimes : %v", msg.GetResult())
	}

}

func doUnary(c greetpb.GreetServiceClient) {

	fmt.Println("Starting todo a Unary RPC")

	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Ammar",
			LastName:  "Anwari",
		},
	}

	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Greet RPC : %v", err)
	}

	log.Printf("Response from Greet; %v", res.Result)
}
