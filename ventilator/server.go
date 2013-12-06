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
    cmds := strings.Split(cmd, " ")
    //fmt.Println(cmds)
    if cmds[0] == "all"{ // publish to all workers
	c := strings.Join(cmds[1:], " ")
        psocket.Send([]byte(c), 0)
        socket.Send([]byte("dummy"), 0)
    } else {
      if len(cmds) < 1{
        socket.Send([]byte("not enough arguments"), 0)
      } else {
	senderSocket.Send([]byte(cmd), 0)
        socket.Send([]byte("dummy"), 0)
      }
    }     
  } 
}

var psocket *zmq.Socket
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

  pcontext, _ := zmq.NewContext()
  psocket, _ = pcontext.NewSocket(zmq.PUB)
  defer pcontext.Close()
  defer psocket.Close()
  psocket.Bind(fmt.Sprintf("%s:16651", address))

  go serve()
  var inp string
  fmt.Scanln(&inp)
}


