# certs and shared libs
FROM alpine:3.6 AS alpine
RUN apk add -U --no-cache ca-certificates

# preparing sane docker image to build app
FROM centos:8 AS build
RUN dnf upgrade -y && dnf install -y golang git make epel-release && dnf install -y upx && dnf clean all

# building app
RUN mkdir -p /opt/gossha
ADD . /opt/gossha/
RUN cd /opt/gossha && make build_without_tests
RUN ls -l /opt/gossha/build

# generating ssh key
RUN ssh-keygen -b 2048 -t rsa -f /opt/gossha/build/id_rsa -q -N ""


# result with size ~ 10mb
FROM scratch
COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=alpine /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
COPY --from=build /opt/gossha/build/gossha /app/gossha

COPY --from=build /opt/gossha/build/id_rsa /data/id_rsa
COPY --from=build /opt/gossha/build/id_rsa.pub /data/id_rsa.pub
COPY --from=build /opt/gossha/homedir/docker.toml /data/.gossha.toml

RUN /app/gossha root gossha gossha
VOLUME /data/

# Listen on 22 port
ENV GOSSHA_PORT=22
ENV GOSSHA_HOMEDIR=/data
ENV GOSSHA_DRIVER=sqlite3
ENV GOSSHA_CONNECTIONSTRING=/data/gossha.db

EXPOSE 22

ENTRYPOINT ["/app/gossha"]

