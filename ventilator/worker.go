package main

import "fmt"
import "flag"
import "strings"
import "os/exec"
import zmq "github.com/alecthomas/gozmq"

var ip *string = flag.String("ip", "127.0.0.1", "server IP address")
var desc *string = flag.String("desc", "this is a worker with only basic libraries", "worker description")
var address string

func main() {
  flag.Parse();
  address = fmt.Sprintf("tcp://%s", *ip)

  serverContext, _ := zmq.NewContext()
  serverSocket, _ := serverContext.NewSocket(zmq.PULL)
  defer serverContext.Close()
  defer serverSocket.Close()
  serverSocket.Connect(fmt.Sprintf("%s:16654", address))

  sinkContext, _ := zmq.NewContext()
  sinkSocket, _ := sinkContext.NewSocket(zmq.PUSH)
  defer sinkContext.Close()
  defer sinkSocket.Close()
  sinkSocket.Connect(fmt.Sprintf("%s:16650", address))

  for {
      datapt, _ := serverSocket.Recv(0)
      st := strings.Replace(string(datapt), "\n", "", -1)
      temps := strings.Split(st, " ")
      cmd := temps[1:]
      var response []byte
      var err error
      var ecmd *exec.Cmd
      if cmd[0] == "execute" {
        ecmd = exec.Command(cmd[1], cmd[2:]...)
        response, err = ecmd.Output()
        if err != nil {
          response = []byte("error when starting a command")        
        }
      }
      sinkSocket.Send(response, 0)
    }
}


