package main

import (
  "os"
  "errors"
  "strconv"
  "fmt"
  "log"
  "net"
  "math/rand"
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

func(s *Papi) GetChunk(ctx context.Context, data *pb.Data) (*pb.ReplyChunk,error){ //verifica si el nodo esta encendido
  titulo := data.Titulo
  var chunk []byte
  root := "./Partes/"
  file, err := os.Open(root + titulo)
  if err != nil {
  log.Fatal(err)
  }
  content, err := ioutil.ReadAll(file)
  if err != nil {
    fmt.Println("Error!:", err)
  }
  chunk  = content
  return &pb.ReplyChunk{Status: true,Chunk: chunk}, nil
}

func remove(s []string, i int) []string { //borrar de un array https://yourbasic.org/golang/delete-element-slice/
  s[len(s)-1], s[i] = s[i], s[len(s)-1]
  return s[:len(s)-1]
}

func generarPropuestaDist(opciones []string, largo int) []string {
  port := "9000" //
  var retorno []string
  for i := 0; i < largo; i++ {
    n_random := rand.Intn(largo-i)
    random := opciones[n_random]
    address := "dist" + random +":"+ port
    opciones = remove(opciones, n_random)
    retorno = append(retorno, address)
  }
  fmt.Println("La nueva propuesta es: ",retorno)
  return retorno
}

func Distribuido() ([]string) { //devuelve una prop válida
  prop := generarPropuesta()
  for i := 0; i < len(prop); i++ {
    conn, err := grpc.Dial(prop[i], grpc.WithInsecure())
    if err != nil {
      fmt.Println("did not connect: %v", err)
    }
    defer conn.Close()
    c := pb.NewLogisticaServiceClient(conn)
    estadito, _ := c.MandarPropuesta(context.Background(), &pb.Propuesta{
    Propuesta: prop,
    })
    if !estadito.GetReplyName() { //si es que rechaza la prop
      i = 0
      fmt.Println("Se rechazó la propuesta por el DataNode con IP '",prop[i],"', se realizará otra propuesta. ")
      prop = generarPropuestaDist(prop, len(prop)-1)
    }else{
      fmt.Println("Se aceptó la propuesta por el DataNode con IP '",prop[i],"'")
    }
  }
  return prop
}

func(s *Papi) MandarPropuesta(ctx context.Context, propuesta *pb.Propuesta) (*pb.ReplyPropuesta,error){
  largo_propuesta := len(propuesta.GetPropuesta())
  //believer := propuesta.GetPropuesta() //propuesta

  if largo_propuesta != 1 {
    if AceptaroRechazar() { //acepta
      var temp []string
      temp = append(temp,"")
      return &pb.ReplyPropuesta{ReplyName: true}, nil
    } else { //rechaza
        return &pb.ReplyPropuesta{ReplyName: false}, nil
    }
  } else { // si la propuesta es de largo 1 siempre se acepta
    var temp []string
    temp = append(temp,"")
    return &pb.ReplyPropuesta{ReplyName: true }, nil
  }
  myErr := errors.New("fallo de pana")
  var temp []string
  temp = append(temp,"fallo de pana")
  return &pb.ReplyPropuesta{ReplyName: false}, myErr
}

func AceptaroRechazar() bool { //Rechaza con 10 % de prob
  n_random:= rand.Intn(100)
  if n_random < 10 {
    return false
  }
  return true
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
      fmt.Println("DataNode", i, "es agregado a la propuesta")
      retorno = append(retorno, address)
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
  ioutil.WriteFile("Partes/" + titulo+"_"+strconv.Itoa(int(parte)), chunk, os.ModeAppend)

  return &pb.ReplySendChunk{Status: true}, nil
}

func EnviarChunk(address string, chunk []byte, titulo string, parte int) {
  conn, err := grpc.Dial(address, grpc.WithInsecure())
  if err != nil {
    fmt.Println("did not connect: %v", err)
  }
  defer conn.Close()
  c := pb.NewLogisticaServiceClient(conn)
  fmt.Println("Enviando chunk a '", address,"'")
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
    fmt.Println("Escritura Correcta")
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

  } else{ // algoritmo distribuido
        prop := Distribuido()
        fmt.Println("Se aceptó la propuesta por todos los DataNodes, se realizará:", prop)
        distribuirChunks(prop, chunks, titulo, address, cantidad)
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
