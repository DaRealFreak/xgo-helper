FROM karalabe/xgo-latest

MAINTAINER Steffen Keuper <steffen.keuper@web.de>

COPY docker/docker-entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

ENTRYPOINT [ "/entrypoint.sh" ]