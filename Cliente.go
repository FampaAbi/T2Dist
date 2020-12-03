package main

import (
  "strings"
  
  "os"
  "strconv"
  "fmt"
  "log"
  "math/rand"
  "math"
  "io/ioutil"
  "path/filepath"
  "bufio"
  "golang.org/x/net/context"
  "google.golang.org/grpc"
  pb "Tareita2/logistica"
  pb2 "Tareita2/logisticaName"
)

func get_chunks(direccion string,titulo string) bool{
  var conn *grpc.ClientConn
  conn, err := grpc.Dial(direccion, grpc.WithInsecure())
  if err != nil {
    log.Fatalf("did not connect: %v", err)
  }
  defer conn.Close()
  c := pb.NewLogisticaServiceClient(conn)
  estadito, _ := c.GetChunk(context.Background(), &pb.Data{
    Titulo: titulo,
  })
  if estadito.GetStatus(){
    fmt.Println("Parte ", titulo," recuperada correctamente" )
  }else{
    fmt.Println("Ocurrió un problema recuperando ", titulo)
  }
  //dejarlos en chunks download
  // write to disk
  fileName := "./ChunksDownload/"+titulo
  _, err1 := os.Create(fileName)
  if err1 != nil {
          fmt.Println(err)
          os.Exit(1)
  }
  file, err2 := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
			if err2 != nil {
				fmt.Println(err2)
			} else {
					_, err3 := file.Write(estadito.GetChunk())

					if err3 != nil {
						fmt.Println(err3)
						os.Exit(1)
					}
					file.Sync()
				  file.Close()
			}
  return true
} // rescata un chunk de cierta address y lo deja en ./ChunksDownload

func split_chunks(titulo string)(int) {

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
}//https://www.socketloop.com/tutorials/golang-recombine-chunked-files-example

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

                  //fmt.Println("Appending at position : [", writePosition, "] bytes")
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

                  _, err1 := file.Write(chunkBufferBytes)

                  if err1 != nil {
                    fmt.Println(err)
                    os.Exit(1)
                  }

                  file.Sync() //flush to disk

                  // free up the buffer for next cycle
                  // should not be a problem if the chunk size is small, but
                  // can be resource hogging if the chunk size is huge.
                  // also a good practice to clean up your own plate after eating

                  chunkBufferBytes = nil // reset or empty our buffer

                  //fmt.Println("Written ", n, " bytes")

                  //fmt.Println("Recombining part [", j, "] into : ", newFileName)
                }
                file.Close()
}

func visit(files *[]string) filepath.WalkFunc {
    return func(path string, info os.FileInfo, err error) error {
        if err != nil {
            log.Fatal(err)
        }
        *files = append(*files, info.Name())
        return nil
    }
}//https://flaviocopes.com/go-list-files/ funciones para mostrar libros a subir

func remove(s []int, i int) []int {
    s[len(s)-1], s[i] = s[i], s[len(s)-1]
    return s[:len(s)-1]
}//borrar de un array https://yourbasic.org/golang/delete-element-slice/

func searchAvailableNode(conn *grpc.ClientConn) string {
  port := "9000" //
  opciones := []int{61,62,63} //datanode
  var inLoop = true
  for inLoop {
    n := len(opciones)
    if n ==0{
      inLoop = false
    }else{
      n_random := rand.Intn(n)
      random := opciones[n_random]
      address := "dist" + strconv.Itoa(random) +":"+ port
      opciones = remove(opciones, n_random)
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
        fmt.Println("Error al conectar: DataNode ",random," no disponible" )
      }else{
        log.Printf("Response from DataNode: %s", response.Mensaje)
        fmt.Println("DataNode", random, "en línea")
        return address
      }
    }
  }
  return ""
}

func agruparChunks(book_name string, total_partes int) [][]byte {
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
  fmt.Println("3. Ver disponibilidad ")
  fmt.Println("4. Salir")
}

func verDisponibilidadLibros() ([]string, []int32, []string, []string) {
  var conn *grpc.ClientConn
  conn, err := grpc.Dial("dist64:9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
  c := pb2.NewLogisticaNameServiceClient(conn)
  estadito, _ := c.Disponibilidad(context.Background(), &pb2.Mensaje{
    Mensaje: true,
  })

  titulos := estadito.GetTitulos()
  cantidad_partes := estadito.GetCantidadPartes()
  subtitulos := estadito.GetSubtitulos()
  direcciones := estadito.GetAddress()

  if len(titulos) == 0 {
    fmt.Println("No hay libros disponibles actualmente")
    return titulos,cantidad_partes,subtitulos,direcciones
  }else{
    fmt.Println("Los libros disponibles son:")
    for i := 0; i < len(titulos); i++ {
      fmt.Println(i+1,".",titulos[i])
    }
  }

  return titulos,cantidad_partes,subtitulos,direcciones
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
        tituloUP := librosUpload() //muestra la lista de opciones a subir
        if tituloUP == "exit" {
          inUpload = false
        } else {
          partes := split_chunks(tituloUP) //separa el libro en chunks y devuelve el total de partes
          lista_de_bytes := agruparChunks(tituloUP, partes) //mete todos los chunks en un array

        fmt.Printf("Qué tipo de algoritmo de exclusión mutua desea utilizar? [0: Distribuido, 1: Centralizado]:")
        fmt.Scanln(&opcionUp)

        address := searchAvailableNode(conn) //encontrar a que nodo mandarle los chunks inicialmente
        //fmt.Println(address)
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

        if estadito.GetStatus() == int32(partes){
          fmt.Println("Respuesta: Se subió correctamente el libro!")
        }else{
          fmt.Println("Respuesta: Se produjo un error en la subida del libro!")
        }
        }
      }
    } else if opcion == 1{
      fmt.Println("Elija un libro a descargar:")
      titulos,cantidad_partes,subtitulos,direcciones := verDisponibilidadLibros()
      base := 0
      if len(titulos) != 0{
        var que_libro int
        fmt.Scanln(&que_libro)
        title := titulos[que_libro-1]
        if que_libro < 0 || que_libro > len(titulos) {
          fmt.Println("Opcion inválida")
        }else {
          for i := 0; i < que_libro-1; i++ {
             base += int(cantidad_partes[i])
          }
          for i := 0; i < int(cantidad_partes[que_libro-1]); i++ {
            parte := strings.Split(subtitulos[i + base], "_")[2]
            fmt.Println("IP: ",direcciones[i + base])
            fmt.Println("Title: ",title+"_"+parte)
            fmt.Println("Parte: ",parte)
            get_chunks(direcciones[i+ base],title+"_"+parte)

          }
          join_chunks(title,int(cantidad_partes[que_libro-1]))
          fmt.Println("Libro descargado correctamente")
        }
      }
    } else if opcion == 4 {
      inMenu = false
    } else if opcion == 3 { // ver disponibilidad
      verDisponibilidadLibros()
    }else {
      fmt.Println("Ingrese una opción válida")
    }
  }
}
