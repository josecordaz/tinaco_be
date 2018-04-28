FROM golang:1.9.2

RUN go get github.com/stianeikeland/go-rpio

RUN go get -u github.com/gorilla/mux

COPY main.go /

COPY run.sh /

EXPOSE 8000

WORKDIR /

RUN pwd

ENTRYPOINT ["sh","run.sh"]
