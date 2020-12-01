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
  "strings"
  "math/rand"
  pb "github.com/Tarea1/Express/logistica"
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

        fmt.Printf("Splitting to %d pieces.\n", totalPartsNum)

        for i := uint64(0); i < totalPartsNum; i++ {

                partSize := int(math.Min(fileChunk, float64(fileSize-int64(i*fileChunk))))
                partBuffer := make([]byte, partSize)

                file.Read(partBuffer)

                // write to disk
                fileName := "./SplitBooks/"+archivo+"_" + strconv.FormatUint(i+1, 10)
                _, err := os.Create(fileName)

                if err != nil {
                        fmt.Println(err)
                        os.Exit(1)
                }

                // write/save buffer to disk
                ioutil.WriteFile(fileName, partBuffer, os.ModeAppend)

                fmt.Println("Split to : ", fileName)
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
            currentChunkFileName := "./ChunksDownload/"+archivo+"_"+strconv.FormatUint(j+1, 10)

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
        *files = append(*files, path)
        return nil
    }
}
func librosDisponibles (){ //tengo duda dee esto
  var files []string
    root := "./Libros/"
    err := filepath.Walk(root, visit(&files))
    if err != nil {
        panic(err)
    }
    for _, file := range files {
        fmt.Println(file)
    }
}

func mostrarMenu() {
  fmt.Println("Bienvenido Cliente!")
  fmt.Println("Seleccione la acción que desea realizar:")
  fmt.Println("1. Download")
  fmt.Println("2. Upload")
  fmt.Println("3. Salir")
}


func main() {
  //conexion
  var conn *grpc.ClientConn
  conn, err := grpc.Dial(":9000", grpc.WithInsecure())
  if err != nil{
    log.Fatalf("could not connect: %s", err)
  }
  defer conn.Close()
  //

  //holamundo
  c := pb.NewLogisticaServiceClient(conn)
  message := pb.Message{
    Body: "Hello from the client!",
  }
  response, err := c.SayHello(context.Background(),&message)
  if err!= nil{
    log.Fatalf("Error when calling SayHello: %s", err)
  }
  //
  var opcion int;
  var flag_menu = true
  var flag_upload = true
  var algoritmoUp string
  
  for flag_menu {
    mostrarMenu()
    fmt.Scanln(&opcion)
    if opcion == 2 {
      fmt.Printf("Qué tipo de algoritmo de exclusión mutua desea utilizar? [0: Distribuido, 1: Centralizado]:")
      fmt.Scanln(&flag_upload)
      if flag_upload == 0 {
        algoritmoUp = "distribuido"
        } else if flag == 1 {
          algoritmoUp = "centralizado"
          }else {
            log.Printf("Opción inválida")
          }
    }else if opcion == 1{
      fmt.Println("A descargar chicos!!")
      //leer registro name node
    }else if opcion == 3 {
      flag_menu = false
    }else {
      fmt.Println("Ingrese una opción válida")
    }
  }
}
