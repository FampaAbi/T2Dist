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

func split_chunks(titulo string)(int) { //https://www.socketloop.com/tutorials/golang-recombine-chunked-files-example

        fileToBeChunked := "./Libros/"+titulo // change here! path

        file, err := os.Open(fileToBeChunked)

        if err != nil {
                fmt.Println(err)
                os.Exit(1)
        }

        defer file.Close()

        fileInfo, _ := file.Stat()

        var fileSize int64 = fileInfo.Size()

        const fileChunk = 250000 // 1 MB, change this to your requirement

        // calculate total number of parts the file will be chunked into

        totalPartsNum := uint64(math.Ceil(float64(fileSize) / float64(fileChunk)))

        fmt.Printf("Separando en %d partes.\n", totalPartsNum)

        for i := uint64(0); i < totalPartsNum; i++ {

                partSize := int(math.Min(fileChunk, float64(fileSize-int64(i*fileChunk))))
                partBuffer := make([]byte, partSize)

                file.Read(partBuffer)

                // write to disk
                fileName := "./SplitBooks/"+titulo+"_" + strconv.FormatUint(i+1, 10)
                _, err := os.Create(fileName)

                if err != nil {
                        fmt.Println(err)
                        os.Exit(1)
                }

                // write/save buffer to disk
                ioutil.WriteFile(fileName, partBuffer, os.ModeAppend)

                //fmt.Println("Split to : ", fileName)
        }
        return int(totalPartsNum)
  }

