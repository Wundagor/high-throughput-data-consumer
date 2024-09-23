FROM golang:1.23.1

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

COPY create_table.sql /docker-entrypoint-initdb.d/

RUN go build -o /goapp .

CMD ["/goapp"]
