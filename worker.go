package main

import "fmt"
import "io"
import "os"
import "flag"
import "strings"
import "io/ioutil"
import "os/exec"
import "time"
import "encoding/json"
import zmq "github.com/alecthomas/gozmq"

var ip *string = flag.String("ip", "127.0.0.1", "server IP address")
var desc *string = flag.String("desc", "this is a worker with only basic libraries", "worker description")
var address string

func main() {
  flag.Parse();
  address = fmt.Sprintf("tcp://%s", *ip)
  rcontext, _ := zmq.NewContext()
  rsocket, _ := rcontext.NewSocket(zmq.REQ)
  rsocket.Connect(fmt.Sprintf("%s:16652", address))
  defer rcontext.Close()
  defer rsocket.Close()
  rsocket.Send([]byte(*desc), 0)
  w_id, _ := rsocket.Recv(0)
  worker_id := string(w_id)
  //fmt.Println(worker_id)

  context, _ := zmq.NewContext()
  socket, _ := context.NewSocket(zmq.SUB)
  defer context.Close()
  defer socket.Close()

  socket.SetSubscribe(string(worker_id))
  socket.Connect(fmt.Sprintf("%s:16654", address))

  wcontext, _ := zmq.NewContext()
  wsocket, _ := wcontext.NewSocket(zmq.REQ)
  wsocket.Connect(fmt.Sprintf("%s:16650", address))
  defer wcontext.Close()
  defer wsocket.Close()

  ccontext, _ := zmq.NewContext()
  csocket, _ := ccontext.NewSocket(zmq.REQ)
  csocket.Connect(fmt.Sprintf("%s:16651", address))
  defer ccontext.Close()
  defer csocket.Close()

  var tasks map[string]time.Time = make(map[string]time.Time)

  for {
    datapt, _ := socket.Recv(0)
    st := strings.Replace(string(datapt), "\n", "", -1)
    temps := strings.Split(st, " ")
    if len(temps) < 2 { //todo: check when this happens
      continue 
    }
    if temps[1] == "checkWorker"{
      csocket.Send([]byte("dummy"), 0)
      _, _ = csocket.Recv(0)
    } else {
      cmd := temps[1:]
      var response []byte
      var err error
      var ecmd *exec.Cmd
      if cmd[0] == "execute" {
        ecmd = exec.Command(cmd[1], cmd[2:]...)
	timeNow := time.Now()
        error := ecmd.Start()
        if error != nil {
          response = []byte("error when starting a command")        
        } else {
          response = []byte("command execution started")
	  identifier := strings.Join(cmd[1:], " ")
	  tasks[identifier] = timeNow
        }
        if err != nil {
          //fmt.Println(err)
        }
      } else if cmd[0] == "wait" {
	filename := "parameters.txt"
        f, err := os.Create(filename)
        if err != nil {
          fmt.Println(err)
        }
        n, err := io.WriteString(f, strings.Join(cmd[3:], " "))
        if err != nil {
          fmt.Println(n, err)
        }
        f.Close()
        ecmd = exec.Command(cmd[1], cmd[2])
	response, err = ecmd.Output()
        if err != nil {
          fmt.Println(err)
        }
      } else if cmd[0] == "results" {
	content, err := ioutil.ReadFile(cmd[1])
	if err == nil {
	  fmt.Println(content)
	  response = []byte(content)
	} else {
	  response = []byte("retrieving results failed")
	}
      } else if cmd[0] == "wget" {
        ecmd := exec.Command("wget", cmd[1])
	error := ecmd.Start()
        if error != nil {
          response = []byte("error when trying to download")
        } else {
          response = []byte("downloading started")
	}
      } else if cmd[0] == "status" {
	statusRepr := make(map[string] string)
  	for identifier, _ := range tasks{
	  fmt.Println(identifier)
	  fmt.Println(tasks[identifier])
    	  parts := strings.Split(identifier, " ")
	  fileName := parts[len(parts)-1]
	  info, err := os.Stat(fileName)
	  if err == nil {
	    modified := tasks[identifier].Before(info.ModTime())
	    //fmt.Println(modified)
	    if modified {
	      statusRepr[identifier] = "done"
	    } else {
	      statusRepr[identifier] = "not done yet"
	    }
	  } else {
	      statusRepr[identifier] = "file does not exist"
	  }
	}
	//fmt.Println(statusRepr)
	response, _ = json.Marshal(statusRepr)
      }
      
      wsocket.Send([]byte(response), 0)
      _, _ = wsocket.Recv(0)
    }
  }
}


