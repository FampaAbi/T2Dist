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

func distribuirChunks(propuesta string[], libro byte[][], titulo string, address string, cantidad int) {
  cuantas_cada_uno := len(libro) / len(propuesta)  
  loquesobra := len(libro) % len(propuesta)
  
  for j := 0; j < len(propuesta) j++ {
    for i := 0; i < cuantas_cada_uno; i++ {
      if j == 0 && i == 0 {
        EscribirEnNameNode(titulo,j+i,propuesta[j],cantidad,true)
        EnviarChunk(propuesta[j],libro[j+i],titulo,j+i)
      }else {
        EscribirEnNameNode(titulo,j+i,propuesta[j],cantidad,false)
        EnviarChunk(propuesta[j],libro[j+i],titulo,j+i)
      }
      //para cada dirección de la propuesta, se le mandan 'cuantas_cada_uno' chunks
      //enviar a address 'j' 
    }
  } 
}

func EnviarChunk(address string, chunk byte[], titulo string, parte int) {
  conn, err := grpc.Dial(address, grpc.WithInsecure())
  if err != nil {
    fmt.Println("did not connect: %v", err)
  }
  defer conn.Close
  c := pb.NewLogisticaNameServiceClient(conn)
  estadito, _ := c.MandarChunk(context.Background(), &pb.SendChunk{
    //campos que se enviaran entre dataNodes
    Titulo: titulo,
    Chunk: chunk,
    Parte: int32(parte),
  })
  fmt.Println("Recibido?:", estadito.GetRecibido())
  }
}

func EscribirEnNameNode(titulo string, chunk int, address string, cantidad int, esPrimero bool)  {
  conn, err := grpc.Dial("dist64:9000", grpc.WithInsecure())
  if err != nil {
    fmt.Println("did not connect: %v", err)
  }
  defer conn.Close
  c := pb.NewLogisticaNameServiceClient(conn)
  estadito, _ := c.MandarLog(context.Background(), &pb.LogMsg{
    NombreLibro: titulo,
    CantidadPartes: strconv.Itoa(cantidad),
    Parte: strconv.Itoa(chunk),
    IpMaquina: address,
    EsPrimero: esPrimero,
  })
  fmt.Println("Recibido?:", estadito.GetStatus())
  }
}

func(s *Papi) SubirLibro(ctx context.Context, dataLibro *pb.Libro) (*pb.SubirLibroReply,error){ // recibe info del libro y sus chunks
  len := 3
  i := len(dataLibro.GetChunks())
  titulo := dataLibro.GetTitulo();
  length := dataLibro.GetLength();
  chunks := dataLibro.GetChunks();
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
    if estadito.GetReplyName() == 1 {
      distribuirChunks(prop, chunks, titulo, address, cantidad) 
    } else {
      distribuirChunks(estadito.GetNuevaProp(), chunks, titulo, address, cantidad)
    }

  }else{ // algoritmo distribuido
    continue
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
