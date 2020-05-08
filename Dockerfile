FROM golang:alpine

COPY app app

CMD ["./app"]