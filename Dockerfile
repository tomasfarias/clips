FROM go:1.14 as builder

WORKDIR /src/
COPY . .

RUN go build .

FROM alpine:latest

COPY --from=builder /src/clips /usr/local/bin/clips

CMD clips -t $TOKEN -c $CLIENT_ID -s $CLIENT_SECRET

