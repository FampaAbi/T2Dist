package main

import (
  //"os"
  //"strconv"
  "fmt"
  "log"
  "net"
  //"math"
  //"io/ioutil"
  //"path/filepath"
  //"bufio"
  "golang.org/x/net/context"
  "google.golang.org/grpc"
  pb "Tareita2/logistica"

)

type Papi struct{
}

func(s *Papi) SayHello(ctx context.Context, message *pb.HelloRequest) (*pb.HelloReply,error){
  log.Printf("Mensaje recibido del cliente: %s", message.Mensaje)
  return &pb.HelloReply{Mensaje: "Estoy disponible!"}, nil
}

func(s *Papi) SubirLibro(ctx context.Context, dataLibro *pb.Libro) (*pb.SubirLibroReply,error){
  i := len(dataLibro.GetChunks())
  return &pb.SubirLibroReply{Status:int32(i)}, nil
}

func main() {
  fmt.Println("Datanode encendido")
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
