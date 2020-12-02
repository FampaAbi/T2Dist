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
type Papi struct{ //struct que maneja la info general de la logistica (EL PAPI)
  libro int
}


func remove(s []string, i int) []string { //borrar de un array https://yourbasic.org/golang/delete-element-slice/
  s[len(s)-1], s[i] = s[i], s[len(s)-1]
  return s[:len(s)-1]
}

func generarPropuesta(opciones []string, largo int) []string {
  port := "9000" //
  retorno = []
  for i := 0; i < largo; i++ {
    n_random := rand.Intn(largo-i)
    random := opciones[n_random]
    address := "dist" + strconv.Itoa(random) +":"+ port
    opciones = remove(opciones, n_random)
    retorno = append(retorno, address)
  }
  return retorno
}

func AceptaroRechazar() bool {
  n_random= rand.Intn(100)
  if n n_random < 20 {
    return false
  }
  return true
}

func(s *Papi) MandarPropuestaName(ctx context.Context, propuesta *pb.PropuestaName) (*pb.ReplyPropuestaName,error){
  largo_propuesta := len(propuesta.GetPropuesta())  
  believer := propuesta.GetPropuesta()

  if largo_propuesta != 1 {
    if AceptaroRechazar() {
      var temp []string
      temp = append(temp,"")
      return &pb.HelloReply{ReplyName: 1 , NuevaProp: temp }, nil
    } else {
      if largo_propuesta > 2 {
        nueva_prop := generarPropuesta(believer, 2)
        return &pb.HelloReply{ReplyName: 0, NuevaProp: nueva_prop}, nil
      }
      else if largo_propuesta > 1 {
        nueva_prop := generarPropuesta(believer, 1)
        return &pb.HelloReply{ReplyName: 0, NuevaProp: nueva_prop}, nil
      }
    }  
  } else {
    var temp []string
    temp = append(temp,"")
    return &pb.HelloReply{ReplyName: 1 , NuevaProp: temp }, nil
  }

}

func(s *Papi) MandarLog(ctx context.Context, LogMsg *pb.LogMsg) (*pb.ReplyLogMsg, error) {
  
  nombre_libro := LogMsg.GetNombreLibro()
  cantidad_partes := LogMsg.GetCantidadPartes()
  parte := LogMsg.GetParte()
  ip := LogMsg.GetIpMaquina()
  esPrimero := LogMsg.GetEsPrimero()
  if esPrimero {
    s.libro++
  }

  libro_actual := strconv.Itoa(s.libro)

  escribirEnLog(nombre_libro,cantidad_partes, parte,ip, esPrimero, libro_actual)
  return &pb.ReplyLogMsg{Recibido: true}, nil
}

func(s *Papi) MandarChunk(ctx context.Context, SendChunk *pb.SendChunk) (*pb.ReplySendChunk, error) {
  titulo := SendChunk.GetTitulo() 
  chunk := SendChunk.GetChunk()
  parte  := SendChunk.GetParte()

  _, err := os.Create("Partes/" + titulo)
  if err != nil {
		os.Exit(1)
  }
  ioutil.WriteFile("Partes/" + titulo, chunk, os.ModeAppend)

  return &pb.ReplySendChunk{Status: true}, nil
}


func escribirEnLog(nombre_libro string, cantidad_partes string, parte string, ip string, esPrimero string, libro_actual string) {

  file, err := os.OpenFile("./LOG.txt", os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		log.Println(err)
  }
  if esPrimero {
	  if _, err := file.WriteString(FileName + " " + strconv.Itoa(int(in.TotalParts)) + "\n"); err != nil {
	  	log.Fatal(err)
    }
  }
  if _, err := file.WriteString("parte_"+libro_actual+"_"+parte+" "+ip+"\n"); err != nil {
    log.Fatal(err)
  }
}

func main() {
  fmt.Println("NameNode encendido")
  lis,err := net.Listen("tcp",":9000")
  if err!= nil {
    log.Fatalf("Failed to listen on port 9000: %v", err)
  }

  s := Papi{}
  s.libro = 0

  grpcServer:= grpc.NewServer()

  pb.RegisterLogisticaServiceServer(grpcServer, &s)

  if err := grpcServer.Serve(lis); err!=nil{
    log.Fatalf("Failed to serve gRPC server over port 9000: %v", err)
  }

}
