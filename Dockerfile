FROM golang:latest 

RUN mkdir /app 
ADD . /app/ 
WORKDIR /app 
ENV GOPATH /app
ENV GOBIN $GOPATH/bin
RUN go get .
RUN go build -o main . 

ENTRYPOINT ["/app/entrypoint.sh"]

EXPOSE 8080
