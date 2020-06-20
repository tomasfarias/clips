FROM golang:1.14 as builder

WORKDIR /src/
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o clips .

FROM alpine:latest

WORKDIR /root/
COPY --from=builder /src/clips .

CMD ./clips -t=$TOKEN -c=$CLIENT_ID -s=$CLIENT_SECRET
