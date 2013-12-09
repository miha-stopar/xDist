package main

import "fmt"
import "flag"
import zmq "github.com/alecthomas/gozmq"

func waitResults() {
  for {
    msg, _ := sinkSocket.Recv(0)
    rSocket.Send(msg, 0)
    rSocket.Recv(0)
  } 
}

var sinkSocket *zmq.Socket
var rSocket *zmq.Socket
var ip *string = flag.String("ip", "127.0.0.1", "public IP address of this very computer")
var address string

func main() {
  flag.Parse();
  address = fmt.Sprintf("tcp://%s", *ip)
  
  sinkContext, _ := zmq.NewContext()
  sinkSocket, _ = sinkContext.NewSocket(zmq.PULL)
  defer sinkContext.Close()
  defer sinkSocket.Close()
  sinkSocket.Bind(fmt.Sprintf("%s:16650", address))
  sinkSocket.Recv(0)

  rContext, _ := zmq.NewContext()
  rSocket, _ = rContext.NewSocket(zmq.REQ) // client won't be able to Bind to this socket, if sink is on a remote machine
  defer rContext.Close()
  defer rSocket.Close()
  rSocket.Connect(fmt.Sprintf("%s:16652", address))

  go waitResults()
  var inp string
  fmt.Scanln(&inp)
}


