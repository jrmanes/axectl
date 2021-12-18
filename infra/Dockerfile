FROM ubuntu:21.10
WORKDIR /app
RUN apt update \
    && apt upgrade -y \
    && apt list --upgradable \
    && apt install -y \
    sudo \
    tzdata

ENV TZ="America/New_York"
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

COPY ./bin/ .
