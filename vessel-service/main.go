package main

import (
	pb "shippy/vessel-service/proto/vessel"
	"github.com/pkg/errors"
	"context"
	"github.com/micro/go-micro"
	"log"
)

type Repository interface {
	FindAvailable(*pb.Specification)(*pb.Vessel,error)
}

type VesselRepository struct {
	vessels []*pb.Vessel
}

//apply interface
func (repo *VesselRepository) FindAvailable(spec *pb.Specification)(*pb.Vessel,error){
	for _,v := range repo.vessels {
		if v.Capacity >= spec.Capacity && v.MaxWeight >= spec.MaxWeight {
			return v,nil
		}
	}
	return nil,errors.New("No vessel can't be use")
}

//define ship service
type service struct {
	repo Repository
}

//apply server
func (s *service) FindAvailable(ctx context.Context,spec *pb.Specification,resp *pb.Response) error{
	//find inner method for search
	v,err := s.repo.FindAvailable(spec)
	if err != nil {
		return err
	}
	resp.Vessel = v
	return nil
}

func main() {
	vessels := []*pb.Vessel{
		{Id: "vessel001", Name: "Boaty McBoatface", MaxWeight: 200000, Capacity: 500},
	}
	repo := &VesselRepository{vessels:vessels}
	server := micro.NewService(
		micro.Name("go.micro.srv.vessel"),
		micro.Version("latest"),
	)
	server.Init()

	pb.RegisterVesselServiceHandler(server.Server(),&service{repo:repo})

	if err := server.Run();err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
