package main

import "fmt"
import "flag"
import "strconv"
import "strings"
import "time"
import zmq "github.com/alecthomas/gozmq"
import "encoding/json"

var workers map[string]string = make(map[string]string) // workerId : description
var statusWorkers map[string]string = make(map[string]string) // workerId : status
var tasksWorkers map[string]int = make(map[string]int) // workerId : number of tasks

func waitRegistrations(){
  rcontext, _ := zmq.NewContext()
  rsocket, _ := rcontext.NewSocket(zmq.REP)
  rsocket.Bind(fmt.Sprintf("%s:16652", address))
  defer rcontext.Close()
  defer rsocket.Close()
  for {
    workerDesc, _ := rsocket.Recv(0)
    count := len(workers)
    workerId := strconv.Itoa(count)
    println(workerId)
    workers[string(workerId)] = string(workerDesc)
    statusWorkers[workerId] = "connected"
    tasksWorkers[workerId] = 0
    println("Got worker: ", string(workerDesc))
    rsocket.Send([]byte(workerId), 0)
  }
}

func serve() {
  for {
    idbla, _ = socket.Recv(0)
    msg, _ := socket.Recv(0)
    cmd := string(msg)
    cmds := strings.Split(cmd, " ")
    fmt.Println(cmds)
    if cmds[0] == "list"{
        workersRepr :=  make(map[string] string)
        for ind, desc := range workers{
          tasks := strconv.Itoa(tasksWorkers[ind])
	  fmt.Println(statusWorkers[ind])
	  if statusWorkers[ind] == "connected"{
            workersRepr[ind] = fmt.Sprintf("tasks: %s | %s", tasks,  desc) 
	  }
        } 
        //data, _ := bson.Marshal(workersRepr)
        data, _ := json.Marshal(workersRepr)
        socket.Send(data, 0)
    } else {
      if len(cmds) < 2{
        socket.Send([]byte("not enough arguments"), 0)
      } else {
        workerId := cmds[1]
        command := cmds[0] + " " + strings.Join(cmds[2:], " ")
        msg := fmt.Sprintf("%s %s", workerId, command)
        psocket.Send([]byte(msg), 0)
      }
    }     
  } 
}

func waitReplies(){
  for {
    reply, err := wsocket.Recv(0)
    wsocket.Send([]byte("dummy"), 0)
    fmt.Println("--------------")
    fmt.Println(string(reply))
    if err != nil {
      //statusWorkers[workerId] = "disconnected"
      socket.Send([]byte("no answer"), 0)
    } else {
      //tasksWorkers[workerId] += 1
      socket.Send(idbla, zmq.SNDMORE) //TODO: reply to ID who actually sent request
      socket.Send([]byte(reply), 0)
    }
  }
}

func checkWorkers(){
  for {
    time.Sleep(1000 * time.Millisecond)
    for ind, _ := range statusWorkers{
      msg := fmt.Sprintf("%s %s", ind, "checkWorker")
      psocket.Send([]byte(msg), 0)
      _, err := csocket.Recv(0)
      csocket.Send([]byte("dummy"), 0)
      if err != nil{
	statusWorkers[ind] = "disconnected"
      } else {
        statusWorkers[ind] = "connected"
      }
    }
  } 
}

var idbla []byte
var socket *zmq.Socket
var psocket *zmq.Socket
var wsocket *zmq.Socket
var csocket *zmq.Socket
var ip *string = flag.String("ip", "127.0.0.1", "public IP address of this very computer")
var address string

func main() {
  flag.Parse();
  address = fmt.Sprintf("tcp://%s", *ip)

  context, _ := zmq.NewContext()
  socket, _ = context.NewSocket(zmq.ROUTER)
  socket.Bind(fmt.Sprintf("%s:16653", address))
  defer context.Close()
  defer socket.Close()

  pcontext, _ := zmq.NewContext()
  psocket, _ = pcontext.NewSocket(zmq.PUB)
  defer pcontext.Close()
  defer psocket.Close()
  psocket.Bind(fmt.Sprintf("%s:16654", address))

  wcontext, _ := zmq.NewContext() // connected to workers
  wsocket, _ = wcontext.NewSocket(zmq.DEALER)
  //wsocket.SetRcvTimeout(1000 * time.Millisecond)
  defer wcontext.Close()
  defer wsocket.Close()
  wsocket.Bind(fmt.Sprintf("%s:16650", address))

  ccontext, _ := zmq.NewContext() // connected to workers
  csocket, _ = ccontext.NewSocket(zmq.REP)
  csocket.SetRcvTimeout(1000 * time.Millisecond)
  defer ccontext.Close()
  defer csocket.Close()
  csocket.Bind(fmt.Sprintf("%s:16651", address))

  go waitRegistrations()
  go serve()
  go checkWorkers()
  go waitReplies()

  var inp string
  fmt.Scanln(&inp)
}


