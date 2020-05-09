#!/bin/bash

# Test de nos fichiers test
ALL_TEST=$(go test ./... -v)
COUNT_FAIL=$(echo $ALL_TEST | grep -o 'FAIL' | wc -l | tr -d '[:space:]')
if [ $COUNT_FAIL = 0 ]
then
    GIT_STATUS=$(git status --porcelain | wc -l | tr -d '[:space:]')
    if [ $GIT_STATUS = 0 ]
    then
      BUILD_CLEAN=yes
    else
      BUILD_CLEAN=no
    fi

    if [ $BUILD_CLEAN = yes ]
    then
        echo "commit ok"
        # Increment version
        VERSION=$(cat 'VERSION')
        NEW_VERSION="${VERSION%.*}.$((${VERSION##*.}+1))"
        printf $NEW_VERSION > 'VERSION'

        git add VERSION && git commit -m "VERSION -> $NEW_VERSION"

        GIT_STATUS=$(git status --porcelain | wc -l | tr -d '[:space:]')
        if [ $GIT_STATUS = 0 ]
        then
          BUILD_CLEAN=yes
        else
          BUILD_CLEAN=no
        fi
        if [ $BUILD_CLEAN = yes ]
        then
            # Make all
            make darwin && make linux && make linux
            echo "Je pense que tout est correct"
            scp bin/* root@192.168.1.10:/media/raid10/app/serverUpdate/static
        else
            echo "Problème de commit juste après un changement de version. A vérifier"
        fi
    else
        echo "Tu as oublié de commit et faire le push"
    fi


else
    echo "Erreur de test - Vérifie avec : go test ./.. -v (--cover)"
fi

