package main

import (
  "os"
  "strconv"
  "fmt"
  "log"
  "math"
  "io/ioutil"
  "path/filepath"
  "bufio"
  "golang.org/x/net/context"
  "google.golang.org/grpc"
  pb "Tareita2/logistica"

)
type Papi struct{
}

func(s *Papi) SayHello(ctx context.Context, message *pb.HelloRequest) (*pb.HelloReply,error){
  log.Printf("Received message body from client: %s", message.Mensaje)
  return &pb.HelloReply{Mensaje: "Hello From DataNode!"}, nil
}

func main() {
  //var conn *grpc.ClientConn
  //conn, err := grpc.Dial(":9000", grpc.WithInsecure())
  //if err != nil{
  //  log.Fatalf("could not connect: %s", err)
  //}
  //defer conn.Close()

  //c := pb.NewLogisticaServiceClient(conn)

fmt.Println("Datanode encendido")

}
