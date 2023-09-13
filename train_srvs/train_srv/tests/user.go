package main

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"

	"gotrains/train_srvs/train_srv/proto"
)

// var userClient proto.UserClient
var conn *grpc.ClientConn
var trainClient proto.TrainClient

func Init() {
	var err error
	conn, err = grpc.Dial("127.0.0.1:50052", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	trainClient = proto.NewTrainClient(conn)
}

func TestGenTrainDaily() {
	dt, _ := time.Parse("2006-01-02", "2023-09-20")
	rsp, err := trainClient.GenDaily(context.Background(), &proto.DateRequest{
		Date: dt.Format("2006-01-02"),
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp)
}

func main() {
	Init()
	TestGenTrainDaily()
	// TestCreateUser()
	// TestGetUserList()
	// TestAddPassenger()
	// TestGetPassengerList()
	conn.Close()
}
