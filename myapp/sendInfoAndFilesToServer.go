package myapp

import (
	"bytes"
	"encoding/json"
	"github.com/denisbrodbeck/machineid"
	"log"
	"net"
	"net/http"
	"regexp"
	"search-and-sort-movies/myapp/model"
	"time"
)

type SendAll struct {
	AllFile  model.AllFiles `json:"all_file"`
	UniqueId string         `json:"unique_id"`
	IP       net.IP         `json:"ip"`
}

func Send() {
	// En Dev
	var url = "http://localhost:9999" + "/info"

	// En prod
	//var url = UrlUpdateURL + "/info"
	var sendAll SendAll
	sendAll.AllFile = getAllFiles()
	sendAll.UniqueId = getUniqueIdPc()
	sendAll.IP = ipLocal()

	j, err := json.Marshal(sendAll)
	if err != nil {
		log.Println(err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(j))
	if err != nil {
		log.Println(err)
	}
	req.Header.Set("X-Custom-Header", "sendAllInfo")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		time.Sleep(1 * time.Minute)
		Send()
	}
	defer resp.Body.Close()

	log.Println("response Status:", resp.Status)
	log.Println("response Headers:", resp.Header)
	//body, _ := ioutil.ReadAll(resp.Body)
	err = json.NewDecoder(resp.Body).Decode(&sendAll)
	if err != nil {
		log.Println(err)
	}
	prettyJson, err := json.MarshalIndent(&sendAll, "", " ")
	if err != nil {
		log.Println(err)
	}
	log.Println("response Body:", string(prettyJson))

}

func getUniqueIdPc() string {
	id, err := machineid.ID()
	if err != nil {
		log.Fatal(err)
	}
	return id
}
func ipLocal() net.IP {
	ifaces, _ := net.Interfaces()
	// handle err
	var re = regexp.MustCompile(`(?m)192.168.\d+.\d+`)
	for _, i := range ifaces {
		addrs, _ := i.Addrs()
		// handle err
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip != nil {
				if re.MatchString(ip.To4().String()) {
					return ip
				}
			}
		}
	}
	return nil
}
