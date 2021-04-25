FROM golang:latest as builder
WORKDIR /app
ADD . .
RUN make linux


FROM ubuntu:latest

VOLUME /be_sorted
VOLUME /movies
VOLUME /series

COPY --from=builder /app/bin/search-and-sort-movies-linux-amd64 .

CMD ["./search-and-sort-movies-linux-amd64", "-scan"]