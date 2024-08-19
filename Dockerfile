FROM golang:1.22.3-alpine

WORKDIR /ascii-art-web

COPY . . 

RUN apk update && apk add bash 

EXPOSE 8080

LABEL authors="Mihit Mohammed - Abdelfattah Elidrissi"
LABEL version="1.0"
LABEL description="This is a docker image for ascii art web project"

CMD ["go", "run", "."] 