package main

import (
	"awesomeProject/proto"
	pf "awesomeProject/proto"
	"context"
	"crypto/sha256"
	"fmt"
	"github.com/linxlib/logs"
	"google.golang.org/grpc"
	"math"
	"net"
)

const (
	port = ":50051"
)

type server struct {
	pf.UnimplementedFileServiceServer
}

func (s *server) verifyFile(file []byte, hash string, length int64) bool {
	h := sha256.New()
	h.Write(file)
	myHash := fmt.Sprintf("%x", h.Sum(nil))
	logs.Info("hash: ", hash, " myHash: ", myHash, " len: ", length, " myLen: ", len(file))
	return hash == myHash
}

func (s *server) Upload(ctx context.Context, in *proto.FSReq) (*proto.FSRep, error) {
	//if !s.verifyFile(in.File, in.Hash, in.Filelen) {
	//	return &proto.FSRep{
	//		Status:  false,
	//		Message: "verify the hash of data failed, please retry",
	//	}, nil
	//}
	return &proto.FSRep{
		Status:  true,
		Message: "received",
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		logs.Fatalf("failed to listen: %v", err)
	}
	//c, err := credentials.NewServerTLSFromFile("./server.crt", "./server.key")
	//if err != nil {
	//	logs.Fatalf("credentials.NewServerTLSFromFile err: %v", err)
	//}

	s := grpc.NewServer(
		//grpc.Creds(c),
		grpc.MaxRecvMsgSize(math.MaxInt64))

	proto.RegisterFileServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		logs.Fatalf("failed to serve: %v", err)
	}
}
