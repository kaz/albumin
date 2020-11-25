FROM alpine:edge

ENV GOPATH /root/go
ENV PATH $GOPATH/bin:$PATH

WORKDIR /albumin
RUN apk add go make libheif-dev

COPY . /albumin
RUN make build

FROM alpine:edge

RUN apk add --no-cache libheif
COPY --from=0 /albumin/albumin /usr/local/bin/albumin

WORKDIR /albumin
EXPOSE 20000
CMD ["albumin"]
