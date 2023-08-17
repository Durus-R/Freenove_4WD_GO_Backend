FROM alpine AS WS281X_BUILDER

RUN apk add --no-cache alpine-sdk cmake git linux-headers

WORKDIR /foundry

RUN git clone https://github.com/jgarff/rpi_ws281x.git \
  && cd rpi_ws281x \
  && mkdir build \
  && cd build \
  && cmake -D BUILD_SHARED=OFF -D BUILD_TEST=OFF .. \
  && cmake --build . \
  && make install

# Stage 1 : Build a go image with the rpi_ws281x C library

FROM golang:1.20-alpine

RUN apk add --no-cache protobuf bash libprotobuf protobuf-dev gcc alpine-sdk

ENV PATH="/go/bin:${PATH}:"


RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.30.0
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0

COPY --from=WS281X_BUILDER /usr/local/lib/libws2811.a /usr/local/lib/
COPY --from=WS281X_BUILDER /usr/local/include/ws2811 /usr/local/include/ws2811
