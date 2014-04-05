# g-wiki

A _KISS_ wiki built on golang with git as the storage back-end. Content
is formatted in [markdown
syntax](http://daringfireball.net/projects/markdown/syntax).

## Install

Ensure that go is installed. Export the GOPATH environment variable to
where you checked out the g-wiki project:

    export GOPATH=$GOPATH:/some/path/g-wiki/

Download dependencies and compile the binary by:

    go get all
    go install mehl.no/wiki

You can now run g-wiki with the standard settings by executing the
binary:

    ./bin/wiki

Point your web browser to `http://localhost:8080/` to see the wiki in
action. The wiki tries to store files in a `files` folder within the
project directory. This folder has to exist and be writeable by the user
running the g-wiki instance.
