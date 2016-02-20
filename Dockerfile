FROM golang:latest 
RUN apt-get install git

RUN mkdir /app 
ADD . /app/ 
WORKDIR /app 

RUN mkdir -p files
RUN git init files

RUN git config --global user.email "system@dockercontainer"
RUN git config --global user.name "system"

ENV GOPATH /app
ENV GOBIN $GOPATH/bin
RUN go get .
RUN go build -o main . 

ENTRYPOINT ["/app/entrypoint.sh"]

EXPOSE 8080
