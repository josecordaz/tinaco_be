FROM golang

RUN go get github.com/stianeikeland/go-rpio

RUN go get -u github.com/gorilla/mux

COPY main.go /

EXPOSE 8000

WORKDIR /

ENTRYPOINT [ "go","run","main.go" ] 