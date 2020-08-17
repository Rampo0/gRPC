package main

import (
	"context"
	"fmt"
	"grpc/blog/blogpb"
	"io"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	fmt.Println("Hello I'm a client")

	tls := false
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

	c := blogpb.NewBlogServiceClient(cc)
	// blogID := createBlog(c)
	// readBlog(c, blogID)
	// updateBlog(c, blogID)
	// deleteBlog(c, "5f3a6d7aba2afb0e8827f5d6")
	listBlog(c)
}

func listBlog(c blogpb.BlogServiceClient) {

	stream, err := c.ListBlog(context.Background(), &blogpb.ListBlogRequest{})
	if err != nil {
		log.Fatalf("error while calling ListBlog RPC: %v", err)
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Something happened: %v", err)
		}
		fmt.Println(res.GetBlog())
	}

}

func deleteBlog(c blogpb.BlogServiceClient, blogID string) {
	// delete Blog
	deleteRes, deleteErr := c.DeleteBlog(context.Background(), &blogpb.DeleteBlogRequest{BlogId: blogID})

	if deleteErr != nil {
		fmt.Printf("Error happened while deleting: %v \n", deleteErr)
	}
	fmt.Printf("Blog was deleted: %v \n", deleteRes)

}

func updateBlog(c blogpb.BlogServiceClient, blogID string) {
	newBlog := &blogpb.Blog{
		Id:       blogID,
		AuthorId: "Changed Author",
		Title:    "My First Blog (edited)",
		Content:  "Content of the first blog, with some awesome additions!",
	}
	updateRes, updateErr := c.UpdateBlog(context.Background(), &blogpb.UpdateBlogRequest{Blog: newBlog})
	if updateErr != nil {
		fmt.Printf("Error happened while updating: %v \n", updateErr)
	}
	fmt.Printf("Blog was updated: %v\n", updateRes)
}

func readBlog(c blogpb.BlogServiceClient, blogID string) {
	fmt.Println("Reading the blog")

	// _, err2 := c.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{BlogId: "5f3a6cb10eacbf1535ac6d12"})
	// if err2 != nil {
	// 	fmt.Printf("Error happened while reading: %v \n", err2)
	// }

	readBlogReq := &blogpb.ReadBlogRequest{BlogId: blogID}
	readBlogRes, readBlogErr := c.ReadBlog(context.Background(), readBlogReq)
	if readBlogErr != nil {
		fmt.Printf("Error happened while reading: %v \n", readBlogErr)
	}

	fmt.Printf("Blog was read: %v \n", readBlogRes)
}

func createBlog(c blogpb.BlogServiceClient) string {

	req := &blogpb.CreateBlogRequest{
		Blog: &blogpb.Blog{
			AuthorId: "authorid3",
			Content:  "content3",
			Title:    "title3",
		},
	}

	res, err := c.CreateBlog(context.Background(), req)
	if err != nil {
		log.Fatalf("Unexpected error : %v", err)
	}

	fmt.Printf("Blog has been created : %v", res)
	blogID := res.GetBlog().GetId()
	return blogID
}
