import zmq

c = zmq.Context()
s = c.socket(zmq.REQ)
s.connect('tcp://127.0.0.1:16653')

msg = "train 2 2 2"
s.send(msg, copy=False)
msg2 = s.recv(copy=False)
print msg2

