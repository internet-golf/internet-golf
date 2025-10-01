FROM golang:1.25

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

# is there a better way to copy in all these files/folders?
COPY pkg /app/pkg
COPY cmd /app/cmd
COPY client-sdk /app/client-sdk
COPY client-cmd /app/client-cmd
COPY magefile.go magefile.go

RUN go tool mage build

EXPOSE 80 443 8888

CMD ["./golf-server"]

# for persistence, create bind mount or volume for "/root/.internetgolf"
