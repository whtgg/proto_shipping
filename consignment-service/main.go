package main

import (
	pb "shippy/consignment-service/proto/consignment"
	"context"
	"github.com/micro/go-micro"
	"log"
	vesselPb "shippy/vessel-service/proto/vessel"
)

const (
	PORT = ":50051"
)

//consignment API
type IRepository interface {
	Create(consignment *pb.Consignment)(*pb.Consignment,error)	//store new cargo
	GetAll() []*pb.Consignment									//get all cargo
}

//repository struct
type Repository struct {
	consignments []*pb.Consignment								//all cargo
}

//achieve
func (repo *Repository) Create(consignment *pb.Consignment)(*pb.Consignment,error){
	repo.consignments = append(repo.consignments,consignment)
	return consignment,nil
}

func (repo *Repository) GetAll()[]*pb.Consignment{
	return repo.consignments
}

//define micro-service
//consignment-service as client use vessel-service --part2
type service struct {
	repo Repository
	vesselClient vesselPb.VesselServiceClient
}

//achieve consignment-service API
//and make service as gRpc server
//consignment new cargo
//func (s *service) CreateConsignment(ctx context.Context,req *pb.Consignment)(*pb.Response,error){
func (s *service) CreateConsignment(ctx context.Context,req *pb.Consignment,resp *pb.Response)error{
	//part2---
	vReq := &vesselPb.Specification{
		Capacity:int32(len(req.Containers)),
		MaxWeight:req.Weight,
	}
	vResp,err := s.vesselClient.FindAvailable(context.Background(),vReq)
	if err != nil {
		return err
	}

	// 货物被承运
	log.Printf("found vessel: %s\n", vResp.Vessel.Name)
	req.VesselId = vResp.Vessel.Id
	//part2--over
	consignment,err := s.repo.Create(req)
	if err != nil {
		log.Printf("create error is %v",err)
		return err
	}
	log.Printf("consignment is %+v \n",consignment)
	//resp = &pb.Response{Created:true,Consignment:consignment}
	resp.Created = true
	resp.Consignment = consignment
	log.Printf("crated is %+v",resp.Created)
	return nil
}

//check consigment-cargo infomation
//func (s *service) GetConsignments(ctx context.Context,req *pb.GetRequest)(*pb.Response,error){
func (s *service) GetConsignments(ctx context.Context,req *pb.GetRequest,resp *pb.Response)error{
	allConsignments := s.repo.GetAll()
	//resp = &pb.Response{Consignments:allConsignments}
	resp.Consignments = allConsignments
	return nil

}


func main() {
	//listener,err := net.Listen("tcp",PORT)
	//if err != nil {
	//	log.Fatalf("failed to listen:%v",err)
	//}
	//log.Printf("listen on: %s\n",PORT)
	//
	//server := grpc.NewServer()
	//repo := Repository{}
	//pb.RegisterShippingServiceServer(server,&service{repo:repo})
	//
	//if err := server.Serve(listener);err != nil {
	//	log.Fatalf("failed to serve: %v", err)
	//}
	server := micro.NewService(
		micro.Name("go.micro.srv.consignment"),
		micro.Version("latest"),
	)

	server.Init()
	repo := Repository{}
	//part2----
	vClient := vesselPb.NewVesselServiceClient("go.micro.srv.vessel",server.Client())
	pb.RegisterShippingServiceHandler(server.Server(),&service{repo,vClient})

	if err := server.Run();err != nil {
		log.Fatalf("failed to serve:",err)
	}
}
