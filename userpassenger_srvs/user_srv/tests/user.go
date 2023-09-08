package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"google.golang.org/grpc"

	"gotrains/userpassenger_srvs/user_srv/model"
	"gotrains/userpassenger_srvs/user_srv/proto"
)

var userClient proto.UserClient
var conn *grpc.ClientConn

func Init() {
	var err error
	conn, err = grpc.Dial("127.0.0.1:50051", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	userClient = proto.NewUserClient(conn)
}

func TestGetUserList() {
	rsp, err := userClient.GetUserList(context.Background(), &proto.PageInfo{
		Pn: 1,
		Ps: 5,
	})
	if err != nil {
		panic(err)
	}
	for _, user := range rsp.Data {
		fmt.Println(user.Mobile, user.NickName, user.Password)
		checkRsp, err := userClient.CheckPassWord(context.Background(), &proto.PasswordCheckInfo{
			Password:          "admin123",
			EncryptedPassword: user.Password,
		})
		if err != nil {
			panic(err)
		}
		fmt.Println(checkRsp.Success)
	}
}

func TestAddPassenger() {
	rand.NewSource(time.Now().UnixNano())
	var typ model.PassengerType
	if rand.Intn(2) == 0 {
		typ = model.PassengerTypeChild
	} else {
		typ = model.PassengerTypeAdult
	}
	for i := 0; i < 10; i++ {
		rep, err := userClient.AddPassenger(context.Background(), &proto.PassengerInfo{
			UserId: int32(rand.Intn(3) + 1),
			Name:   fmt.Sprintf("passenger_%d", i),
			IdCard: fmt.Sprintf("36232520021563456%d", i),
			Type:   int64(typ),
		})
		if err != nil {
			panic(err)
		}
		fmt.Println(rep.Id)
	}
}

func TestCreateUser() {
	for i := 0; i < 10; i++ {
		rsp, err := userClient.CreateUser(context.Background(), &proto.CreateUserInfo{
			NickName: fmt.Sprintf("user_name_%d", i),
			Mobile:   fmt.Sprintf("1878222222%d", i),
			Password: "admin123",
		})
		if err != nil {
			panic(err)
		}
		fmt.Println(rsp.Id)
	}
}

func TestGetPassengerList() {
	rsp, err := userClient.GetPassengerList(context.Background(), &proto.PassengerPageInfo{
		UserId: 1,
	})
	if err != nil {
		panic(err)
	}
	for _, passenger := range rsp.Data {
		fmt.Println(passenger.Name, passenger.IdCard, passenger.Type)
	}
}

func main() {
	Init()
	// TestCreateUser()
	// TestGetUserList()
	// TestAddPassenger()
	TestGetPassengerList()
	conn.Close()
}
