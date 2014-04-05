# g-wiki

A KISS wiki built on golang with git as the storage back-end.

## Install

Ensure that go is installed. Export the GOPATH environment variable to
where you checked out the g-wiki project:

    export GOPATH=$GOPATH:/some/path/g-wiki/

Create a binary in ./bin by installing g-wiki:

    go install mehl.no/wiki

You can now run g-wiki with the standard settings by executing the
binary:

    ./bin/wiki

Point your web browser to http://localhost:8080/ to see the wiki in
action.
