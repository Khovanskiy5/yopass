FROM golang:1.24-bookworm AS app
WORKDIR /yopass
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o yopass ./cmd/yopass && \
    CGO_ENABLED=0 go build -o yopass-server ./cmd/yopass-server

FROM oven/bun:1 AS website
WORKDIR /website
COPY website/package.json website/bun.lock ./
RUN bun install --frozen-lockfile
COPY website .
RUN bun run build

FROM gcr.io/distroless/static
COPY --from=app /yopass/yopass /yopass/yopass-server /
COPY --from=website /website/dist /public
USER 1000
ENTRYPOINT ["/yopass-server"]
