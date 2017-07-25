FROM golang

WORKDIR /go/src/auth

RUN go get github.com/dgrijalva/jwt-go

COPY . .

RUN go build

ENTRYPOINT ./auth