#!/bin/bash

# Test de nos fichiers test
ALL_TEST=$(go test ./... -v)
COUNT_FAIL=$(echo "$ALL_TEST" | grep -o 'FAIL' | wc -l | tr -d '[:space:]')
if [ $# = 0 ]
then
    echo "N'oublie pas de mettre un premier paramètre. Ex: "
    printf "\t-M (Major) : 0.6.3.2 -> 1.0.0.0\n"
    printf "\t-m (minor) : 0.6.3.2 -> 0.7.0.0\n"
    printf "\t-p (patch) : 0.6.3.2 -> 0.6.4.0\n"
    printf "\t-s (special when beta) : 0.6.3.2 -> 0.6.3.3\n"
    printf "\t-n (debug rapid) pas de changement de version\n"
    exit;
fi

function increment() {

  while getopts ":Mmpsn" Option
  do
    case $Option in
      M ) major=true;;
      m ) minor=true;;
      p ) patch=true;;
      s ) special=true;;
      n ) none=true;;
    esac
  done

  shift $(($OPTIND - 1))

  version=$1

  # Build array from version string.

  a=( ${version//./ } )

  # If version string is missing or has the wrong number of members, show usage message.

  if [ ${#a[@]} -ne 5 ]
  then
    echo "usage: $(basename $0) [-Mmpsn] Major.minor.patch.special.none"
    exit 1
  fi

  # Increment version numbers as requested.

  if [ ! -z $major ]
  then
    ((a[0]++))
    a[1]=0
    a[2]=0
    a[3]=0
  fi

  if [ ! -z $minor ]
  then
    ((a[1]++))
    a[2]=0
    a[3]=0
  fi

  if [ ! -z $patch ]
  then
    ((a[2]++))
    a[3]=0
  fi

  if [ ! -z $special ]
  then
    ((a[3]++))
  fi
  if [ ! -z $none ]
  then
    make linux
    echo "Je pense que tout est correct"
    exit
  fi

  NEW_VERSION="${a[0]}.${a[1]}.${a[2]}.${a[3]}"
}

if [ "$COUNT_FAIL" = 0 ]
then
    GIT_STATUS=$(git status --porcelain | wc -l | tr -d '[:space:]')
    if [ "$GIT_STATUS" = 0 ]
    then
      BUILD_CLEAN=yes
    else
      BUILD_CLEAN=no
    fi

    if [ $BUILD_CLEAN = yes ]
    then
        echo "Aucun commit à faire. Tu es à jour!"
        # Increment version
        VERSION=$(cat 'VERSION')

        increment "$1" "$VERSION" NEW_VERSION

#        NEW_VERSION="${VERSION%.*}.$((${VERSION##*.}+1))"

        echo "Attention, la version a changé en ${NEW_VERSION}. Souhaites-tu continuer à commit puis make le project?"
        select yn in "Yes" "No"; do
            case $yn in
                Yes ) printf "%s" "$NEW_VERSION" > 'VERSION'; git add VERSION && git commit -m "VERSION -> $NEW_VERSION" ; break;;
                No ) exit;;
            esac
        done

        git add VERSION && git commit -m "VERSION -> $NEW_VERSION"

        GIT_STATUS=$(git status --porcelain | wc -l | tr -d '[:space:]')
        if [ "$GIT_STATUS" = 0 ]
        then
          BUILD_CLEAN=yes
        else
          BUILD_CLEAN=no
        fi
        if [ $BUILD_CLEAN = yes ]
        then
            # Make all
            make linux
            echo "Je pense que tout est correct"
#            scp bin/* root@vps.olprog.fr:/media/app/static
        else
            echo "Problème de commit juste après un changement de version. A vérifier"
        fi
    else
        echo "Tu as oublié de commit et faire le push"
    fi


else
    echo "Erreur de test - Vérifie avec : go test ./... -v (--cover)"
fi

