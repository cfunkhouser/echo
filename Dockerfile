FROM golang:1.22-alpine AS build
ENV GOBIN=/build
COPY . /echo
RUN mkdir -p /build && \
  go -C /echo install ./cmd/echoserver

FROM scratch AS echo
COPY --from=build /build/echoserver /bin/echoserver
LABEL org.opencontainers.image.authors="Christian Funkhouser <christian@funkhouse.rs>"
LABEL org.opencontainers.image.description="Echo server will help you \
understand your execution environment by echoing requests back to you."
LABEL org.opencontainers.image.licenses=MIT
LABEL org.opencontainers.image.source=https://github.com/cfunkhouser/echo
LABEL org.opencontainers.image.title="cfunkhouser/echoserver"
LABEL org.opencontainers.image.url="https://idontfixcomputers.com/echo"
EXPOSE 8080
ENTRYPOINT [ "/bin/echoserver" ]
CMD [ "-address", ":8080"]
