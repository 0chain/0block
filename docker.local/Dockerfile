FROM golang:1.12.4-alpine3.9 as zblock_build

RUN apk add --update --no-cache build-base linux-headers git cmake bash perl grep

# Install Herumi's cryptography
RUN apk add gmp gmp-dev openssl-dev && \
    cd /tmp && \
    wget -O - https://github.com/herumi/mcl/archive/master.tar.gz | tar xz && \
    mv mcl* mcl && \
    make -C mcl -j $(nproc) lib/libmclbn256.so install && \
    cp mcl/lib/libmclbn256.so /usr/local/lib && \
    rm -R /tmp/mcl && \
    git clone https://github.com/herumi/bls /tmp/bls && \
    cd /tmp/bls && \
    git submodule init && \
    git submodule update && \
    make -j $(nproc) install && \
    cd - && \
    rm -R /tmp/bls

ENV SRC_DIR=/0block
ENV GO111MODULE=on

# Download the dependencies:
# Will be cached if we don't change mod/sum files
COPY ./code/go/0block/core/go.mod          ./code/go/0block/core/go.sum          $SRC_DIR/go/0block/core/
COPY ./code/go/0block/zblockcore/go.mod     ./code/go/0block/zblockcore/go.sum     $SRC_DIR/go/0block/zblockcore/

WORKDIR $SRC_DIR/go/0block/zblockcore/zblock
RUN go mod download

# Add the source code:
ADD ./code/go/0block $SRC_DIR/go/0block

RUN go build -v -tags "bn256 development" -ldflags "-X 0chain.net/core/build.BuildTag=$GIT_COMMIT"

# Copy the build artifact into a minimal runtime image:
FROM golang:1.12.4-alpine3.9
ENV APP_DIR=/0block
WORKDIR $APP_DIR
RUN apk add gmp gmp-dev openssl-dev
COPY --from=zblock_build  /usr/local/lib/libmcl*.so \
    /usr/local/lib/libbls*.so \
    /usr/local/lib/
COPY --from=zblock_build $APP_DIR/go/0block/zblockcore/zblock/zblock $APP_DIR/bin/zblock

