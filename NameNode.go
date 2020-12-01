package main

import (
  //"os"
  //"strconv"
  "fmt"
  "log"
  "net"
  "math"
  //"io/ioutil"
  //"path/filepath"
  //"bufio"
  "golang.org/x/net/context"
  "google.golang.org/grpc"
  pb "Tareita2/logisticaName"

)

type Papi struct{
}

func(s *Papi) MandarPropuestaName(ctx context.Context, propuesta *pb.PropuestaName) (*pb.ReplyPropuestaName,error){
  n_random= rand.Intn(100)
  res = 1
  if n_random <20 {
    res = 0
    //generar una nueva propuesta y la reenvia
    return &pb.HelloReply{ReplyName: 0, NuevaProp: nuevaProp}, nil
  }else{
    var temp []string
    temp = append(temp,"")
    return &pb.HelloReply{ReplyName: 1 , NuevaProp: temp }, nil
  }


}

func main() {
  fmt.Println("NameNode encendido")
  lis,err := net.Listen("tcp",":9000")
  if err!= nil {
    log.Fatalf("Failed to listen on port 9000: %v", err)
  }

  s := Papi{}


  grpcServer:= grpc.NewServer()

  pb.RegisterLogisticaServiceServer(grpcServer, &s)

  if err := grpcServer.Serve(lis); err!=nil{
    log.Fatalf("Failed to serve gRPC server over port 9000: %v", err)
  }

}
