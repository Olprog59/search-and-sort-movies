package myapp

import (
	"bytes"
	"encoding/json"
	"github.com/denisbrodbeck/machineid"
	"log"
	"net"
	"net/http"
	"reflect"
	"regexp"
	"search-and-sort-movies/myapp/model"
	"time"
)

type User struct {
	Video    model.Video `json:"video"`
	UniqueId string      `json:"unique_id"`
	IPLocal  net.IP      `json:"ip_local"`
	IPWan    string      `json:"ip_wan"`
}

var user User

func PostInfo() {
	if send() != nil {
		time.Sleep(1 * time.Minute)
		PostInfo()
	}
}

func send() error {
	var url = UrlUpdateURL + "/info"
	var user2 User

	user2.Video = getVideos()
	user2.UniqueId = getUniqueIdPc()
	user2.IPLocal = ipLocal()
	user2.IPWan = ipWan()

	if reflect.DeepEqual(user, user2) {
		log.Println("Pas de changements")
		user2 = User{}
		return nil
	}
	user = user2

	j, err := json.Marshal(user)
	if err != nil {
		log.Println(err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(j))
	if err != nil {
		log.Println(err)
		return err
	}
	defer req.Body.Close()
	req.Header.Set("X-Custom-Header", "sendAllInfo")
	req.Header.Set("Content-Type", "application/json")

	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout: 10 * time.Second,
		}).DialContext,
	}
	client := &http.Client{Transport: transport}
	resp, err := client.Do(req)

	if err != nil {
		log.Println(err)
		return err
	}
	defer resp.Body.Close()

	log.Println("response Status:", resp.Status)
	log.Println("response Headers:", resp.Header)

	//err = json.NewDecoder(resp.Body).Decode(&user)
	//if err != nil {
	//	log.Println(err)
	//	return err
	//}
	//prettyJson, err := json.MarshalIndent(&user, "", " ")
	//if err != nil {
	//	log.Println(err)
	//	return err
	//}
	//log.Println("response Body:", string(prettyJson))
	return nil
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

func ipWan() string {
	resp, err := http.Get("https://ifconfig.me/all.json")
	if err != nil {
		log.Println("Pas possible d'accéder à ifconfig.me/all.json")
	}
	defer resp.Body.Close()
	var ip model.IPWan
	_ = json.NewDecoder(resp.Body).Decode(&ip)
	return ip.IPAddr
}
