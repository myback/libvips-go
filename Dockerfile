FROM debian:bookworm-slim as builder

WORKDIR /build

ENV DEBIAN_FRONTEND noninteractive
ENV PATH /usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/usr/local/go/bin

RUN apt update \
    && apt install gnupg2 tzdata --no-install-recommends -y \
    && dpkg-reconfigure --frontend noninteractive tzdata \
    && apt-key adv --keyserver keyserver.ubuntu.com --recv-keys 1E9377A2BA9EF27F \
    && echo 'deb http://ppa.launchpad.net/ubuntu-toolchain-r/test/ubuntu focal main' > /etc/apt/sources.list.d/ubuntu-toolchain-r-ubuntu-test-focal.list \
    && apt update \
    && apt-get install -y --no-install-recommends curl git ca-certificates libvips-dev libc6-dev pkg-config make gcc-11 g++-11 \
    && ln -s /usr/bin/gcc-11 /usr/bin/gcc \
    && ln -s /usr/bin/g++-11 /usr/bin/g++ \
    && curl -sSfLO https://go.dev/dl/go1.18.1.linux-amd64.tar.gz \
    && tar -C /usr/local -xzf go1.18.1.linux-amd64.tar.gz

COPY . .

RUN make converter && \
    for lib in $(ldd bin/converter | awk -F' ' '/=>/ {print $3}' | xargs); do [ -d libs$(dirname $lib) ] || mkdir -p libs$(dirname $lib) && cp -v $lib libs$lib ; done

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /build/bin/converter .
COPY --from=builder /build/libs/ /
RUN apt-get update && apt-get install -y --no-install-recommends libmagickcore-6.q16-6 && \
    rm -f /var/log/apt/* /var/log/{dpkg.log,faillog,lastlog}

ENTRYPOINT ["/app/converter"]
