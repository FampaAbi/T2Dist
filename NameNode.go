package main

import (
  "os"
  "strconv"
  "fmt"
  "log"
  "net"
  "math/rand"
  "errors"
  //"io/ioutil"
  //"path/filepath"
  //"bufio"
  "golang.org/x/net/context"
  "google.golang.org/grpc"
  pb2 "Tareita2/logisticaName"
  //pb "Tareita2/logistica"

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
  var retorno []string
  for i := 0; i < largo; i++ {
    n_random := rand.Intn(largo-i)
    random := opciones[n_random]
    address := "dist" + random +":"+ port
    opciones = remove(opciones, n_random)
    retorno = append(retorno, address)
  }
  return retorno
}

func AceptaroRechazar() bool {
  n_random:= rand.Intn(100)
  if n_random < 20 {
    return false
  }
  return true
}

func(s *Papi) MandarPropuestaName(ctx context.Context, propuesta *pb2.PropuestaName) (*pb2.ReplyPropuestaName,error){
  largo_propuesta := len(propuesta.GetPropuesta())
  believer := propuesta.GetPropuesta() //propuesta

  if largo_propuesta != 1 {
    if AceptaroRechazar() { //acepta
      var temp []string
      temp = append(temp,"")
      return &pb2.ReplyPropuestaName{ReplyName: 1 , NuevaProp: temp }, nil
    } else { //rechaza
      if largo_propuesta > 2 {
        nueva_prop := generarPropuesta(believer, 2)
        return &pb2.ReplyPropuestaName{ReplyName: 0, NuevaProp: nueva_prop}, nil
      }else if largo_propuesta > 1 {
        nueva_prop := generarPropuesta(believer, 1)
        return &pb2.ReplyPropuestaName{ReplyName: 0, NuevaProp: nueva_prop}, nil
      }
    }
  } else { // si la propuesta es de largo 1 siempre se acepta
    var temp []string
    temp = append(temp,"")
    return &pb2.ReplyPropuestaName{ReplyName: 1 , NuevaProp: temp }, nil
  }
  myErr := errors.New("fallo de pana")
  var temp []string
  temp = append(temp,"fallo de pana")
  return &pb2.ReplyPropuestaName{ReplyName: 0 , NuevaProp: temp }, myErr
}

func(s *Papi) MandarLog(ctx context.Context, LogMsg *pb2.LogMsg) (*pb2.ReplyLogMsg, error) {

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
  return &pb2.ReplyLogMsg{Recibido: true}, nil
}

func escribirEnLog(nombre_libro string, cantidad_partes string, parte string, ip string, esPrimero bool, libro_actual string) {

  file, err := os.OpenFile("./LOG.txt", os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		log.Println(err)
  }
  if esPrimero {
	  if _, err := file.WriteString(nombre_libro + " " + cantidad_partes + "\n"); err != nil {
	  	log.Fatal(err)
    }
  }
  if _, err := file.WriteString("parte_"+libro_actual+"_"+parte+" "+ip+"\n"); err != nil {
    log.Fatal(err)
  }
  fmt.Println("Se registró correctamente la parte ",parte," del libro ",nombre_libro," en el LOG")
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

  pb2.RegisterLogisticaNameServiceServer(grpcServer, &s)

  if err := grpcServer.Serve(lis); err!=nil{
    log.Fatalf("Failed to serve gRPC server over port 9000: %v", err)
  }

}
