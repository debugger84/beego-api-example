FROM golang:1.8

ARG PUID=
RUN usermod -u ${PUID} www-data
ARG PROJECT_PATH=
WORKDIR ${PROJECT_PATH}
ARG PORT=8081
EXPOSE ${PORT}

# Install selected extensions and other stuff
#
RUN go get -u github.com/astaxie/beego
RUN go get -u github.com/beego/bee
RUN go get -u github.com/govend/govend
RUN govend

CMD ["bee", "run"]