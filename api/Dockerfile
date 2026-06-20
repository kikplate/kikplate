FROM --platform=$BUILDPLATFORM golang:1.25-alpine AS builder

ARG TARGETOS
ARG TARGETARCH

RUN apk add --no-cache git

WORKDIR /src

COPY api/go.mod api/go.sum ./api/
COPY pkg/generator/go.mod ./pkg/generator/

WORKDIR /src/api
RUN go mod download

WORKDIR /src
COPY api ./api
COPY pkg ./pkg

WORKDIR /src/api
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -ldflags="-s -w" -o /app/kikplate-api .

FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata

RUN addgroup -S app && adduser -S app -G app

WORKDIR /app

COPY --from=builder /app/kikplate-api .

USER app

EXPOSE 3001

ENTRYPOINT ["./kikplate-api"]
CMD ["app:serve"]
