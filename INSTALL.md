Dependencies
============

* go get -u github.com/astaxie/beego
* go get -u github.com/beego/bee
* go get -u github.com/govend/govend


Update vendors
==============

Run command
`govend -v -u -l`
if you want to update dependencies


Database
========
Change connection string at config/app.conf to connect to the Redis database.
Current value of host is pointing to the Redis, that was set up in the docker.


Run server
==========
Do shell command `bee run`
or `./tournamentAPI`, if you has not installed the bee application


Running in the Docker
=====================
If you has the Docker-compose program on your computer, you can run the project with all dependencies in the container.

Go to the `docker` directory, using command
`cd docker`

and run container, using command
`docker-compose up --build`


Enjoy the API on address http://localhost:8081