# Étape de build
FROM golang:1.22rc2 AS builder

# Définis le répertoire de travail
WORKDIR /app

# Copie les fichiers de dépendances et télécharge les dépendances
COPY go.mod go.sum ./
RUN go mod download && go mod tidy && go mod verify

# Copie le code source dans l'image
COPY . .

# Compile l'application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .


# Étape de création de l'image finale
FROM scratch

# Copie l'exécutable depuis l'étape de build
COPY --from=builder /app/main .

# Copie les certificats CA
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

VOLUME /be_sorted
VOLUME /movies
VOLUME /series
VOLUME /all

ENV A_TRIER=/movies/be_sorted
ENV MOVIES=/movies/movies
ENV SERIES=/movies/series
ENV ALL=/movies
ENV REGEX_MOVIES='{name}-{resolution} ({year})'
ENV REGEX_SERIES='{name}-s{season}e{episode}-{resolution} ({year})'

COPY --from=builder /app/bin/search-and-sort-movies-linux-amd64 /app

CMD ["./main", "-scan"]
