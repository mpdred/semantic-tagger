# ---
FROM golang AS src
WORKDIR /app
RUN apt update > /dev/null \
  && apt install -y make > /dev/null

COPY . .



# ---
FROM src AS test
RUN make test




# ---
FROM src AS build
RUN make build




# ---
FROM alpine AS final
WORKDIR /app
RUN apk add --no-cache \
      bash git openssh-client ca-certificates


COPY --from=build /tmp/semtag /usr/local/bin
RUN chmod +x /usr/local/bin/semtag
ENTRYPOINT "/bin/bash"
