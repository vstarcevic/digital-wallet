FROM golang:1.22 as builder

RUN mkdir /usr/src/app

WORKDIR /usr/src/app/cmd/transaction-service

COPY /transaction-service /usr/src/app

RUN CGO_ENABLED=0 go build -o transaction-service /usr/src/app/cmd/transaction-service
 
CMD [ "/usr/src/app/cmd/transaction-service/transaction-service"]