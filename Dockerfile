FROM alpine:latest
RUN apk update && apk upgrade
RUN apk add zsh vim curl git libstdc++ libc6-compat go
# RUN ln -s /lib64/ld-linux-x86-64.so.2 /lib/ld-linux-x86-64.so.2
RUN sh -c "$(curl -fsSL https://raw.github.com/ohmyzsh/ohmyzsh/master/tools/install.sh)"
WORKDIR /root
ENTRYPOINT ["sed", "-i", "-e", "'s/robbyrussell/ys/g'", ".zshrc"]

RUN mkdir -p /be_sorted
RUN mkdir -p /movies
RUN mkdir -p /series

ADD .vimrc .

WORKDIR /app
ADD bin/search-and-sort-movies-linux-amd64 .

RUN chmod +x search-and-sort-movies-linux-amd64

VOLUME /be_sorted
VOLUME /movies
VOLUME /series

ENTRYPOINT ["/app/search-and-sort-movies-linux-amd64", "-scan"]
