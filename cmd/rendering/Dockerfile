FROM golang:1.9.2 as builder
WORKDIR /home/src
COPY . /home/src
RUN go get github.com/lib/pq
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o run .

FROM golang:1.9.2
WORKDIR /root/
EXPOSE 80
COPY --from=builder /home/src/run .
COPY --from=builder /home/src/static/ ./static/
CMD ["./run"]