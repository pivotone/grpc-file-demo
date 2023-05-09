package main

import (
	"awesomeProject/proto"
	"context"
	"crypto/sha256"
	"fmt"
	"github.com/linxlib/conv"
	"github.com/linxlib/logs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io/ioutil"
	"time"
)

const (
	address = "localhost:50051"
)

func main() {
	//c, err := credentials.NewClientTLSFromFile("./server.crt", "deploy")
	//if err != nil {
	//	logs.Fatalf("credentials.NewClientTLSFromFile err: %v", err)
	//}
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logs.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := proto.NewFileServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)

	defer cancel()

	bs, _ := ioutil.ReadFile("./redis.pdf")
	filelen := conv.Int64(len(bs))

	h := sha256.New()
	h.Write(bs)
	myhash := fmt.Sprintf("%x", h.Sum(nil))
	logs.Info("myhash: ", myhash)

	start := time.Now()
	r, err := client.Upload(ctx, &proto.FSReq{
		DstDir:   "main",
		ProjName: "grpc-file-demo",
		Name:     "redis",
		ProjType: 1,
		Hash:     myhash,
		Filelen:  filelen,
		IfReboot: false,
		File:     bs,
	})
	end := time.Now().Sub(start).Seconds()
	kb := filelen / 1024
	logs.Info("time: ", end, " file size: ", kb, "KB")
	if err != nil {
		logs.Fatalf("could not upload: %v", err)
	}
	logs.Printf("upload: %s", r.Message)
}
