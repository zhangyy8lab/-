FROM octahub.8lab.cn:5000/base/ubt22-go:1.21.5 as builder
ENV GO111MODULE on
ENV GOPROXY https://goproxy.cn,direct
ENV GOPRIVATE https://bitbucket.org
ENV PACKJECT octa_accountserver

ADD .gitconfig /root/.gitconfig
ADD .git-credentials /root/.git-credentials

ADD ./ /data/go/src/$PACKJECT
WORKDIR /data/go/src/$PACKJECT
COPY sources.list /etc/apt/sources.list
RUN set -x \
    && apt update \
    && apt-get install -y git \
    && go mod tidy \
    && go mod vendor \
    && go build -ldflags -s -trimpath -v -o /data/go/bin/accountServer /data/go/src/$PACKJECT/src/


FROM octahub.8lab.cn:5000/base/ubuntu:22.04
ENV PACKJECT accountServer
ENV TEMPLOGPATH /data/oscro-v2/logs

MAINTAINER $PACKJECT:v2.0

COPY sources.list /etc/apt/sources.list

RUN set -x \
    && useradd -u 1000 ubuntu -s /bin/bash \
    && apt update && apt install -y tzdata ntpdate \
    && ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && mkdir -p $TEMPLOGPATH \
    && chmod 755 $TEMPLOGPATH

USER ubuntu

COPY --chown=ubuntu entrypoint.sh /usr/bin/entrypoint.sh
COPY --from=builder --chown=ubuntu /data/go/bin/$PACKJECT /usr/local/server/$PACKJECT

WORKDIR /usr/local/server/
EXPOSE 9091

RUN chmod +x /usr/bin/entrypoint.sh

ENTRYPOINT [ "/usr/bin/entrypoint.sh" ]