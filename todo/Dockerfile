FROM golang
WORKDIR /go/src/todo

RUN go get github.com/gin-gonic/gin github.com/jinzhu/gorm github.com/jinzhu/gorm/dialects/mysql

COPY . .

RUN go build

ENTRYPOINT ./todo