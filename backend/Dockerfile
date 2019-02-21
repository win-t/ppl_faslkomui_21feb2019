FROM golang:alpine
RUN apk -U add git
WORKDIR /app
ADD go.mod .
RUN go mod download
ADD . .
RUN CGO_ENABLED=0 go build -o output
ENTRYPOINT [ "/app/output" ]
