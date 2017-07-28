FROM golang

WORKDIR /go/src/auth

RUN go get github.com/go-sql-driver/mysql golang.org/x/crypto/bcrypt github.com/dgrijalva/jwt-go


COPY . .

RUN go build

ENTRYPOINT ./auth