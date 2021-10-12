FROM golang:latest

RUN mkdir /build
WORKDIR /build

RUN export GO111MODULE=on

COPY go.mod /build
COPY go.sum /build/

RUN cd /build/ && git clone https://github.com/Restitutor-Orbis/DISYS-MiniProject2.git
RUN cd /build/DISYS-MiniProject2/server

EXPOSE 9080

ENTRYPOINT [ "/build/DISYS-MiniProject2/server/server" ]