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

func generarPropuesta(address string, len int) []string {
  port := "9000" //
  retorno = [address]
  for i := 61; i < 64; i++ {
    address := "dist" + strconv.Itoa(i) +":"+ port
    conn, err := grpc.Dial(address, grpc.WithInsecure())
    if err != nil {
      fmt.Println(err)
    }
    defer conn.Close()

    message := pb.HelloRequest{
      Mensaje: "Estas disponible?",
    }
    c := pb.NewLogisticaServiceClient(conn)
    response, err := c.SayHello(context.Background(),&message)
    if err!= nil{
      fmt.Println("Error al conectar: DataNode ",i," no disponible" )
    }else{
      fmt.Println("DataNode", i, "en línea")
      if len(retorno) != len {
        retorno = append(retorno, address)
      }

    }

  }
  return retorno
}

func(s *Papi) SayHello(ctx context.Context, message *pb.HelloRequest) (*pb.HelloReply,error){ //verifica si el nodo esta encendido
  log.Printf("Mensaje recibido del cliente: %s", message.Mensaje)
  return &pb.HelloReply{Mensaje: "Estoy disponible!"}, nil
}

func(s *Papi) SubirLibro(ctx context.Context, dataLibro *pb.Libro) (*pb.SubirLibroReply,error){ // recibe info del libro y sus chunks
  len := 3
  i := len(dataLibro.GetChunks())
  algoritmo := dataLibro.GetAlgoritmo() // 0:distribuido 1:centralizado
  address := dataLibro.GetIp() // ip de maquina actual
  if algoritmo == 1 {
       prop := generarPropuesta(address,len)
       conn, err := grpc.Dial("dist64:9000", grpc.WithInsecure())
       if err != nil {
         fmt.Println("did not connect: %v", err)
       }
       defer conn.Close()

       c := pb.NewLogisticaNameServiceClient(conn)
       estadito, _ := c.MandarPropuestaName(context.Background(), &pb.PropuestaName{
         Propuesta: prop,
       })
       fmt.Println("Respuesta:", estadito) // 1 : prop anterior 0: sacar nueva propuesta

  }else{ // algoritmo distribuido // si rechazan 2 prpuestas se acpeta siempre la 3ra que es dejar todo en este datanode

  }
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
