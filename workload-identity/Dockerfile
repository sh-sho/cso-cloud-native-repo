FROM golang:1.21.1

WORKDIR /app

COPY go.mod gp.sum ./
RUN go mod download
COPY main.go ./

EXPOSE 8080
CMD ["go", "run", "main.go"]
