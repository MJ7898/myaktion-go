FROM golang:1.16-buster
WORKDIR /go/src/app
COPY . .
RUN go mod download
RUN go install

# RUN wget root for wait-for-it.sh Shell script as an dependencie
RUN chmod +x ./wait-for-it.sh ./docker-entrypoint.sh
ENTRYPOINT ["./docker-entrypoint.sh"]
CMD ["myaktion"]
EXPOSE 8000