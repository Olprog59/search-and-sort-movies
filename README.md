# Media Organizer

**Media Organizer** est une application robuste conçue pour l'organisation et le renommage automatisés de collections de fichiers multimédias, tels que les films et les séries télévisées.

## Fonctionnalités Principales
- **Tri Automatique** : Classe les fichiers multimédias dans des répertoires appropriés basés sur leur type.
- **Renommage Intuitif** : Applique des conventions de nommage personnalisables aux fichiers.
- **Support plateforme** : Fonctionne sur Linux (`amd64`, `arm64`).
- **Intégration FFmpeg** : Utilise `ffprobe` pour l'extraction de la durée des films et des séries. Cela permet de voir s'il y a une incohérence entre la durée du fichier et la durée du média. Par exemple, si un fichier fait moins de **60 minutes** et avec mon algo, il détecte un film alors il ne se déplace pas. Dans ce cas, vous avez la possibilité d'aller sur l'interface **Web** afin de corriger le nom d'origine en y ajoutant par exemple `Saison 1 Episode 1` ou `s01e1095` etc.

## Technologies
- Langage : Go (Golang).
- Conteneurisation : Docker et Docker Compose.
- Traitement Multimédia : FFmpeg.

## Architecture Docker
Le déploiement utilise un `Dockerfile` multi-stage pour optimiser la taille de l'image et l'efficacité du build.

```dockerfile

FROM golang:1.21.6 AS builder
ARG GOOS=linux
ARG GOARCH=amd64
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download && go mod tidy && go mod verify
COPY . .
RUN CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} go build -a -installsuffix cgo -o main .
```

Final Stage
```dockerfile
FROM alpine:latest
RUN apk add --no-cache ffmpeg
COPY --from=builder /app/main /app/main
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
VOLUME ["/be_sorted", "/movies", "/series"]
EXPOSE 8080/tcp
CMD ["/app/main", "-scan"]
```

Déploiement avec Docker Compose
Créez un fichier docker-compose.yml :

```yaml
version: "3.8"
services:
  media-organizer:
    image: olprog/media-organizer:latest
    ports:
      - 1574:8080
    volumes:
      - /mnt/medias/be_sorted:/be_sorted
      - /mnt/medias/medias:/medias
      - /mnt/medias/series:/series
    environment:
      REGEX_MOVIE: "{name}-{resolution} ({year})"
      REGEX_SERIE: "{name}-s{season}e{episode}-{resolution} ({year})"
      UID: "0"
      GID: "0"
      CHMOD: "0755"
    restart: always
```

Lancez l'application :

```bash
docker-compose up
```

## Exemple d'utilisation avec Docker Run

Pour exécuter directement avec Docker :
    
```bash
docker run -d \
  --name media-organizer \
  -p 1574:8080 \
  -v /mnt/medias/be_sorted:/be_sorted \
  -v /mnt/medias/medias:/medias \
  -v /mnt/medias/series:/series \
  -e REGEX_MOVIE="{name}-{resolution} ({year})" \
  -e REGEX_SERIE="{name}-s{season}e{episode}-{resolution} ({year})" \
  -e UID="0" \
  -e GID="0" \
  -e CHMOD="0755" \
  --restart always \
  olprog/media-organizer:latest
```

| Variable      | Description                                                                                                                                                                                                                                                                                         | Valeur par Défaut                                    | Exemples de Valeurs / Formats Acceptés                                                                                                                                                                |
|---------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|------------------------------------------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `REGEX_MOVIE` | Expression régulière pour le formatage des noms de films.                                                                                                                                                                                                                                           | `"{name}-{resolution} ({year})"`                     | `"{name}-{resolution} ({year})"`, `"{year}-{name}"` <br/>- nom du fichier : `{name}`<br/>- résolution: `{resolution}`<br/>- année de sortie(si présente): `{year}`<br/>- langage: {language}          |
| `REGEX_SERIE` | Expression régulière pour le formatage des noms de séries.                                                                                                                                                                                                                                          | `"{name}-s{season}e{episode}-{resolution} ({year})"` | `"{name}-s{season}e{episode}-{resolution} ({year})"`, `"{name} S{season}E{episode}"`<br/>- même que `REGEX_MOVIE` avec en plus:<br/>- saison(numérique): {season}<br/>- épisode(numérique): {episode} |
| `UID`         | Identifiant utilisateur pour la gestion des permissions.                                                                                                                                                                                                                                            | `"0"` (root)                                         | `"0"`, `"1000"` (ou autre UID utilisateur)                                                                                                                                                            |
| `GID`         | Identifiant de groupe pour la gestion des permissions.                                                                                                                                                                                                                                              | `"0"` (root)                                         | `"0"`, `"1000"` (ou autre GID groupe)                                                                                                                                                                 |
| `CHMOD`       | Permissions par défaut pour les fichiers et dossiers créés.                                                                                                                                                                                                                                         | `"0755"`                                             | `"0755"`, `"0777"` (pour un accès plus ouvert)                                                                                                                                                        |

## Notes
> **Note** : Les variables `UID` et `GID` sont utilisées pour définir les permissions des fichiers et dossiers créés par l'application. Elles sont utiles pour les systèmes de fichiers montés avec des permissions spécifiques.

> **Note** : Les variables `REGEX_MOVIE` et `REGEX_SERIE` sont utilisées pour définir les conventions de nommage des fichiers multimédias. Elles sont utiles pour personnaliser le format des noms de fichiers.
> <br>
> - `{language}` : Langage du film ou de la série détecté dans le titre. (actuellement configuré: `french`, `multi`, `vostfr`, `subfrench`, `vo`)
 
> **Note** : La détection des films se fait avec cette expression régulière : `(?mi)-(french|vf|dvdrip|multi|vostfr|subfrench|dvd-r|bluray|bdrip|brrip|cam|ts|tc|vcd|md|ld|r[0-9]|xvid|divx|scr|dvdscr|repack|hdlight|720p|480p|1080p|2160p|uhd|4k|1920x1080)` 
> <br><br>
> **Note** : La détection des séries se fait avec cette expression régulière : `(?mi)((s\d{1,2})(?:\W+)?(e?\d{1,4}))|e\d{1,4}|(episode-(\d{2,4})-?)|((\d{1,2})-(\d{2,4}))|((saison|season)-(\d{1,2})-episode-(\d{1,4}))`
> <br><br>
> Si vous souhaitez tester ces expressions régulières, vous pouvez utiliser [regex101.com](https://regex101.com/).


## A venir
- [ ] Ajout de la traduction.
- [ ] Ajout de la suppression des dossiers en plus des fichiers autres comme maintenant.
- [ ] Variables d'environnement pour les dossiers principaux.
- [ ] Ajout de possibilités de connecter un discord webhooks pour les notifications en cas de problèmes.

