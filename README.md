# Media Organizer

**Media Organizer** est une application robuste conçue pour l'organisation et le renommage automatisés de collections de fichiers multimédias, tels que les films et les séries télévisées.

## Fonctionnalités Principales
- **Tri Automatique** : Classe les fichiers multimédias dans des répertoires appropriés basés sur leur type.
- **Renommage Intuitif** : Applique des conventions de nommage personnalisables aux fichiers.
- **Support plateforme** : Fonctionne sur Linux (`amd64`, `arm64`).
- **Intégration FFmpeg** : Utilise `ffprobe` pour l'extraction de la durée des films et des séries. Cela permet de voir s'il y a une incohérence entre la durée du fichier et la durée du média. Par exemple, si un fichier fait moins de **60 minutes** et avec mon algo, il détecte un film alors il ne se déplace pas. Dans ce cas, vous avez la possibilité d'aller sur l'interface **Web** afin de corriger le nom d'origine en y ajoutant par exemple `Saison 1 Episode 1` ou `s01e1095` etc.

## Interface Web pour la Gestion des Fichiers

- **Gestion Manuelle** : Pour les fichiers qui n'ont pas été automatiquement triés ou dont le classement semble
  inapproprié, l'application offre une **interface web** conviviale. Cela vous permet de modifier manuellement les
  informations des fichiers, comme en ajoutant `Saison 1 Episode 1` ou `s01e01`, pour garantir que chaque média soit
  correctement identifié et classé. Cette interface sert de complément parfait au tri automatique, offrant une
  flexibilité totale dans la gestion de votre bibliothèque multimédia.

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
      # Par défaut, le conteneur monte le dossier /medias/be_sorted, /medias/movies et /medias/series
      # Vous pouvez les changer en fonction de vos besoins avec les variables d'environnement.
      - /mnt/medias:/medias
    environment:
      # Par défaut, BE_SORTED est /medias/be_sorted, MOVIES est /medias/movies et SERIES est /medias/series
      # BE_SORTED : "/medias/a_trier"
      # MOVIES : "/medias/films"
      # SERIES : "/medias/series"
      REGEX_MOVIE: "{name}-{resolution} ({year})"
      REGEX_SERIE: "{name}-s{season}e{episode}-{resolution} ({year})"
      UID: "0"
      GID: "0"
      CHMOD: "0755"
    restart: always
```

Autre possibilité docker-compose.yml :

```yaml
version: "3.8"
services:
  media-organizer:
    image: olprog/media-organizer:latest
    ports:
      - 1574:8080
    volumes:
      # Par défaut, le conteneur monte le dossier /medias/be_sorted, /medias/movies et /medias/series
      - /mnt/medias/a_trier:/medias/be_sorted
      - /mnt/medias/films:/medias/movies
      - /mnt/medias/séries:/medias/series
    environment:
      REGEX_MOVIE: "{name}-{resolution} ({year})"
      REGEX_SERIE: "{name}-s{season}e{episode}-{resolution} ({year})"
      UID: "0"
      GID: "0"
      CHMOD: "0755"
    restart: always
```

L'existence de multiples options de configuration vise à s'aligner au mieux sur vos exigences spécifiques. Vous avez la
liberté soit de monter des dossiers particuliers selon votre choix, soit d'opter pour le montage d'un unique dossier,
lequel est ensuite accessible via des variables d'environnement.

La première option offre un avantage notable : elle permet un accès direct à l'hôte, favorisant ainsi une performance
accrue pour les opérations de renommage et de déplacement de fichiers. À l'opposé, la seconde option, tout en vous
donnant la flexibilité de choisir et de partager des dossiers spécifiques via des variables d'environnement, implique
une gestion de multiples volumes, ce qui, pour les transferts de fichiers entre différents dossiers, peut se traduire
par une opération similaire à une copie de fichiers d'un emplacement à un autre.

Comparaison de performance entre les deux configurations (basée sur mes tests sur mon serveur, les résultats pouvant
varier selon votre propre infrastructure) :

> Première configuration : Transfert de 22Go en 10ms

> Deuxième configuration : Transfert de 22Go en 2m30s


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
  -v /mnt/medias/be_sorted:/medias/be_sorted \
  -v /mnt/medias/movies:/medias/movies \
  -v /mnt/medias/series:/medias/series \
  -e REGEX_MOVIE="{name}-{resolution} ({year})" \
  -e REGEX_SERIE="{name}-s{season}e{episode}-{resolution} ({year})" \
  -e UID="0" \
  -e GID="0" \
  -e CHMOD="0755" \
  --restart always \
  olprog/media-organizer:latest
```

```bash
docker run -d \
  --name media-organizer \
  -p 1574:8080 \
  -v /mnt/medias:/medias \
  -e BE_SORTED="/medias/a_trier" \
  -e MOVIES="/medias/films" \
  -e SERIES="/medias/series" \
  -e REGEX_MOVIE="{name}-{resolution} ({year})" \
  -e REGEX_SERIE="{name}-s{season}e{episode}-{resolution} ({year})" \
  -e UID="0" \
  -e GID="0" \
  -e CHMOD="0755" \
  --restart always \
  olprog/media-organizer:latest
```

| Variable      | Description                                                   | Valeur par Défaut                                    | Exemples de Valeurs / Formats Acceptés                                                                                                                                                                |
|---------------|---------------------------------------------------------------|------------------------------------------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `BE_SORTED`   | Dossier de destination pour les fichiers multimédias à trier. | `"/medias/be_sorted"`                                | `"/mnt/multimedia/a_trier"`                                                                                                                                                                           |
| `MOVIES`      | Dossier de destination pour les films triés.                  | `"/medias/movies"`                                   | `"/films"`                                                                                                                                                                                            |
| `SERIES`      | Dossier de destination pour les séries triées.                | `"/medias/series"`                                   | `"/medias/series"`                                                                                                                                                                                    |
| `REGEX_MOVIE` | Expression régulière pour le formatage des noms de films.     | `"{name}-{resolution} ({year})"`                     | `"{name}-{resolution} ({year})"`, `"{year}-{name}"` <br/>- nom du fichier : `{name}`<br/>- résolution: `{resolution}`<br/>- année de sortie(si présente): `{year}`<br/>- langage: {language}          |
| `REGEX_SERIE` | Expression régulière pour le formatage des noms de séries.    | `"{name}-s{season}e{episode}-{resolution} ({year})"` | `"{name}-s{season}e{episode}-{resolution} ({year})"`, `"{name} S{season}E{episode}"`<br/>- même que `REGEX_MOVIE` avec en plus:<br/>- saison(numérique): {season}<br/>- épisode(numérique): {episode} |
| `UID`         | Identifiant utilisateur pour la gestion des permissions.      | `"0"` (root)                                         | `"0"`, `"1000"` (ou autre UID utilisateur)                                                                                                                                                            |
| `GID`         | Identifiant de groupe pour la gestion des permissions.        | `"0"` (root)                                         | `"0"`, `"1000"` (ou autre GID groupe)                                                                                                                                                                 |
| `CHMOD`       | Permissions par défaut pour les fichiers et dossiers créés.   | `"0755"`                                             | `"0755"`, `"0777"` (pour un accès plus ouvert)                                                                                                                                                        |

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
- [ ] Ajout de possibilités de connecter un discord webhooks pour les notifications en cas de problèmes.

