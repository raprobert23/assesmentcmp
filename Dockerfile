FROM golang:latest

RUN mkdir -p /go/src/CMP

WORKDIR /go/src/CMP

COPY . /go/src/CMP

RUN go install CMP

CMD /go/bin/CMP

EXPOSE 8080