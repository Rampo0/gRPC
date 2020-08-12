package main

import (
	"context"
	"fmt"
	"grpc/calculator/calculatorpb"
	"io"
	"log"
	"time"

	"google.golang.org/grpc/codes"

	"google.golang.org/grpc/status"

	"google.golang.org/grpc"
)

func main() {

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatal("Failed to Dial : %v", err)
	}

	defer cc.Close()

	c := calculatorpb.NewCalculatorServiceClient(cc)

	// doUnary(c)

	// doServerStreaming(c)

	// doClientStreaming(c)

	// doBiDiStreaming(c)

	doErrorUnary(c)
}

func doErrorUnary(c calculatorpb.CalculatorServiceClient) {

	// correct call
	doErrorCall(c, 10)

	// incorrect call
	doErrorCall(c, -2)
}

func doErrorCall(c calculatorpb.CalculatorServiceClient, n int32) {
	res, err := c.SquareRoot(context.Background(), &calculatorpb.SquareRootRequest{Number: n})

	if err != nil {
		respErr, ok := status.FromError(err)
		if ok {
			// actual error from grpc (user error)
			fmt.Printf("Error message from server : %v\n", respErr.Message())
			fmt.Println(respErr.Code())
			if respErr.Code() == codes.InvalidArgument {
				fmt.Println("We probably sent a negative number !")
				return
			}

		} else {
			log.Fatalf("Big error calling square root : %v", err)
			return
		}
	}

	fmt.Printf("Result of square root of %v : %v\n", n, res.GetNumberRoot())
}

func doBiDiStreaming(c calculatorpb.CalculatorServiceClient) {

	stream, err := c.MaximumNumber(context.Background())

	if err != nil {
		log.Fatalf("Error while invoking function : %v\n", err)
	}

	waitc := make(chan struct{})

	go func() {
		// function to send req
		numbers := []int64{1, 5, 3, 6, 2, 20}

		for _, number := range numbers {
			fmt.Printf("Sending number : %v\n", number)
			stream.Send(&calculatorpb.MaximumNumberRequest{
				Number: number,
			})
			time.Sleep(1000 * time.Millisecond)
		}

		stream.CloseSend()

	}()

	go func() {
		// function to recv response

		for {
			res, err := stream.Recv()

			if err == io.EOF {
				break
			}

			if err != nil {
				log.Fatalf("Error while receiving response : %v\n", err)
				break
			}

			fmt.Printf("%v \n", res.GetResult())
		}
		close(waitc)
	}()

	<-waitc

}

func doClientStreaming(c calculatorpb.CalculatorServiceClient) {

	stream, err := c.AverageNumber(context.Background())
	if err != nil {
		log.Fatalf("Error while invoke function average number : %v \n", err)
	}

	requests := []*calculatorpb.AverageNumberRequest{
		&calculatorpb.AverageNumberRequest{
			Number: 2,
		},
		&calculatorpb.AverageNumberRequest{
			Number: 3,
		},
		&calculatorpb.AverageNumberRequest{
			Number: 4,
		},
		&calculatorpb.AverageNumberRequest{
			Number: 6,
		},
		&calculatorpb.AverageNumberRequest{
			Number: 1,
		},
	}

	for _, req := range requests {
		fmt.Printf("Sending request : %v \n", req)
		err := stream.Send(req)
		if err != nil {
			log.Fatalf("Error while sending a request : %v \n", err)
		}

	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while getting response : %v", err)
	}

	fmt.Printf("The average number is : %v \n", res.GetResult())

}

func doServerStreaming(c calculatorpb.CalculatorServiceClient) {
	req := &calculatorpb.PrimeNumberDecompositionRequest{
		Number: 120,
	}

	resStream, err := c.PrimeNumberDecomposition(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling PrimeNumberDecompostion RPC: %v", err)
	}

	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("Error while reading the stream : %v", err)
		}

		log.Printf("%v", msg.GetResult())
	}

}

func doUnary(c calculatorpb.CalculatorServiceClient) {

	req := &calculatorpb.CalculatorRequest{
		Sum: &calculatorpb.SumNumber{
			FirstNumber:  10,
			SecondNumber: 4,
		},
	}

	res, _ := c.Sum(context.Background(), req)

	fmt.Println(res.Result)

}
