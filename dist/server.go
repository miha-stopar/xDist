package main

import "fmt"
import "flag"
import "strconv"
import "strings"
import "time"
import "math/rand"
import zmq "github.com/alecthomas/gozmq"
import "encoding/json"

var workers map[string]string = make(map[string]string) // workerId : description
var statusWorkers map[string]string = make(map[string]string) // workerId : status
var tasksCounter map[string]int = make(map[string]int) // workerId : number of tasks
//var workerTasks map[string]string = make(map[string]string) // workerId : client_id+" "+cmd

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
    tasksCounter[workerId] = 0
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
    //fmt.Println(cmds)
    if cmds[0] == "list"{
        workersRepr :=  make(map[string] string)
        for ind, desc := range workers{
          //tasks := strconv.Itoa(tasksCounter[ind])
	  if statusWorkers[ind] == "connected"{
            //workersRepr[ind] = fmt.Sprintf("tasks: %s | %s", tasks,  desc) 
            workersRepr[ind] = fmt.Sprintf("%s", desc) 
	  }
        } 
        //data, _ := bson.Marshal(workersRepr)
        data, _ := json.Marshal(workersRepr)
        socket.Send(idbla, zmq.SNDMORE) //TODO: reply to ID which actually sent request
        socket.Send(data, 0)
    } else {
      if len(cmds) < 2{
        socket.Send(idbla, zmq.SNDMORE) //TODO: reply to ID which actually sent request
        socket.Send([]byte("not enough arguments"), 0)
      } else {
 	if cmds[1] == "all" {
          for ind, _ := range workers{
            //command := cmds[0] + " " + string(idbla) + " "  + strings.Join(cmds[2:], " ")
            command := cmds[0] + " " + strings.Join(cmds[2:], " ")
	    delegate(command, ind)
	  }
 	} else if cmds[1] == "-1" {
	  workerId := chooseId()
	  /* TODO
	  workerId := "-1"
    	  maxTasks := 1000      
    	  for ind, _ := range tasksCounter{
      	    status := statusWorkers[ind]
      	    if status == "connected"{
              if tasksCounter[ind] < maxTasks {
                workerId = ind
                maxTasks = tasksCounter[ind]
              }
            }
    	  }
	  */
    	  if workerId != "-1" {
            command := cmds[0] + " " + strings.Join(cmds[2:], " ")
	    delegate(command, workerId)
          }
	} else {
          workerId := cmds[1]
          command := cmds[0] + " " + strings.Join(cmds[2:], " ")
	  delegate(command, workerId)
	}
      }
    }     
  } 
}

func chooseId() string{
  // TODO
  workerId := "0"
  c := 0
  for ind, _ := range tasksCounter{
    status := statusWorkers[ind]
    if status == "connected"{
      c = c + 1
    }
  }
  rand.Seed(time.Now().Unix())
  i := rand.Intn(c)

  k := 0
  for ind, _ := range tasksCounter{
    status := statusWorkers[ind]
    if status == "connected"{
      if k == i {
	  workerId = ind
	  break
      }
    }
    k += 1
  }
  if lastId == workerId && c > 1{
	  for ind, _ := range tasksCounter{
	    status := statusWorkers[ind]
	    if status == "connected"{
	      if ind != lastId {
		  workerId = ind
		  break
	      }
	    }
	  }
  }
  
  lastId = workerId
  return workerId
}

func delegate(command string, workerId string){
  msg := fmt.Sprintf("%s %s", workerId, command)
  psocket.Send([]byte(msg), 0)
  tasksCounter[workerId] += 1
  //workerTasks[workerId] //TODO
}

func waitReplies(){
  for {
    reply, err := wsocket.Recv(0)
    wsocket.Send([]byte("dummy"), 0)
    //fmt.Println(string(reply))
    if err != nil {
      //statusWorkers[workerId] = "disconnected"
      socket.Send([]byte("no answer"), 0)
    } else {
      socket.Send(idbla, zmq.SNDMORE) //TODO: reply to ID which actually sent request and update (-1) tasksCounter counter
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

var lastId string = "-1"
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


