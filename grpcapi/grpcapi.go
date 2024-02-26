package grpcapi

import (
	"context"
	"database/sql"
	"log"
	"net"
	"time"

	"github.com/MFarkha/my-mailinglist-microservice/mdb"
	pb "github.com/MFarkha/my-mailinglist-microservice/proto"
	"google.golang.org/grpc"
)

type MailServer struct {
	pb.UnimplementedMailingListServiceServer
	db *sql.DB
}

func pbEntryToMdbEntry(pbEntry *pb.EmailEntry) mdb.EmailEntry {
	t := time.Unix(pbEntry.ConfirmedAt, 0)
	return mdb.EmailEntry{
		Id:          pbEntry.Id,
		Email:       pbEntry.Email,
		ConfirmedAt: &t,
		OptOut:      pbEntry.OptOut,
	}
}

func mdbEntryTopbEntry(mdbEntry *mdb.EmailEntry) pb.EmailEntry {
	return pb.EmailEntry{
		Id:          mdbEntry.Id,
		Email:       mdbEntry.Email,
		ConfirmedAt: mdbEntry.ConfirmedAt.Unix(),
		OptOut:      mdbEntry.OptOut,
	}
}

func emailResponse(db *sql.DB, email string) (*pb.EmailResponse, error) {
	mdbEntry, err := mdb.GetEmailEntry(db, email)
	if err != nil {
		return nil, err
	}
	if mdbEntry == nil {
		return &pb.EmailResponse{}, nil
	}
	pbEntry := mdbEntryTopbEntry(mdbEntry)
	return &pb.EmailResponse{EmailEntry: &pbEntry}, nil
}

func (s *MailServer) GetEmail(ctx context.Context, req *pb.GetEmailRequest) (*pb.EmailResponse, error) {
	log.Printf("gRPC GetEmail: %v\n", req.EmailAddr)
	return emailResponse(s.db, req.EmailAddr)
}

func (s *MailServer) CreateEmail(ctx context.Context, req *pb.CreateEmailRequest) (*pb.EmailResponse, error) {
	log.Printf("gRPC CreateEmail: %v\n", req.EmailAddr)
	err := mdb.CreateEmailEntry(s.db, req.EmailAddr)
	if err != nil {
		return nil, err
	}
	return emailResponse(s.db, req.EmailAddr)
}

func (s *MailServer) UpdateEmail(ctx context.Context, req *pb.UpdateEmailRequest) (*pb.EmailResponse, error) {
	log.Printf("gRPC UpdateEmail: %v\n", req.EmailEntry)
	mdbEntry := pbEntryToMdbEntry(req.EmailEntry)
	err := mdb.UpdateEmailEntry(s.db, &mdbEntry)
	if err != nil {
		return nil, err
	}
	return emailResponse(s.db, req.EmailEntry.Email)
}

func (s *MailServer) DeleteEmail(ctx context.Context, req *pb.DeleteEmailRequest) (*pb.EmailResponse, error) {
	log.Printf("gRPC DeleteEmail: %v\n", req.EmailAddr)
	err := mdb.DeleteEmailEntry(s.db, req.EmailAddr)
	if err != nil {
		return nil, err
	}
	return emailResponse(s.db, req.EmailAddr)
}

func (s *MailServer) GetEmailBatch(ctx context.Context, req *pb.GetEmailBatchRequest) (*pb.GetEmailBatchResponse, error) {
	log.Printf("gRPC GetEmailBatch: %v %v\n", req.Count, req.Page)
	queryOptions := mdb.GetEmailBatchQueryParams{
		Page:  int(req.Page),
		Count: int(req.Count),
	}
	mdbEntries, err := mdb.GetEmailBatch(s.db, queryOptions)
	if err != nil {
		return nil, err
	}
	pbEntries := make([]*pb.EmailEntry, 0, len(mdbEntries))
	for _, e := range mdbEntries {
		pbEntry := mdbEntryTopbEntry(&e)
		pbEntries = append(pbEntries, &pbEntry)
	}
	return &pb.GetEmailBatchResponse{
		EmailEntries: pbEntries,
	}, nil
}

func Serve(db *sql.DB, bind string) {
	listener, err := net.Listen("tcp", bind)
	if err != nil {
		log.Fatalf("gRPC Server failure to bind: %v, error: %v\n", bind, err)
	}
	grpcServer := grpc.NewServer()
	mailServer := MailServer{db: db}
	pb.RegisterMailingListServiceServer(grpcServer, &mailServer)
	log.Printf("gRPC server listening on %v\n", bind)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("gRPC Server failure to start: %v, error: %v\n", bind, err)
	}
}
