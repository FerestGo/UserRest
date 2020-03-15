FROM golang:1.12.7-alpine AS builder

WORKDIR $GOPATH/src/user

ARG BITBUCKET_SSH_KEY

RUN apk add git build-base gettext curl openssh-client \
 && addgroup -g 1000 www-data \
 && adduser -u 1000 -G www-data -D www-data \
 && mkdir /root/.ssh \
 && echo "${BITBUCKET_SSH_KEY}" > /root/.ssh/id_rsa \
 && chmod 700 /root/.ssh \
 && chmod 600 /root/.ssh/id_rsa \
 && git config --global url."git@bitbucket.org:".insteadOf "https://bitbucket.org/" \
 && ssh-keyscan bitbucket.org > /root/.ssh/known_hosts

COPY --chown=1000:1000 . $GOPATH/src/user/

RUN chown www-data:www-data $GOPATH/src/user

RUN go get -d -v ./...

RUN go get github.com/stretchr/testify/assert
RUN go get github.com/golangci/golangci-lint/cmd/golangci-lint

RUN go test

RUN go install -v ./...

RUN golangci-lint run

FROM alpine AS production

USER root

COPY --from=builder /etc/passwd /etc/passwd

COPY --from=builder /etc/group /etc/group

COPY --from=builder /go/bin/uesr /go/bin/user

RUN apk update \
 && apk add curl

USER www-data

CMD ["/go/bin/user"]