FROM golang:1.16-buster
WORKDIR /go/src/app
COPY ./banktransfer .

RUN apt update && apt install -y protobuf-compiler
RUN go get google.golang.org/protobuf/cmd/protoc-gen-go google.golang.org/grpc/cmd/protoc-gen-go-grpc

RUN go mod download
RUN go generate ./...
RUN go install

RUN wget https://raw.githubusercontent.com/vishnubob/wait-for-it/81b1373f17855a4dc21156cfe1694c31d7d1792e/wait-for-it.sh
RUN chmod +x ./wait-for-it.sh ./docker-entrypoint.sh
ENTRYPOINT ["./docker-entrypoint.sh"]

CMD ["banktransfer"]

EXPOSE 9111