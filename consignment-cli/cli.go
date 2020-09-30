package main

import (
	pb "shippy/consignment-service/proto/consignment"
	"io/ioutil"
	"encoding/json"
	"errors"
	"log"
	"os"
	"context"
	"github.com/micro/go-micro"
)

const (
	ADDRESS				= "localhost:50051"
	DEFAULT_INFO_FILE	= "consignment.json"
)

//parse file from json
func parseFile(filename string)(*pb.Consignment,error){
	data,err := ioutil.ReadFile(filename)
	if err != nil {
		return nil,err
	}
	var consignment *pb.Consignment
	err = json.Unmarshal(data,&consignment)
	if err != nil {
		return nil,errors.New("json content is error")
	}
	return consignment,nil
}

func main() {
	//connect gRPC
	//conn,err := grpc.Dial(ADDRESS,grpc.WithInsecure())
	//if err != nil {
	//	log.Fatalf("connect err : %v :",err)
	//}
	//defer conn.Close()
	service := micro.NewService(micro.Name("go.micro.srv.consignment"))
	service.Init()
	//initialize gRPC client
	client := pb.NewShippingServiceClient("go.micro.srv.consignment",service.Client())

	infofile := DEFAULT_INFO_FILE
	if len(os.Args) > 1{
		infofile = os.Args[1]
	}

	//parse
	consignment,err := parseFile(infofile)
	if err != nil {
		log.Fatalf("parse info file error :%v",err)
	}

	//invoke gRPC
	resp,err := client.CreateConsignment(context.Background(),consignment)
	if err != nil {
		log.Fatalf("create consignment error %v:",err)
	}
	//log.Printf("created: %t :",resp.Created)
	log.Printf("created: %t", resp.Created)
	//log.Printf("resp: %v", resp)

	//list all consignment
	resp, err = client.GetConsignments(context.Background(), &pb.GetRequest{})
	if err != nil {
		log.Fatalf("failed to list consignments: %v", err)
	}


	for idx, c := range resp.Consignments {
		log.Printf("%+v", c)
		log.Printf("%v",idx)
	}

}
