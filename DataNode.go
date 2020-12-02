package main

import (
  "os"
  "strconv"
  "fmt"
  "log"
  "net"
  //"math"
  "io/ioutil"
  //"path/filepath"
  //"bufio"
  "golang.org/x/net/context"
  "google.golang.org/grpc"
  pb "Tareita2/logistica"
  pb2 "Tareita2/logisticaName"

)

type Papi struct{
}

func generarPropuesta(/*address string, longitud int*/) []string {
  port := "9000" //
  var retorno []string
  //retorno = append(retorno,address)
  for i := 61; i < 64; i++ {
    address := "dist" + strconv.Itoa(i) +":"+ port
    conn, err := grpc.Dial(address, grpc.WithInsecure())
    if err != nil {
      fmt.Println(err)
    }
    defer conn.Close()

    message := pb.HelloRequest{
      Mensaje: "Estas disponible para recibir distribucion de chunks?",
    }
    c := pb.NewLogisticaServiceClient(conn)
    response , err := c.SayHello(context.Background(),&message)
    if err!= nil{
      fmt.Println("Error al conectar: DataNode ",i," no disponible" )
    }else{
      fmt.Println(response)
      fmt.Println("DataNode", i, "en línea")
      /*if len(retorno) != longitud {
        retorno = append(retorno, address)
      }*/

    }

  }
  return retorno
}

func(s *Papi) SayHello(ctx context.Context, message *pb.HelloRequest) (*pb.HelloReply,error){ //verifica si el nodo esta encendido
  log.Printf("Mensaje recibido del cliente: %s", message.Mensaje)
  return &pb.HelloReply{Mensaje: "Estoy disponible!"}, nil
}

func distribuirChunks(propuesta []string, libro [][]byte, titulo string, address string, cantidad int32){
  fmt.Println("Se comienza a distribuir los chunks")
  fmt.Println("Len(libro): ",len(libro))
  fmt.Println("Len(Propuesta): ",len(propuesta))
  fmt.Println("loquesobra: ",len(libro) % len(propuesta))
  fmt.Println("cuantocadauno: ",(len(libro)-loquesobra) / len(propuesta))
  loquesobra := len(libro) % len(propuesta)
  cuantas_cada_uno := (len(libro)-loquesobra) / len(propuesta) //cuantas para cada nodo
  index := cuantas_cada_uno * len(propuesta) //las que se reparten inicialmente
  cont := 0
  for j := 0; j < len(propuesta) ; j++ {
    for i := 0; i < cuantas_cada_uno; i++ {
      if j == 0 && i == 0 {
        EscribirEnNameNode(titulo,cont+1,propuesta[j],cantidad,true)
        EnviarChunk(propuesta[j],libro[cont],titulo,cont+1)
      }else {
        EscribirEnNameNode(titulo,cont+1,propuesta[j],cantidad,false)
        EnviarChunk(propuesta[j],libro[cont],titulo,cont+1)
      }
      cont++
      //para cada dirección de la propuesta, se le mandan 'cuantas_cada_uno' chunks
      //enviar a address 'j'
    }
  }
  for i := 0;  i < loquesobra; i++ {
    EscribirEnNameNode(titulo,index+1,propuesta[i],cantidad,false)
    EnviarChunk(propuesta[i],libro[index],titulo,index+1)
    index++
  }
}


func(s *Papi) MandarChunk(ctx context.Context, SendChunk *pb.SendChunk) (*pb.ReplySendChunk, error) {
  titulo := SendChunk.GetTitulo()
  chunk := SendChunk.GetChunk()
  parte  := SendChunk.GetParte()

  _, err := os.Create("Partes/" + titulo+"_"+strconv.Itoa(int(parte)))
  if err != nil {
		os.Exit(1)
  }
  ioutil.WriteFile("Partes/" + titulo, chunk, os.ModeAppend)

  return &pb.ReplySendChunk{Status: true}, nil
}

func EnviarChunk(address string, chunk []byte, titulo string, parte int) {
  conn, err := grpc.Dial(address, grpc.WithInsecure())
  if err != nil {
    fmt.Println("did not connect: %v", err)
  }
  defer conn.Close()
  c := pb.NewLogisticaServiceClient(conn)
  fmt.Println("Enviando chunk a '", address,"''")
  estadito, _ := c.MandarChunk(context.Background(), &pb.SendChunk{
    //campos que se enviaran entre dataNodes
    Titulo: titulo,
    Chunk: chunk,
    Parte: int32(parte),
  })
  if estadito.GetStatus(){
    fmt.Println("Chunk enviado con éxito")
  }else{
    fmt.Println("Ocurrió un problema en el envío del chunk")
  }

}

func EscribirEnNameNode(titulo string, parte int, address string, cantidad int32, esPrimero bool)  {
  conn, err := grpc.Dial("dist64:9000", grpc.WithInsecure())
  if err != nil {
    fmt.Println("did not connect: %v", err)
  }
  defer conn.Close()
  c := pb2.NewLogisticaNameServiceClient(conn)
  estadito, _ := c.MandarLog(context.Background(), &pb2.LogMsg{
    NombreLibro: titulo,
    CantidadPartes: strconv.Itoa(int(cantidad)),
    Parte: strconv.Itoa(parte),
    IpMaquina: address,
    EsPrimero: esPrimero,
  })
  if estadito.GetRecibido(){
    fmt.Println("Se registró correctamente la parte ",strconv.Itoa(parte)," del libro ",titulo," en el LOG")
  }


}

func(s *Papi) SubirLibro(ctx context.Context, dataLibro *pb.Libro) (*pb.SubirLibroReply,error){ // recibe info del libro y sus chunks
  //longitud := 3

  //fmt.Println("Respuesta Len de partes DataNode:", i)
  titulo := dataLibro.GetTitulo();
  cantidad := dataLibro.GetLength(); // cuantas partes se dividio
  chunks := dataLibro.GetChunks();
  algoritmo := dataLibro.GetAlgoritmo() // 0:distribuido 1:centralizado
  address := dataLibro.GetIp() // ip de maquina actual

  if algoritmo == 1 {
    prop := generarPropuesta(/*address,longitud*/)
    conn, err := grpc.Dial("dist64:9000", grpc.WithInsecure())
    if err != nil {
      fmt.Println("did not connect: %v", err)
    }
    defer conn.Close()

    c := pb2.NewLogisticaNameServiceClient(conn)
    estadito, _ := c.MandarPropuestaName(context.Background(), &pb2.PropuestaName{
      Propuesta: prop,
    })

    if estadito.GetReplyName() == 1 {// 1 : prop anterior 0: sacar nueva propuesta
      fmt.Println("Se aceptó la propuesta en el NameNode, se realizará:", prop)
      distribuirChunks(prop, chunks, titulo, address, cantidad)
    } else {
      fmt.Println("Se rechazó la propuesta en el NameNode, se realizará: ",estadito.GetNuevaProp())
      distribuirChunks(estadito.GetNuevaProp(), chunks, titulo, address, cantidad)
    }

  }else{ // algoritmo distribuido
    fmt.Println("lolerio")
  }
  return &pb.SubirLibroReply{Status:cantidad}, nil //Devuelve el largo del array de chunks recibidos
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
