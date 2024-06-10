FROM 1.22.4-alpine3.20

WORKDIR /app

COPY . .

RUN go mod download && go build -o main .

CMD ["/app/main"]