FROM golang:1.26-alpine AS build

WORKDIR /src
COPY go.mod go.sum ./
COPY vendor ./vendor
COPY cmd ./cmd
COPY pkg ./pkg

ENV CGO_ENABLED=0 GOOS=linux
RUN go build -mod=vendor -trimpath -o /out/server ./cmd/server
RUN go build -mod=vendor -trimpath -o /out/worker ./cmd/worker

FROM alpine:3.22
RUN apk --no-cache add ca-certificates
COPY --from=build /out/server /server
COPY --from=build /out/worker /worker

USER 65532:65532
EXPOSE 8080
CMD ["/server"]
