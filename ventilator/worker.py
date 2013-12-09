import sys
from subprocess import call
import threading
import zmq
from gradientcmd import command

# python worker.py 192.168.1.12

#ip = "127.0.0.1"
ip = "198.101.154.21"
address = "tcp://%s" % ip
desc = "this is a worker ..."
tasks = {}
imported_modules = {}

class Subscriber(threading.Thread): # listening for commands to be executed on all workers
    def __init__(self):
        threading.Thread.__init__(self)
        context = zmq.Context()
        self.socket = context.socket(zmq.SUB)
        self.socket.connect("%s:16651" % address)
        self.socket.setsockopt(zmq.SUBSCRIBE, "")

    def run(self):
        while True:
            msg = str(self.socket.recv(copy=False))
            cmd = msg.split(" ")
	    print "+++++"
	    print cmd
            if cmd[0] == "addCommand":
                source_url = cmd[1]
                file_name = source_url.split("/")[-1]
                file_name = file_name[:-3] # remove ".py"
                call(["wget", source_url])
                module = __import__(file_name)
                imported_modules[file_name] = module           
            elif cmd[0] == "execute":
                call(cmd[1:])

if __name__ == "__main__":
    if len(sys.argv) == 2:
        ip = sys.argv[1]
    subscriber = Subscriber()
    subscriber.daemon = True
    subscriber.start()
    serverContext = zmq.Context()
    serverSocket = serverContext.socket(zmq.PULL)
    serverSocket.connect("%s:16654" % address)
    
    sinkContext = zmq.Context()
    sinkSocket = sinkContext.socket(zmq.PUSH)
    sinkSocket.connect("%s:16650" % address)
   
    while True: 
        datapt = str(serverSocket.recv(copy=False))
	if len(datapt) == 0:
	    continue
        if datapt[-1] == "\n":
            datapt = datapt[:-1]
        cmd = datapt.split(" ")
	print "----"
	print cmd
        response = ""
        if cmd[0] == "execute":
            call(cmd[1:])
            response = "command executed"
        elif cmd[0] == "gradient":
            response = command(" ".join(cmd[1:]))
        else:
            # assumptions: 
            # - command name is the same as downloaded file name (via addCommand)
            # - file that was downloaded contains function "command"
            file_name = cmd[0]
            response = imported_modules[file_name].__dict__["command"](cmd[1:])
        sinkSocket.send(response, copy=False)
            
            
            