func join_chunks(titulo string,totalPartsNum int){
    newFileName := "./JoinBooks/"+titulo
    _, err := os.Create(newFileName)
    if err != nil {
            fmt.Println(err)
            os.Exit(1)
          }
          //set the newFileName file to APPEND MODE!!
          // open files r and w
          file, err := os.OpenFile(newFileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
          if err != nil {
            fmt.Println(err)
            os.Exit(1)
          }
          // IMPORTANT! do not defer a file.Close when opening a file for APPEND mode!
          // defer file.Close()

          // just information on which part of the new file we are appending

          var writePosition int64 = 0
          for j := uint64(0); j < uint64(totalPartsNum); j++ {
            //read a chunk
            currentChunkFileName := "./ChunksDownload/"+titulo+"_"+strconv.FormatUint(j+1, 10)

            newFileChunk, err := os.Open(currentChunkFileName)

            if err != nil {
                    fmt.Println(err)
                    os.Exit(1)
                  }

                  defer newFileChunk.Close()

                  chunkInfo, err := newFileChunk.Stat()

                  if err != nil {
                    fmt.Println(err)
                    os.Exit(1)
                  }

                  // calculate the bytes size of each chunk
                  // we are not going to rely on previous data and constant

                  var chunkSize int64 = chunkInfo.Size()
                  chunkBufferBytes := make([]byte, chunkSize)

                  fmt.Println("Appending at position : [", writePosition, "] bytes")
                  writePosition = writePosition + chunkSize

                  // read into chunkBufferBytes
                  reader := bufio.NewReader(newFileChunk)
                  _, err = reader.Read(chunkBufferBytes)

                  if err != nil {
                    fmt.Println(err)
                    os.Exit(1)
                  }

                  // DON't USE ioutil.WriteFile -- it will overwrite the previous bytes!
                  // write/save buffer to disk
                  //ioutil.WriteFile(newFileName, chunkBufferBytes, os.ModeAppend)

                  n, err := file.Write(chunkBufferBytes)

                  if err != nil {
                    fmt.Println(err)
                    os.Exit(1)
                  }

                  file.Sync() //flush to disk

                  // free up the buffer for next cycle
                  // should not be a problem if the chunk size is small, but
                  // can be resource hogging if the chunk size is huge.
                  // also a good practice to clean up your own plate after eating

                  chunkBufferBytes = nil // reset or empty our buffer

                  fmt.Println("Written ", n, " bytes")

                  fmt.Println("Recombining part [", j, "] into : ", newFileName)
                }
                file.Close()
}

//https://flaviocopes.com/go-list-files/ funciones para mostrar libros a subir
func visit(files *[]string) filepath.WalkFunc {
    return func(path string, info os.FileInfo, err error) error {
        if err != nil {
            log.Fatal(err)
        }
        *files = append(*files, info.Name())
        return nil
    }
}

func librosUpload() string{
  var files []string
    root := "./Libros/"
    err := filepath.Walk(root, visit(&files))
    if err != nil {
        panic(err)
    }
    var i int
    var libro int
	  i = 0
    fmt.Println("Seleccione el libro que desea subir:")
    for _, file := range files {
        if i !=0{
          fmt.Println(i, ".", file)
        }
        i++
    }
    fmt.Println(i, ". Salir")
    fmt.Scanln(&libro)
    if libro == i {
      return "exit"
    } else if libro > i || libro < 1 {
      fmt.Println("Opcion invalida, hazlo bien sapo culiao")
      return librosUpload()
    }
    return files[libro]
} //retorna el libro

func mostrarMenu() {
  fmt.Println("Bienvenido Cliente!")
  fmt.Println("Seleccione la acción que desea realizar:")
  fmt.Println("1. Download")
  fmt.Println("2. Upload")
  fmt.Println("3. Salir")
}



func searchAvailableNode(conn *grpc.ClientConn) string { //agregar sayhello
  port := "9000" //
  for i := 61; i < 64; i++ {
    address := "dist" + strconv.Itoa(i) +":"+ port
    conn, err := grpc.Dial(address, grpc.WithInsecure())
    if err != nil {
      fmt.Println(err)
    }
    defer conn.Close()

    message := pb.HelloRequest{
      Mensaje: "Hello from cliente",
    }
    c := pb.NewLogisticaServiceClient(conn)
    response, err := c.SayHello(context.Background(),&message)
    if err!= nil{
      fmt.Println("Error al conectar: DataNode ",i," no disponible" )
    }else{
      log.Printf("Response from DataNode: %s", response.Mensaje)
      fmt.Println("DataNode", i, "en línea")
      return address
    }
  }
  return ""
}

func getChunks(book_name string, total_partes int) [][]byte {
  var retorno [][]byte
  root := "./SplitBooks/"
  for i := 1; i < total_partes + 1; i++ {
    file, err := os.Open(root + book_name + "_" + strconv.Itoa(i))
	  if err != nil {
		log.Fatal(err)
	  }
    content, err := ioutil.ReadAll(file)
	  if err != nil {
		  fmt.Println("Error!:", err)
	  }
    retorno = append(retorno, content) //posible error
  }
  return retorno
}


func main() {
  //conexion
  var conn *grpc.ClientConn
  //
  var opcion int;
  var opcionUp int;
  var inMenu = true

  for inMenu {
    mostrarMenu()
    fmt.Scanln(&opcion)
    if opcion == 2 {
      var inUpload = true
      for inUpload {
        tituloUP := librosUpload()
        if tituloUP == "exit" {
          inUpload = false
        } else {
          partes := split_chunks(tituloUP)
          lista_de_bytes := getChunks(tituloUP, partes)

        fmt.Printf("Qué tipo de algoritmo de exclusión mutua desea utilizar? [0: Distribuido, 1: Centralizado]:")
        fmt.Scanln(&opcionUp)

        address := searchAvailableNode(conn)
        fmt.Println(address)
        conn, err := grpc.Dial(address, grpc.WithInsecure())
        if err != nil {
          fmt.Println("did not connect: %v", err)
        }
        defer conn.Close()

        c := pb.NewLogisticaServiceClient(conn)
        estadito, _ := c.SubirLibro(context.Background(), &pb.Libro{
          Titulo: tituloUP,
          Length: int32(partes),
          Chunks: lista_de_bytes,
          Ip: address,
          Algoritmo: int32(opcionUp),
        })
        fmt.Println("Respuesta:", estadito)
        }

      }
    } else if opcion == 1{
      fmt.Println("A descargar chicos!!")
      //leer registro name node
    } else if opcion == 3 {
      inMenu = false
    } else {
      fmt.Println("Ingrese una opción válida")
    }
  }
}
