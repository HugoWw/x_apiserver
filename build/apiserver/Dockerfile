FROM debian:buster-slim

LABEL name="X-APIServer" \
      vendor="Gouhuan." \
      image="x-apiserver/x-apiserver" \
      git.rev=${GIT_REV}


RUN set -x && \
         echo "deb http://mirrors.aliyun.com/debian/ buster main non-free contrib" > /etc/apt/sources.list && \
         echo "deb-src http://mirrors.aliyun.com/debian/ buster main non-free contrib" >> /etc/apt/sources.list && \
         apt-get update && apt-get install -y --no-install-recommends \
         vim \
         curl \
         ca-certificates \
         netcat-openbsd \
         libltdl7 \
         telnet \
         net-tools \
         && \
         apt-get clean && \
         rm -rf /var/lib/apt/lists/* && \
         ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

COPY x_apiserver /app/x_apiserver

RUN chmod +x /app/x_apiserver

CMD /app/x_apiserver

