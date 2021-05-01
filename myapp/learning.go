package myapp

import (
	"context"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
	"log"
	"search-and-sort-movies/myapp/logger"
)

// learning to firestore
func (m *myFile) learningFirestore(videosTry bool) {
	ctx := context.Background()
	sa := option.WithCredentialsJSON([]byte(`{
  "type": "service_account",
  "project_id": "search-and-sort-movies",
  "private_key_id": "7231c7ecbdaf1ff4bc4dc6e076e62969acfd3b95",
  "private_key": "-----BEGIN PRIVATE KEY-----\nMIIEvAIBADANBgkqhkiG9w0BAQEFAASCBKYwggSiAgEAAoIBAQCO048PjWhhzMGb\nY72+59dMaHOuwlwxfF0YIt9Ji9jmBh0jUPgFN5NC1AWifCQxQzNpO2KN3VpzXGVu\no2VOIAM1r2d53Mb7tk82tmqz+G7q5EvbDGYVhOtVmfrBV4HmAgtNZN0fOaXgBkPA\nTJasxTSt6JMVE8squBO6v9Rjsm4/ABXgpUErPVa/M02wGvcu3na7WpmN8RWuvsty\nFXbwOD3WeS1JbD0GsNobNH+EdoBKvZFw4sWIzkvJAzpJ7HUsGeEARVNuKbk+IygJ\n05vUhWZV/yrnnYKnsIzpxbH+oUhsuJx/MtOkO9J/xeQ0vMASD8ryIp4CzeRJpPvu\nV8y2VlLZAgMBAAECgf8Ft/LCbKcBQ3F9B0FRo3a7WHNJS/7k/FW94aePCkpk61m/\nUD6d9tuukU6octAAKPCmBfqsM3Crrcdh5qjnt+BpBdzexDgzW9lQF32uNwQ542om\nGxj7Q5colkGz/Az3aJ2LO76ewXAiyB1uVaKyyDKTw/wfm+tGipg/oe6fD0Xav9bM\nkzIqTVS6OfKtsvAGTE2syi0MZ255aj4TBwb1XcmTUQibbIn/yY1gDxZWWcIb0xHW\nEhkG+iDLVtvrhC8NFWnb3Nrja5gwWD6+yLq/sj9pPvG9pNORw6xRUExuiR3DfHMp\nJWi6IBgyFLzgzn01zQidDvF6fGKxeIWDDTv1qmECgYEAx/Dh9tYf1Z41I025RIIS\nPqgw1Zdqe0AB4buudNm0FAt0JKWQeMetkhU3kCSMmKXPZvgxxrtdhHRVH8+mJZEk\n4PooXwQXW+acdt6h07Ol3fQbOSad8koUDLVJ6fIpx+rjdLPVpoHaExfmHzXwY27S\na2EC8cM46x6M6rxvWijpAekCgYEAtt8tj0JuxRMSwp0mmlv8OL+qXduNWMn8Lqlq\nt4NSULWxI4WypvIttFqgNqYl9SX8cWJ0NI0h+L+kB4OhYKcCXx5vBx6DLuinj7PU\nwQ/+1t7tzJ+TvP2va20CPsvs76r7iAdn655mNERFVaAH31dQQnwYAr46tkFqHh1Y\n7yXIw3ECgYEAnLeR1nFlyIHWctKUOj+d32Djzjd97hdwoigDCXIu9Vs48RSZFiKl\nSRC6WZBcZ7XnyHUYRwZLueuZYXLYby/CcVmDVV6WlKFA2OeOfqqcg0m4IObE/MnV\nx1Q+GFKJLztMiAgBmh7D+R1Ncf9MahPOeP40WZ3Eun5axA9pVIkmgikCgYEArAAj\n8O/SEeiLp6J58Yt9Ir8bdaYQPyfT2uucJTkODj2me7u/ughk9pKayGvjnb15wAeT\nNu5buoQ9upeTDL5om6CbWz3WsyM+nwnMnT33OpB5aBHbulF9UfQ4vWm+0/mlFV+p\n3dKhXJ2t/QhE/0s3gSEI0GOuA0hpkCOYR27pcvECgYAuFLa9MNbs00WVdCHzgmC5\ng7ZeCuaS7iUsQzUhXUonzjnn9K1TVVbctTrWihgsDxe3qLtZfp+5CkUY3sza4GEz\nXJMpz0HrRPteA72uMTgxZTBTiSIebAtlxVLpv+yBL5aEVX/e5v1r0OfmQkf6fYTg\nP4LEZlaVlV/uX2y+HrUPBg==\n-----END PRIVATE KEY-----\n",
  "client_email": "search-and-sort-movies@appspot.gserviceaccount.com",
  "client_id": "100119263118153870570",
  "auth_uri": "https://accounts.google.com/o/oauth2/auth",
  "token_uri": "https://oauth2.googleapis.com/token",
  "auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
  "client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/search-and-sort-movies%40appspot.gserviceaccount.com"
}`))
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		logger.L(logger.Red, "%s", err)
	}
	_, _, err = client.Collection("learning").Add(ctx, map[string]interface{}{
		"original": m.fileWithoutDir,
		"complete": m.complete,
		"isOk":     videosTry,
	})

	if err != nil {
		logger.L(logger.Red, "%s", err)
	}
	err = client.Close()
	if err != nil {
		return
	}
}
