FROM golang:1.22
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o /main
ENV TODO_PORT=7540 \
    TODO_DBFILE=/app/data/scheduler.db \
    TODO_PASSWORD=your_secret_password
EXPOSE 7540
CMD ["/main"]