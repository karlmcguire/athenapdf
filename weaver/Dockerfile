FROM golang:1.16-alpine
WORKDIR /go/src/github.com/lachee/athenapdf/weaver

RUN apk add --update git

# Go-Dep is deprecated
# RUN go get -u github.com/golang/dep/cmd/dep
# COPY Gopkg.lock Gopkg.toml ./
# RUN dep ensure --vendor-only -v

COPY . ./
RUN go build -v -o weaver .

# Dont need dumb-init
# ENV GIN_MODE release
# RUN \
#   wget https://github.com/Yelp/dumb-init/releases/download/v1.0.0/dumb-init_1.0.0_amd64.deb \
#   && dpkg -i dumb-init_*.deb \
#   && rm dumb-init_*.deb \
#   && mkdir -p /athenapdf-service/tmp/

COPY weaver /athenapdf-service/
WORKDIR /athenapdf-service/

ENV PATH /athenapdf-service/:$PATH

COPY conf/ /athenapdf-service/conf/

EXPOSE 8080

CMD ["weaver"]

ENTRYPOINT ["sh", "/athenapdf-service/conf/entrypoint.sh"]