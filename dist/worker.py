import sys
from subprocess import check_output, call
import zmq
from gradientcmd import command

# python worker.py 192.168.1.12

#ip = "127.0.0.1"
ip = "198.101.154.21"
address = "tcp://%s" % ip
desc = "this is a worker ..."
tasks = {}
imported_modules = {}
        
if __name__ == "__main__":
    if len(sys.argv) == 2:
        ip = sys.argv[1]
    
    rcontext = zmq.Context()
    rsocket = rcontext.socket(zmq.REQ)
    rsocket.connect("%s:16652" % address)
    rsocket.send(desc, copy=False)
    worker_id = str(rsocket.recv(copy=False))

    context = zmq.Context()
    socket = context.socket(zmq.SUB)
    socket.connect("%s:16654" % address)
    socket.setsockopt(zmq.SUBSCRIBE, worker_id)

    wcontext = zmq.Context()
    wsocket = wcontext.socket(zmq.DEALER)
    wsocket.connect("%s:16650" % address)

    ccontext = zmq.Context()
    csocket = ccontext.socket(zmq.REQ)
    csocket.connect("%s:16651" % address)

    while True:
	try:
            msg = str(socket.recv(copy=False))
	    if msg[-1] == "\n":
            	msg = msg[:-1]
            cmd = msg.split(" ")
	    cmd = cmd[1:] # remove worker_id
	    #if cmd[0] != "checkWorker":
	    #	print cmd[0]

	    response = ""
	    if cmd[0] == "checkWorker":
		csocket.send("dummy", copy=False)
	        csocket.recv(copy=False)
		continue
            elif cmd[0] == "addCommand":
                source_url = cmd[1]
                file_name = source_url.split("/")[-1]
                file_name = file_name[:-3] # remove ".py"
                call(["wget", source_url])
                module = __import__(file_name)
                imported_modules[file_name] = module           
            elif cmd[0] == "execute":
                response = check_output(cmd[1:])
	    elif cmd[0] == "gradient":
                response = command(" ".join(cmd[1:]))
            else:
                # assumptions: 
                # - command name is the same as downloaded file name (via addCommand)
                # - file that was downloaded contains function "command"
                file_name = cmd[0]
                response = imported_modules[file_name].__dict__["command"](cmd[1:])
	    wsocket.send(response, copy=False)
	except Exception as e:
	    print e

           


