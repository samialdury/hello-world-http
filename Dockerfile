ARG GO_VERSION=1.22
ARG WORK_DIR="/app"

ARG VERSION_SHA="unknown-sha"
ARG PORT=8080

FROM golang:${GO_VERSION}-alpine as builder

ARG WORK_DIR
ARG VERSION_SHA

WORKDIR ${WORK_DIR}

COPY go.mod ./

COPY main.go ./

RUN CGO_ENABLED=0 GOOS=linux \ 
    go build \ 
    -o bin/main \
    -ldflags "-X main.VersionSHA=${VERSION_SHA}" \
    main.go

FROM gcr.io/distroless/static-debian12 as runtime

ARG WORK_DIR
ARG PORT

WORKDIR ${WORK_DIR}

COPY --from=builder --chown=nonroot:nonroot ${WORK_DIR}/bin/main ./bin/main

ENV TZ=UTC
ENV PORT=${PORT}

EXPOSE ${PORT}

USER nonroot:nonroot

ENTRYPOINT ["./bin/main"]
