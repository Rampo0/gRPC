package main

import (
	"context"
	"fmt"
	"grpc/calculator/calculatorpb"
	"io"
	"log"
	"math"
	"net"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	"google.golang.org/grpc"
)

type server struct{}

func (*server) Sum(ctx context.Context, req *calculatorpb.CalculatorRequest) (*calculatorpb.CalculatorResponse, error) {
	firstNum := req.GetSum().GetFirstNumber()
	secondNum := req.GetSum().GetSecondNumber()

	result := firstNum + secondNum
	response := &calculatorpb.CalculatorResponse{
		Result: result,
	}

	return response, nil
}

func (*server) PrimeNumberDecomposition(req *calculatorpb.PrimeNumberDecompositionRequest, stream calculatorpb.CalculatorService_PrimeNumberDecompositionServer) error {

	var number int64 = req.GetNumber()
	var k int64 = 2

	for {

		if number <= 1 {
			break
		}

		if number%k == 0 {
			response := &calculatorpb.PrimeNumberDecompositionResponse{
				Result: k,
			}

			stream.Send(response)

			//time.Sleep(1000 * time.Millisecond)

			number = number / k
		} else {
			k++
		}
	}

	return nil
}

func (*server) AverageNumber(stream calculatorpb.CalculatorService_AverageNumberServer) error {

	count := int64(0)
	result := int64(0)

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			average := float32(result) / float32(count)
			return stream.SendAndClose(&calculatorpb.AverageNumberResponse{
				Result: average,
			})
		}

		if err != nil {
			log.Fatalf("Error while receiving request : %v \n", err)
		}

		result += req.GetNumber()
		count++
	}

}

func (*server) MaximumNumber(stream calculatorpb.CalculatorService_MaximumNumberServer) error {

	max := int64(-1)
	var sendErr error
	for {
		req, err := stream.Recv()

		if err == io.EOF {
			return nil
		}

		if err != nil {
			log.Fatalf("Error while receiving request : %v", err)
			return err
		}

		if req.GetNumber() > max {
			max = req.GetNumber()
			sendErr = stream.Send(&calculatorpb.MaximumNumberResponse{
				Result: max,
			})
		}

		if sendErr != nil {
			log.Fatalf("Error while sending a response : %v", sendErr)
			return nil
		}

	}

}

func (*server) SquareRoot(ctx context.Context, req *calculatorpb.SquareRootRequest) (*calculatorpb.SquareRootResponse, error) {
	number := req.GetNumber()

	if number < 0 {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Received negative number : %v\n", number),
		)
	}

	return &calculatorpb.SquareRootResponse{
		NumberRoot: math.Sqrt(float64(number)),
	}, nil

}

func main() {

	fmt.Println("Application running...")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatal("failed to listen : %v", err)
	}

	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	reflection.Register(s)

	if err := s.Serve(lis); err != nil {
		log.Fatal("failed to Serve : %v", err)
	}

}
