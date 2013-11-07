About
=====

NOTE: this is work in progress.

.. image:: https://raw.github.com/miha-stopar/xNeural/master/img/xcloud.png


Run server
=====

* install Go (for this and the steps below you can check docker/Dockerfile and see which libraries are needed)
* install ZeroMQ
* install gozmq (you will need Git for this):

::

	go get github.com/alecthomas/gozmq
	
* download xNeural

* build server.go:

::

	go build server.go

* run *server* (example command for running server on a local subnet): 

::

	./server -ip=192.168.1.12

Run worker
=====

* install docker
* download xCloud/docker directory
* modify *ip* parameter at the end of the Dockerfile (needs to be the central *server* IP) and SSH username/password (you can as well disable SSH access)
* build docker container from a Dockerfile (execute the following command when in folder *xCloud/docker*):

::

	docker build -t xneural-img .

* run docker container (you might want to limit the CPU and RAM of the container using *-c* and *-m* options):

::

	docker run -d xneural-img

*Worker* will be automatically started. You can connect to the container using SSH:

::

        ssh root@localhost -p 49164

Find out the port number using the command:

::

        docker ps

Run client
=====

There are two possibilities:

Run client from within Docker container:
-------------------------------

* install docker
* download xCloud/docker-client directory
* build docker container from a Dockerfile (execute the following command when in folder *xCloud/docker-client*):

::

	docker build -t xclient .

* run docker container:

::

	docker run -d xclient

* go into Docker container and set GOPATH variable:

::

	export GOPATH=/srv/gocode

* build client.go:

::

	go build client.go

* start *client*

Run client without Docker container:
-------------------------------

* install Go
* install ZeroMQ
* install gozmq
* download xCloud
* build client.go:

::

	go build client.go

* start *client*

How to start and use client
-------------------------------

* run *client* - ip has to be the IP of a *server*: 

::

	./client -ip=192.168.1.12

Note
=====

Use ZeroMQ version 2.2 or higher (due to SetRcvTimeout call in server.go).



