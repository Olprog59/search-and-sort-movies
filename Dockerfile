# Étape de build
FROM golang:latest AS builder

# Définis le répertoire de travail
WORKDIR /app

# Copie les fichiers de dépendances et télécharge les dépendances
COPY go.mod go.sum ./
RUN go mod download && go mod tidy && go mod verify

# Copie le code source dans l'image
COPY . .

# Compile l'application
RUN GOOS=linux go build -o main .

# Étape de création de l'image finale
FROM alpine:latest

# Installe ffprobe
RUN apk add --no-cache ffmpeg

# Copie l'exécutable depuis l'étape de build
COPY --from=builder /app/main /app/main

# Copie les certificats CA (si nécessaire)
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

VOLUME ["/be_sorted", "/movies", "/series", "/all"]

ENV A_TRIER=/mnt/medias/be_sorted \
    MOVIES=/mnt/medias/movies \
    SERIES=/mnt/medias/series \
    ALL="" \
    REGEX_MOVIES='{name}-{resolution} ({year})' \
    REGEX_SERIES='{name}-s{season}e{episode}-{resolution} ({year})' \
    UID=0 \
    GID=0 \
    CHMOD=0755

EXPOSE 8080/tcp

CMD ["/app/main", "-scan"]
