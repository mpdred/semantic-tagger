# ---
FROM golang AS dependencies
WORKDIR /app
RUN apt update > /dev/null \
  && apt install -y make > /dev/null

COPY go.* ./
COPY Makefile ./
RUN make dependencies




# ---
FROM dependencies AS base
COPY . .



# ---
FROM base AS test
RUN make test




# ---
FROM base AS build
RUN make build-linux




# ---
FROM debian:stable-slim AS final

WORKDIR /app
RUN apt update && apt install -y bash git openssh-client ca-certificates


COPY --from=build /app/bin/semtag-linux-amd64 /usr/local/bin/semtag-linux-amd64
COPY --from=build /app/bin/semtag-linux-386 /usr/local/bin/semtag-linux-386
COPY --from=test /app/README.md ./
RUN chmod +x /usr/local/bin/semtag*
RUN ln -s /usr/local/bin/semtag-linux-386 /usr/local/bin/semtag

LABEL description="Set up versioning by using git tags or files"
LABEL documentation="https://github.com/mpdred/semantic-tagger/blob/master/README.md"
LABEL maintainer="https://github.com/mpdred"
LABEL name="semantic-tagger"
LABEL source="https://github.com/mpdred/semantic-tagger"
LABEL url="https://hub.docker.com/r/mpdred/semantic-tagger"

ENTRYPOINT "/bin/bash"
