FROM golang:alpine

RUN mkdir /app
WORKDIR /app

COPY go.mod .
COPY go.mod .
COPY . .
RUN go get -u github.com/gin-gonic/gin
RUN go get -u gorm.io/gorm
RUN go get -u gorm.io/driver/mysql

EXPOSE 8000

CMD ["go", "run", "."]
