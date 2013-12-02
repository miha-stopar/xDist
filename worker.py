from subprocess import call
import zmq
from gradientcmd import command

#var ip *string = flag.String("ip", "127.0.0.1", "server IP address")
ip = "127.0.0.1"
#var desc *string = flag.String("desc", "this is a worker with only basic libraries", "worker description")
desc = "this is a worker ..."
tasks = {}

if __name__ == "__main__":
    address = "tcp://%s" % ip
    rcontext = zmq.Context()
    rsocket = rcontext.socket(zmq.REQ)
    rsocket.connect("%s:16652" % address)
    rsocket.send(desc, copy=False)
    w_id = rsocket.recv(copy=False)
    worker_id = str(w_id)
    
    context = zmq.Context()
    socket = context.socket(zmq.SUB)
    socket.connect("%s:16654" % address)
    socket.setsockopt(zmq.SUBSCRIBE, worker_id)
    
    wcontext = zmq.Context()
    wsocket = wcontext.socket(zmq.REQ)
    wsocket.connect("%s:16650" % address)
    
    ccontext = zmq.Context()
    csocket = ccontext.socket(zmq.REQ)
    csocket.connect("%s:16651" % address)
    
    while True: 
        datapt = str(socket.recv(copy=False))
        if datapt[-1] == "\n":
            datapt = datapt[:-1]
        temps = datapt.split(" ")
        if len(temps) < 2:
            continue 
        if temps[1] == "checkWorker":
            csocket.send("dummy", copy=False)
            _ = csocket.recv(copy=False)
        else:
            cmd = temps[1:]
            response = ""
            if cmd[0] == "execute":
                call(cmd[1:])
                response = "command executed"
            elif cmd[0] == "gradient":
                response = command(" ".join(cmd[1:]))
            wsocket.send(response, copy=False)
            _ = wsocket.recv(copy=False)
            
            
            


