package main

import "fmt"
import "flag"
import "strings"
import zmq "github.com/alecthomas/gozmq"

func serve() {
  context, _ := zmq.NewContext()
  socket, _ := context.NewSocket(zmq.REP)
  socket.Bind(fmt.Sprintf("%s:16653", address))
  defer context.Close()
  defer socket.Close()
  
  for {
    msg, _ := socket.Recv(0)
    cmd := string(msg)
    senderSocket.Send([]byte(cmd), 0)
    socket.Send([]byte("task sent to worker"), 0)
  } 
}

func catchSinkResults() {
  context, _ := zmq.NewContext()
  socket, _ := context.NewSocket(zmq.REP)
  socket.Bind(fmt.Sprintf("%s:16653", address))
  defer context.Close()
  defer socket.Close()
  
  for {
    msg, _ := socket.Recv(0)
    cmd := string(msg)
    //fmt.Println("-----------")
    cmds := strings.Split(cmd, " ")
    //fmt.Println(cmds)
    if cmds[0] == "list"{
        socket.Send([]byte("not available for ventilator"), 0)
    } else {
      if len(cmds) < 1{
        socket.Send([]byte("not enough arguments"), 0)
      } else {
	senderSocket.Send([]byte(cmd), 0)
      }
    }     
  } 
}


var senderSocket *zmq.Socket
var sinkSocket *zmq.Socket
var ip *string = flag.String("ip", "127.0.0.1", "public IP address of this very computer")
var address string

func main() {
  flag.Parse();
  address = fmt.Sprintf("tcp://%s", *ip)
  senderContext, _ := zmq.NewContext()
  senderSocket, _ = senderContext.NewSocket(zmq.PUSH)
  defer senderContext.Close()
  defer senderSocket.Close()
  senderSocket.Bind(fmt.Sprintf("%s:16654", address))

  sinkContext, _ := zmq.NewContext()
  sinkSocket, _ = sinkContext.NewSocket(zmq.PUSH)
  defer sinkContext.Close()
  defer sinkSocket.Close()
  sinkSocket.Connect(fmt.Sprintf("%s:16650", address))
  sinkSocket.Send([]byte("0"), 0) // start batch

  go serve()
  var inp string
  fmt.Scanln(&inp)
}


