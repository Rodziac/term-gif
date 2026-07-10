# ---- build stage ----
FROM golang:1.23-bookworm AS build

WORKDIR /src

# Cache modules first.
COPY go.mod go.sum ./
RUN go mod download

# Build a fully static binary so it runs on a minimal runtime image.
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -o /out/term-gif .

# ---- runtime stage ----
FROM alpine:3.20

WORKDIR /app

# CA certificates are required for the outbound HTTPS calls to the Giphy API
# and the media CDN; without them every fetch fails and you only ever get the
# "oops" fallback gif.
RUN apk add --no-cache ca-certificates

COPY --from=build /out/term-gif /app/term-gif
# Static assets: the browser landing page and the oops.gif fallback, both
# loaded via relative paths at runtime.
COPY static/ /app/static/
# Config defaults; real secrets (APIKEY) are supplied as env / config vars and
# override these at runtime.
COPY env.sample /app/env.sample

# Documentation only. Heroku/Render/Fly inject the real port via $PORT, which
# the app reads through viper (see BindEnv in main.go).
EXPOSE 9000

CMD ["/app/term-gif"]
