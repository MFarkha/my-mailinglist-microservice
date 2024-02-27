package main

import (
	"context"
	"log"
	"time"

	pb "github.com/MFarkha/my-mailinglist-microservice/proto"
	"go.wit.com/dev/alexflint/arg"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func logResponse(res *pb.EmailResponse, err error) {
	if err != nil {
		log.Fatalf(" error: %v\n", err)
	}
	if res.EmailEntry == nil {
		log.Println(" email not found")
	} else {
		log.Printf(" response: %v\n", res.EmailEntry)
	}
}

func createEmail(client pb.MailingListServiceClient, emailAddr string) *pb.EmailEntry {
	log.Println("createEmail")
	// gRPC server requires context with request, timeout with 1 second
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	res, err := client.CreateEmail(ctx, &pb.CreateEmailRequest{EmailAddr: emailAddr})
	logResponse(res, err)
	return res.EmailEntry
}

func getEmail(client pb.MailingListServiceClient, emailAddr string) *pb.EmailEntry {
	log.Println("getEmail")
	// gRPC server requires context with request, timeout with 1 second
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	res, err := client.GetEmail(ctx, &pb.GetEmailRequest{EmailAddr: emailAddr})
	logResponse(res, err)
	return res.EmailEntry
}

func getEmailBatch(client pb.MailingListServiceClient, page int32, count int32) []*pb.EmailEntry {
	log.Println("getEmailBatch")
	// gRPC server requires context with request, timeout with 1 second
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	res, err := client.GetEmailBatch(ctx, &pb.GetEmailBatchRequest{Page: page, Count: count})
	if err != nil {
		log.Fatalf(" error: %v\n", err)
	}
	log.Println(" response:")
	for i, e := range res.EmailEntries {
		log.Printf("item: [%v] %v\n", i, e)
	}
	return res.EmailEntries
}

func updateEmail(client pb.MailingListServiceClient, emailEntry *pb.EmailEntry) {
	log.Println("updateEmail")
	// gRPC server requires context with request, timeout with 1 second
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	res, err := client.UpdateEmail(ctx, &pb.UpdateEmailRequest{EmailEntry: emailEntry})
	logResponse(res, err)

}

func deleteEmail(client pb.MailingListServiceClient, emailAddr string) *pb.EmailEntry {
	log.Println("deleteEmail")
	// gRPC server requires context with request, timeout with 1 second
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	res, err := client.DeleteEmail(ctx, &pb.DeleteEmailRequest{EmailAddr: emailAddr})
	logResponse(res, err)
	return res.EmailEntry
}

var args struct {
	GrpcAddr string `arg:"env:MAILINGLIST_GRPC_ADDR"`
}

func main() {
	arg.MustParse(&args)

	if args.GrpcAddr == "" {
		args.GrpcAddr = ":3001"
	}
	// this microservice is running on a backend on preauthorized service,
	// therefore no need for secure connection
	opts := grpc.WithTransportCredentials(insecure.NewCredentials())
	conn, err := grpc.Dial(args.GrpcAddr, opts)
	if err != nil {
		log.Fatalf("unable to connect to a server: %v\n", args.GrpcAddr)
	}
	defer conn.Close()
	client := pb.NewMailingListServiceClient(conn)
	// newEmailEntry := createEmail(client, "mark.matt@example.com")
	// newEmailEntry.ConfirmedAt = 10000
	// updateEmail(client, newEmailEntry)
	// deleteEmail(client, "jason.smith@example.com")
	getEmailBatch(client, 1, 2)
	getEmailBatch(client, 2, 2)
}
