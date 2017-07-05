FROM golang

ARG app_env
ENV APP_ENV $app_env

COPY ./ /go/src/github.com/synnick/go-playground
WORKDIR /go/src/github.com/synnick/go-playground

RUN go get ./
RUN go build

CMD if [ ${APP_ENV} = production ]; \
	then \
	go-playground; \
	else \
	go get github.com/pilu/fresh && \
	fresh; \
	fi
