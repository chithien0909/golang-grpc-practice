package main

import (
	"contact/contactpb"
	"fmt"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net"

	"google.golang.org/grpc"
)

type server struct {
}

func init()  {
	err := orm.RegisterDriver("mysql", orm.DRMySQL)
	if err != nil {
		log.Panicf("Register driver err %v", err)
	}
	connectStr := "root:123456@tcp(localhost:3306)/demo_contact?charset=utf8"
	err = orm.RegisterDataBase("default", "mysql", connectStr, 100, 100)
	if err != nil {
		log.Panicf("Register db err %v", err)
	}
	orm.RegisterModel(new(ContactInfo))
	_ = orm.RunSyncdb("default", false, true)
	fmt.Println("Connect db successfully")
}
func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:50069")

	s := grpc.NewServer()

	contactpb.RegisterContactServiceServer(s, &server{})

	fmt.Println("Contact is running ....")

	err = s.Serve(lis)

	if err != nil {
		log.Fatalf("Error while server %v", err)
	}

}
