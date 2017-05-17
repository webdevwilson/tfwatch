FROM hashicorp/terraform:full
MAINTAINER "Kerry Wilson <kwilson@goodercode.com>"

RUN apk add --update make nodejs

ENV APP_PATH=$GOPATH/src/github.com/webdevwilson/terraform-ci

WORKDIR $APP_PATH

# Copy sources
ADD . $APP_PATH

RUN make build install

VOLUME /var/lib/terraform-ci

ENTRYPOINT [ "terraform-ci" ]