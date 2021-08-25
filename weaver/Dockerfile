# ==== Building
FROM golang:1.16-alpine AS build
WORKDIR /go/src/github.com/lachee/athenapdf/weaver

RUN apk add --update git

# Go-Dep is deprecated
# RUN go get -u github.com/golang/dep/cmd/dep
# COPY Gopkg.lock Gopkg.toml ./
# RUN dep ensure --vendor-only -v

COPY . ./
RUN go build -v -o weaver .

  
# ==== Running
FROM arachnysdocker/athenapdf AS run

# Setup dumb-init
ENV GIN_MODE release
RUN \
  wget https://github.com/Yelp/dumb-init/releases/download/v1.0.0/dumb-init_1.0.0_amd64.deb \
  && dpkg -i dumb-init_*.deb \
  && rm dumb-init_*.deb \
  && mkdir -p /athenapdf-service/tmp/

# Copy the items across
COPY --from=build /go/src/github.com/lachee/athenapdf/weaver/weaver /athenapdf-service/weaver
COPY --from=build /go/src/github.com/lachee/athenapdf/weaver/conf/ /athenapdf-service/conf/

WORKDIR /athenapdf-service/

# Setup path
ENV PATH /athenapdf-service/:$PATH

# Expose
EXPOSE 8080

CMD ["dumb-init", "weaver"]

ENTRYPOINT ["/athenapdf-service/conf/entrypoint.sh"]