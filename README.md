# go-wiki

A _KISS_ wiki built on golang with git as the storage back-end. Content
is formatted in [markdown
syntax](http://daringfireball.net/projects/markdown/syntax). The wiki is
rendered with go templates and [bootstrap](http://getbootstrap.com) css.

Current running example: [wiki.mehl.no](https://wiki.mehl.no/)

## Build and run locally

Ensure that go is installed. Export the GOPATH environment variable to
where you checked out the g-wiki project:

    export GOPATH=$GOPATH:/some/path/g-wiki/

Download dependencies and compile the binary by:

    go get -a
    go build -o wiki .
    
Create a git repository in the files folder:

    git init files/

You can now run g-wiki with the standard settings by executing the
binary:

    ./wiki -local=":8080"

Point your web browser to `http://localhost:8080/` to see the wiki in
action. The wiki tries to store files in a `files` folder within the
project directory. This folder has to exist and be writeable by the user
running the g-wiki instance.

## Docker

Ensure that docker is installed. The docker file will create a `files` directory for you, and initialize a git repository there. Rembember that these files are not persistent. Dependent on your environement run docker as root (or not) and execute the following commands:

    docker build -t go-wiki:latest .

If this executes succesfully your container is ready:

    docker run -d -p 8080:8080 go-wiki:latest
    
This starts the web application in deamon mode, and the application should be accessible on `http://localhost:8080/`
