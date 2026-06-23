FROM golang:1.25.8-bookworm

WORKDIR /app

COPY . .

RUN go mod tidy
RUN GOOS=linux GOARCH=amd64 go build -o /app/exe ./cmd

CMD [ "/app/exe" ]