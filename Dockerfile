FROM golang:1

LABEL org.opencontainers.image.authors="Christian Funkhouser <christian@funkhouse.rs>"
LABEL org.opencontainers.image.description="Echo server will help you \
understand your execution environment by echoing requests back to you."
LABEL org.opencontainers.image.licenses=MIT
LABEL org.opencontainers.image.source=https://github.com/cfunkhouser/echo
LABEL org.opencontainers.image.title="cfunkhouser/echoserver"
LABEL org.opencontainers.image.url="https://idontfixcomputers.com/echo"

ENV GOBIN /bin

COPY . /echo
WORKDIR /echo

RUN go install ./cmd/echoserver

EXPOSE 8080
CMD [ "/bin/echoserver" ]
