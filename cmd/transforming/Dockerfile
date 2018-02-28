FROM golang:1.9.2 as builder
WORKDIR /home/src
COPY . /home/src
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o run .

FROM golang:1.9.2
WORKDIR /root/
COPY --from=builder /home/src/run .
CMD ["./run"]