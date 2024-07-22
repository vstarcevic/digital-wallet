FROM golang:1.22 as builder

RUN mkdir /usr/src/app

WORKDIR /usr/src/app/cmd/user-service

COPY /user-service /usr/src/app

RUN CGO_ENABLED=0 go build -o user-service /usr/src/app/cmd/user-service
 
CMD [ "/usr/src/app/cmd/user-service/user-service"]