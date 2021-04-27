FROM debian:bullseye-slim as builder

WORKDIR /build

RUN apt-get update && \
    apt-get install -y --no-install-recommends git golang ca-certificates libvips-dev gcc libc6-dev pkg-config make

COPY . .

RUN make converter && \
    for lib in $(ldd bin/converter | awk -F' ' '/=>/ {print $3}' | xargs); do [ -d libs$(dirname $lib) ] || mkdir -p libs$(dirname $lib) && cp -v $lib libs$lib ; done

FROM debian:bullseye-slim

WORKDIR /app

COPY --from=builder /build/bin/converter .
COPY --from=builder /build/libs/ /
RUN apt-get update && apt-get install -y --no-install-recommends libmagickcore-6.q16-6 && \
    rm -f /var/log/apt/* /var/log/{dpkg.log,faillog,lastlog}

ENTRYPOINT ["/app/converter"]
