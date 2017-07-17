FROM    alpine:3.6

RUN     apk -U add \
            python \
            py-pip \
            go \
            bash \
            git \
            gcc \
            musl-dev

ENV     GOPATH=/go
RUN     git config --global http.https://gopkg.in.followRedirects true
RUN     go get github.com/dnephin/filewatcher && \
        cp /go/bin/filewatcher /usr/bin/ && \
        rm -rf /go/src/* /go/pkg/* /go/bin/*

RUN     pip install sphinx==1.4.5

WORKDIR /go/src/github.com/dnephin/dobi
ENV     PS1="# "
