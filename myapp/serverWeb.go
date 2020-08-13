package myapp

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"syscall"
	"time"
)

type HttpCheckPort struct {
	Port    string
	Message string
}

func checkIfPortIsUsed(port string) bool {
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return true
	}
	defer ln.Close()

	return false
}

func ServerHttp() {
	var port = checkEnvPort()

	fmt.Printf("Link Application: http://%s:%s/config\n", IpLocal(), port)

	//http.HandleFunc("/", homeHandler)

	http.HandleFunc("/config", configHandler)
	http.HandleFunc("/port", portHandler)
	http.HandleFunc("/", errorHandler)

	log.Printf("Serveur http démarré sur le port %s\n", port)
	log.Println(http.ListenAndServe(":"+port, nil))
}

func errorHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	_, _ = fmt.Fprint(w, errorHtml)
}
func configHandler(w http.ResponseWriter, r *http.Request) {
	pwd, _ := os.Getwd()
	if r.Method == http.MethodGet {
		_, _ = fmt.Fprintf(w, configHtml, pwd, GetEnv("port"), GetEnv("dlna"), GetEnv("movies"), GetEnv("series"))
	} else if r.Method == http.MethodPost {
		dlna := r.FormValue("dlna")
		port := r.FormValue("port")
		movies := r.FormValue("movies")
		series := r.FormValue("series")

		changeInfo := false

		if dlna != "" {
			if GetEnv("dlna") != dlna {
				SetEnv("dlna", dlna)
				changeInfo = true
			}
		}
		if port != "" {
			if GetEnv("port") != port {
				SetEnv("port", port)
				log.Printf("new link : http://%s:%s/config\n", IpLocal(), port)
				changeInfo = true
			}
		}
		if movies != "" {
			if GetEnv("movies") != movies {
				SetEnv("movies", movies)
				changeInfo = true
			}
		}
		if series != "" {
			if GetEnv("series") != series {
				SetEnv("series", series)
				changeInfo = true
			}
		}
		if changeInfo {
			go func() {
				time.Sleep(2 * time.Second)
				restartApp()
			}()
		}
		_, _ = fmt.Fprintf(w, configHtml, pwd, GetEnv("port"), GetEnv("dlna"), GetEnv("movies"), GetEnv("series"))
	}

}

func portHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	if r.Method == http.MethodGet {
		keys, ok := r.URL.Query()["port"]
		if !ok || len(keys[0]) < 1 {
			js, err := json.Marshal(HttpCheckPort{Port: "", Message: "Url Param 'port' is missing."})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			_, _ = w.Write(js)
			return
		}
		portHttp := keys[0]
		i, _ := strconv.Atoi(portHttp)
		if i > 65535 || i < 1 {
			js, err := json.Marshal(HttpCheckPort{Port: "", Message: fmt.Sprintf("This port %s is not good : 1 - 65 535", portHttp)})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			//_, _ = fmt.Fprintf(w, "{ 'port': '', 'response': This port %s is already used or is > 65 535}", portHttp)
			_, _ = w.Write(js)
			return
		}
		if checkIfPortIsUsed(portHttp) {

			js, err := json.Marshal(HttpCheckPort{Port: "", Message: fmt.Sprintf("This port %s is already used.", portHttp)})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			//_, _ = fmt.Fprintf(w, "{ 'port': '', 'response': This port %s is already used or is > 65 535}", portHttp)
			_, _ = w.Write(js)
			return
		}

		js, err := json.Marshal(HttpCheckPort{Port: portHttp, Message: ""})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		//_, _ = fmt.Fprintf(w, "{ 'port': '', 'response': This port %s is already used or is > 65 535}", portHttp)
		_, _ = w.Write(js)
		return
	}
}

func checkEnvPort() (port string) {
	if GetEnv("port") == "" {
		SetEnv("port", "4567")
		port = "4567"
	} else {
		port = GetEnv("port")
	}
	return port
}

func restartApp() {
	var binary string
	var lookErr error

	var args []string
	if runtime.GOOS == "linux" {
		binary, lookErr = exec.LookPath("systemctl")
		args = []string{
			"systemctl",
			"restart",
			"searchAndSortMovies.service",
		}
	} else if runtime.GOOS == "darwin" {
		binary, lookErr = exec.LookPath("bash")
		args = []string{
			binary,
			"-c",
			"./search-and-sort-movies",
			//"./search-and-sort-movies-darwin-amd64",
		}
	}

	if lookErr != nil {
		panic(lookErr)
	}
	err := syscall.Exec(binary, args, os.Environ())
	if err != nil {
		log.Println(err.Error())
	}
}

