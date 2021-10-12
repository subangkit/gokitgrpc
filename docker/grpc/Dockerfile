FROM golang:alpine

RUN go env -w GO111MODULE=on
RUN apk add --update git
RUN rm -rf /var/cache/apk/*; rm -rf /tmp/*
#RUN go get -u -v github.com/golang/dep/cmd/dep


WORKDIR /go/src/gokitgrpc
ADD . .
RUN go mod download
RUN go get github.com/pilu/fresh
#COPY . /go/src/gokitgrpc
#RUN dep ensure

CMD ["fresh"]
EXPOSE 3000
