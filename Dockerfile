# Reference: https://docs.docker.com/engine/reference/builder/

ARG gover=1.14

FROM debian:latest
FROM golang:$gover

ARG user=omakoto
ARG group=user
ARG home=/home/$user
ARG shell=/bin/bash

ARG go_target=github.com/omakoto/zenlog/

ENV GOPATH=$home/go/
ARG copy_target=$GOPATH/src/$go_target

RUN apt-get update
RUN apt-get install -y git-core zsh vim less psmisc sudo procps libpcre++-dev man-db

RUN go get -v -t golang.org/x/lint/golint honnef.co/go/tools/cmd/... $go_target/zenlog/cmd/zenlog/
RUN go install $go_target/zenlog/cmd/zenlog/

WORKDIR $home
ENV HOME=$home

RUN groupadd -g 1000 $group && \
    useradd -r -u 1000 -g $group -s $shell $user

RUN mkdir -p $GOPATH
ENV PATH=$GOPATH/bin:$PATH

RUN mkdir -p $copy_target
COPY --chown=1000:1000 ./ $copy_target

ENV SHELL=$shell

RUN echo "PATH=$PATH" >> .profile ;\
    echo "if [ -n \"\$BASH_VERSION\" -a -f .bashrc ] ; then source .bashrc ; fi" >> .profile ;\
    echo "if [ -n \"\$ZSH_VERSION\" -a -f .zshrc ] ; then source .zshrc ; fi" >> .profile

RUN chown -R $user:$group $home

USER $user


#ENTRYPOINT $SHELL -l
