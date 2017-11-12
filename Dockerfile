FROM scratch

MAINTAINER Guilherme Silveira <xguiga@gmail.com>

COPY simplesurance-api /

EXPOSE 8080

ENTRYPOINT ["/simplesurance-api"]
