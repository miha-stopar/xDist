import sys
from subprocess import call
import zmq
from gradientcmd import command

# python worker.py 192.168.1.12

#var ip *string = flag.String("ip", "127.0.0.1", "server IP address")
ip = "127.0.0.1"
#var desc *string = flag.String("desc", "this is a worker with only basic libraries", "worker description")
desc = "this is a worker ..."
tasks = {}

if __name__ == "__main__":
    if len(sys.argv) == 2:
        ip = sys.argv[1]
    address = "tcp://%s" % ip
    serverContext = zmq.Context()
    serverSocket = serverContext.socket(zmq.PULL)
    serverSocket.connect("%s:16654" % address)
    
    sinkContext = zmq.Context()
    sinkSocket = sinkContext.socket(zmq.PUSH)
    sinkSocket.connect("%s:16650" % address)
   
    while True: 
        datapt = str(serverSocket.recv(copy=False))
        if datapt[-1] == "\n":
            datapt = datapt[:-1]
        cmd = datapt.split(" ")
        response = ""
        if cmd[0] == "execute":
            call(cmd[1:])
            response = "command executed"
        elif cmd[0] == "gradient":
            response = command(" ".join(cmd[1:]))
        sinkSocket.send(response, copy=False)
            
            
            


