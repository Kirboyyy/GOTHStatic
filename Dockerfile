FROM golang:latest
WORKDIR /app
COPY . .
RUN go mod download
RUN go install github.com/a-h/templ/cmd/templ@latest
RUN go generate
RUN go build -o main .
ENV GIN_MODE=release
EXPOSE 8080
CMD ["./main"]
