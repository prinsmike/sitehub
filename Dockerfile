FROM debian:latest
MAINTAINER Michael Prinsloo

ENV SH_PORT 80
ENV SH_WORKDIR /var/sitehub

ADD sitehub /sitehub

RUN chmod a+x /sitehub
RUN mkdir /var/sitehub

CMD ["/sitehub"]
