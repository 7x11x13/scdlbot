FROM golang:1.22 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /scdlbot

FROM python:3.12-slim as runtime

WORKDIR /

RUN apt-get update \
    && apt-get install -y --no-install-recommends \
    ffmpeg \
    && rm -fr /var/lib/apt/lists/*

RUN python3 -m pip install scdl==2.12.3

COPY --from=builder /scdlbot /scdlbot

RUN addgroup \
    --gid 1001 \
    appgroup \
    &&  adduser \
    --disabled-password \
    --gecos "" \
    --ingroup appgroup \
    --uid 1001 \
    appuser

USER appuser

ENTRYPOINT [ "/scdlbot" ]