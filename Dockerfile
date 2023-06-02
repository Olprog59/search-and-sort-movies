FROM alpine:latest as builder

RUN apk update && apk add go make

WORKDIR /app

ADD . .

RUN go mod download github.com/Olprog59/fsnotify
RUN go get github.com/Machiel/slugify

RUN make linux


FROM alpine:latest

VOLUME /be_sorted
VOLUME /movies
VOLUME /series

RUN mkdir /app

COPY --from=builder /app/bin/search-and-sort-movies-linux-amd64 /app

WORKDIR /app

ENV FORMAT_FILE "-, name, resolution, year"

CMD ["./search-and-sort-movies-linux-amd64", "-scan"]
