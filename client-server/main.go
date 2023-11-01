package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gofiber/fiber/v2"
	pb "github.com/nauish/go-grpc-server/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	addr     = flag.String("addr", "localhost:50051", "the address to connect to")
	grpcName = flag.String("name", "gRPC", "Name to greet")
)

type ResponseData struct {
	Data string `json:"data"`
}

func sendHTTPRequest(url string, wg *sync.WaitGroup) ResponseData {
	defer wg.Done()

	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return ResponseData{Data: "Error: " + err.Error()}
	}
	defer resp.Body.Close()

	var responseData ResponseData

	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&responseData); err != nil {
		fmt.Printf("Failed to decode JSON response: %v\n", err)
		return ResponseData{Data: "Failed to decode JSON response"}
	}

	return responseData
}

func main() {
	app := fiber.New()

	app.Get("/grpc", func(c *fiber.Ctx) error {
		flag.Parse()
		// Set up a connection to the server.
		numRequests := 10

		var wg sync.WaitGroup
		wg.Add(numRequests)

		responses := make([]*pb.HelloReply, numRequests)

		sendRequest := func(i int) {
			defer wg.Done()
			conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				log.Printf("did not connect: %v", err)
				return
			}
			defer conn.Close()
			client := pb.NewGreeterClient(conn)

			// remove cancel
			ctx := context.Background()
			r, err := client.SayHello(ctx, &pb.HelloRequest{Name: *grpcName})
			if err != nil {
				log.Printf("could not greet: %v", err)
				return
			}
			responses[i] = r
		}

		for i := 0; i < numRequests; i++ {
			go sendRequest(i)
		}

		wg.Wait()

		return c.JSON(responses)
	})

	app.Get("/go-http", func(c *fiber.Ctx) error {
		resp, err := http.Get("http://localhost:5000")
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("request failed with status: %s", resp.Status)
		}

		var responseData ResponseData
		err = json.NewDecoder(resp.Body).Decode(&responseData)
		if err != nil {
			return err
		}

		return c.JSON(responseData)
	})

	app.Get("/goroutine-http", func(c *fiber.Ctx) error {

		url := "http://localhost:5000"
		numRequests := 10

		var wg sync.WaitGroup
		var responses []ResponseData

		for i := 0; i < numRequests; i++ {
			wg.Add(1)
			go func() {
				responseData := sendHTTPRequest(url, &wg)
				responses = append(responses, responseData)
			}()
		}

		wg.Wait()

		return c.JSON(responses)
	})

	app.Get("/node-http", func(c *fiber.Ctx) error {
		resp, err := http.Get("http://localhost:8080")
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("request failed with status: %s", resp.Status)
		}

		var responseData ResponseData
		err = json.NewDecoder(resp.Body).Decode(&responseData)
		if err != nil {
			return err
		}

		return c.JSON(responseData)
	})

	app.Listen(":3000")
}