const configHtml = `
<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<title>Config Application</title>
	<link href="https://fonts.googleapis.com/css?family=Nunito+Sans&display=swap" rel="stylesheet">
	<style>
		:root{
			--bgc: #101010;
			--color: #FEFEFE;
			--font-size: 2vw;
		}
		*{
			box-sizing: border-box;
			//transition: all 50ms ease-in;
		}

		html,
		body{
			margin: 0 auto;
			padding: 0;
			font-size: 2vw;
			font-family: 'Nunito Sans', sans-serif;
			background-color: var(--bgc);
		}

		.container{
			margin: 0 auto;
			max-width: 1200px;
		}

		.container h1, .container h3{
			text-align: center;
			filter: invert(1);
			mix-blend-mode: difference;
		}

		.container form label{
			filter: invert(1);
			mix-blend-mode: difference;
		}

		.container form .form-row{
			min-height: 120px;
			margin: 5px;
			padding: 5px;
		}

		.container form .form-row div{
			width: 100%%;
		}
		.container form .form-row div input{
			line-height: 50px;
			font-size: 2vw;
			outline: none;
			width: 100%%;
			border: 4px solid transparent;
			background-color: var(--color);
			color: var(--bgc);
		}

		.container form .form-row div p {
			margin: 0;
			padding: 0;
			font-size: 0.9em;
			color: var(--color);
		}

		.container form .form-row div button{
			line-height: 50px;
			width: 200px;
			outline: none;
			border: none;
			color: #666;
			background-color: #34f19a;
			font-size: 20px;
		}

		.container form .form-row div button:disabled {
			opacity: 0.3;
		}

		@media screen and (min-width: 1200px) {
			html,body {
				font-size: 25px;
			}
			.container form .form-row div input{
				line-height: 50px;
				font-size: 20px;
				width: 100%%;
			}
		}
		@media screen and (max-width: 800px) {
			html,body {
				font-size: 15px;
			}
			.container form .form-row div input{
				line-height: 30px;
				font-size: 15px;
				width: 100%%;
			}
			.container form .form-row{
				min-height: 60px;
			}
		}
	</style>
</head>
<body>
<div class="container">
	<h1>Configuration de l'application Search and sort Movies</h1>
	<hr>
	<h3>Emplacement de l'application : <div id="location-app">%s</div></h3>
	<hr>
	<form action="/config" method="post" id="form">
		<div class="form-row">
			<label for="port">Port actuellement utiliser pour les pages Web</label>
			<div>
				<input type="text" name="port" id="port" placeholder="Port des pages Web" value="%s">
				<p id="error-port"></p>
			</div>
		</div>
		<div class="form-row">
			<label for="dlna">Emplacement des fichiers à trier</label>
			<div>
				<input type="text" name="dlna" id="dlna" placeholder="Emplacement des fichiers à trier" value="%s">
			</div>
		</div>
		<div class="form-row">
			<label for="movies">Dossier des films triés</label>
			<div>
				<input type="text" name="movies" id="movies" placeholder="Dossier des films triés" value="%s">
			</div>
		</div>
		<div class="form-row">
			<label for="series">Dossier des séries triées</label>
			<div>
				<input type="text" name="series" id="series" placeholder="Dossier des séries triées" value="%s">
			</div>
		</div>
		<div class="form-row">
			<div>
				<button type="submit" id="submit" disabled>Valider</button>
			</div>
		</div>
	</form>
</div>
<script>
    (function () {
        const idForm = document.getElementById("form");
        const idButton = document.getElementById("submit");
        const idPort = document.getElementById("port");
        const idDlna = document.getElementById("dlna");
        const idMovies = document.getElementById("movies");
        const idSeries = document.getElementById("series");

        localStorage.setItem("idButton", idButton.value);
        localStorage.setItem("idPort", idPort.value);
        localStorage.setItem("idDlna", idDlna.value);
        localStorage.setItem("idMovies", idMovies.value);
        localStorage.setItem("idSeries", idSeries.value);

        idForm.onsubmit = (e) => {
            const c =
                confirm("Attention, en changeant les informations, l'application va redémarrer quelques secondes après la validation.");
            if (!c) {
                e.preventDefault();
            }

        };

        idDlna.addEventListener('keyup', (e) => {
            if (localStorage.getItem("idDlna") === e.target.value) {
                idButton.setAttribute("disabled", "true");
            } else {
                idButton.removeAttribute("disabled");
            }
        })
        idMovies.addEventListener('keyup', (e) => {
            if (localStorage.getItem("idMovies") === e.target.value) {
                idButton.setAttribute("disabled", "true");
            } else {
                idButton.removeAttribute("disabled");
            }
        })
        idSeries.addEventListener('keyup', (e) => {
            if (localStorage.getItem("idSeries") === e.target.value) {
                idButton.setAttribute("disabled", "true");
            } else {
                idButton.removeAttribute("disabled");
            }
        })

        // Check Port if ok
        idPort.addEventListener('keyup', (e) => {
            fetch("http://"+window.location.host+"/port?port=" + e.target.value)
                .then(response => {
                    const contentType = response.headers.get("content-type");
                    if (contentType && contentType.indexOf("application/json") !== -1) {
                        return response.json().then(function (json) {
                            const port = json.Port;
                            const message = json.Message;
                            console.log(port);
                            if (port === "") {
                                idButton.setAttribute("disabled", "true");
                                idPort.style.borderColor = "red";
                                document.getElementById("error-port").innerText = message;
                            } else {
                                document.getElementById("error-port").innerText = '';
                                idButton.removeAttribute("disabled");
                                idPort.style.borderColor = "green";
                            }
                        });
                    } else {
                        console.log("Oops, nous n'avons pas du JSON!");
                    }
                })
        });
    })();
</script>
</body>
</html>
`

const errorHtml = `
<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<title>Error</title>
	<link href="https://fonts.googleapis.com/css?family=Nunito+Sans&display=swap" rel="stylesheet">
	<style>
		:root {
			--bgc: #101010;
			--color: #FEFEFE;
		}

		* {
			box-sizing: border-box;
		}

		html,
		body {
			margin: 0 auto;
			padding: 0;
			font-size: 1.5vw;
			font-family: 'Nunito Sans', sans-serif;
			background-color: var(--bgc);
			color: var(--color);
			text-align: center;
		}

		.container div {
			position: fixed;
			top: 50vh;
			left: 50vw;
			width: 90vw;
			padding: 0;
			margin: 0;
			transform: translate(-50%, -50%);
			font-size: 40px;
		}
	</style>
</head>
<body>
<div class="container">
	<div>
		<p>
			Cette page n'est pas activée. Peut-être plus tard.
		</p>
		<p>
			Demande à l'administrateur de ce site. Il t'en dira plus.
		</p>
		<p>😜</p>
	</div>
</div>
</body>
</html>
`
