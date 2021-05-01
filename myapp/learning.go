package myapp

import (
	"cloud.google.com/go/firestore"
	"context"
	firebase "firebase.google.com/go"
	"log"
	"search-and-sort-movies/myapp/logger"
)

func connectFirestore() (*firestore.Client, context.Context) {
	ctx := context.Background()
	conf := &firebase.Config{ProjectID: "search-and-sort-movies"}
	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()
	return client, ctx
}

// learning to firestore
func (m *myFile) learningFirestore(videosTry bool) {
	client, ctx := connectFirestore()
	_, _, err := client.Collection("learning").Add(ctx, map[string]interface{}{
		"original": m.fileWithoutDir,
		"complete": m.complete,
		"isOk":     videosTry,
	})

	if err != nil {
		logger.L(logger.Red, "%s", err)
	}
}
