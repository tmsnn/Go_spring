FROM golang:latest
WORKDIR test/

COPY ./ ./
RUN go build -o main
CMD ["./main"]