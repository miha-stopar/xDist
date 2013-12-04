package main

import "fmt"
import "flag"
import "bufio"
import "os"
import "strings"
import "encoding/json"
import zmq "github.com/alecthomas/gozmq"

func enterCmd(socket *zmq.Socket){
  reader := bufio.NewReader(os.Stdin)
  fmt.Print("Enter command: ")
  command, _ := reader.ReadString('\n')
  parts := strings.Split(string(command), " ")
  if strings.Contains(parts[0], "list") {
    c := "list"
    socket.Send([]byte(c), 0)
    reply, _ := socket.Recv(0)
    var output map[string]string
    _ = json.Unmarshal(reply, &output)
    fmt.Println("Workers:")
    for k,v := range output{
      fmt.Printf("%s: %s\n\n", string(k), v)
    }
  } else {
    socket.Send([]byte(command), 0)
    reply, _ := socket.Recv(0)
    fmt.Println(string(reply))
    /*
    if strings.Contains(parts[0], "results") {
      _ = json.Unmarshal(reply, &output)
      for k,v := range output{
        fmt.Printf("%s: %s\n\n", string(k), v)
      }
    } else {
      fmt.Println(string(reply))
    }
    */
  }   
  enterCmd(socket)
}

var workerId string
var ip *string = flag.String("ip", "127.0.0.1", "server IP address")
var address string 

func main() {
  flag.Parse();
  context, _ := zmq.NewContext()
  socket, _ := context.NewSocket(zmq.REQ)
  address = fmt.Sprintf("tcp://%s", *ip)
  add := fmt.Sprintf("%s:16653", address)
  socket.Connect(add)
  enterCmd(socket)
}
