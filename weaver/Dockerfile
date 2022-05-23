# ==== Building
FROM golang:1.16 AS build
WORKDIR /go/src/salucro-weaver

COPY . .
RUN go build -v -o build/weaver .

# ==== Running
FROM arachnysdocker/athenapdf AS run

RUN \
  wget https://github.com/Yelp/dumb-init/releases/download/v1.0.0/dumb-init_1.0.0_amd64.deb \
  && dpkg -i dumb-init_*.deb \
  && rm dumb-init_*.deb \
  && mkdir -p /athenapdf-service/tmp/

RUN apt-get update -y \
  && apt-get -y --force-yes install xvfb \
  && rm -rf /var/lib/apt/lists/* /var/cache/apt/*

COPY --from=build /go/src/salucro-weaver/build/ ./
COPY --from=build /go/src/salucro-weaver/conf/ ./conf/

EXPOSE 8080

CMD [ "./weaver" ]
ENTRYPOINT ["./conf/entrypoint.sh"]
