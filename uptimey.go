package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/shirou/gopsutil/host"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
)

type BingImage struct {
	Images []struct {
		Bot           int    `json:"bot"`
		Copyright     string `json:"copyright"`
		Copyrightlink string `json:"copyrightlink"`
		Drk           int    `json:"drk"`
		Enddate       string `json:"enddate"`
		Fullstartdate string `json:"fullstartdate"`
		Hs            []struct {
			Desc  string `json:"desc"`
			Link  string `json:"link"`
			Locx  int    `json:"locx"`
			Locy  int    `json:"locy"`
			Query string `json:"query"`
		} `json:"hs"`
		Hsh       string        `json:"hsh"`
		Msg       []interface{} `json:"msg"`
		Startdate string        `json:"startdate"`
		Top       int           `json:"top"`
		URL       string        `json:"url"`
		Urlbase   string        `json:"urlbase"`
		Wp        bool          `json:"wp"`
	} `json:"images"`
	Tooltips struct {
		Loading  string `json:"loading"`
		Next     string `json:"next"`
		Previous string `json:"previous"`
		Walle    string `json:"walle"`
		Walls    string `json:"walls"`
	} `json:"tooltips"`
}

type Page struct {
	Title     string
	StartDate time.Time
}

func home(c web.C, w http.ResponseWriter, r *http.Request) {
	p := &Page{Title: "Uptimey"}
	sys, _ := host.HostInfo()

	now := time.Now()
	past := now.Add(-1 * time.Duration(sys.Uptime) * time.Second)

	p.StartDate = past
	tmpl, err := FSString(false, "/assets/index.html")
	if err != nil {
		fmt.Println(err.Error())
	}

	t, _ := template.New("index").Parse(tmpl)
	t.Execute(w, p)
}

func ajax(c web.C, w http.ResponseWriter, r *http.Request) {
	sys, _ := host.HostInfo()
	action := r.FormValue("action")
	if action == "time" {
		now := time.Now()
		past := now.Add(-1 * time.Duration(sys.Uptime) * time.Second)
		fmt.Fprintf(w, "%s;%s;%s", now.Format("January 2, 2006"), now.Format("3:04 pm"), past.Format("January 2, 2006"))
	}
	if action == "uptime" {
		uptimeDuration := time.Duration(sys.Uptime) * time.Second
		days := int(uptimeDuration.Hours() / 24)
		hours := int(uptimeDuration.Hours()) % 24
		minutes := int(uptimeDuration.Minutes()) % 60
		fmt.Fprintf(w, "%d;%d;%d", days, hours, minutes)
	}
	if action == "image" {
		bingResponse, err := http.Get("http://www.bing.com/HPImageArchive.aspx?format=js&idx=0&n=1&mkt=en-US")
		if err != nil {
			fmt.Printf("%s", err)
			os.Exit(1)
		} else {
			defer bingResponse.Body.Close()
			contents, err := ioutil.ReadAll(bingResponse.Body)
			if err != nil {
				fmt.Printf("%s", err)
				os.Exit(1)
			}

			var animals BingImage
			err = json.Unmarshal(contents, &animals)
			if err != nil {
				fmt.Println("error:", err)
			}
			fmt.Fprintf(w, "%s;%s", "https://www.bing.com/"+animals.Images[0].URL, animals.Images[0].Copyright)
		}
	}
	if action == "location" {
		ipResponse, err := http.Get("http://icanhazip.com")
		if err != nil {
			fmt.Printf("%s", err)
			os.Exit(1)
		} else {
			defer ipResponse.Body.Close()
			contents, err := ioutil.ReadAll(ipResponse.Body)
			if err != nil {
				fmt.Printf("%s", err)
				os.Exit(1)
			}
			fmt.Fprintf(w, "%s", string(contents))
		}
	}
}

func main() {

	goji.Get("/", home)
	goji.Get("/script/ajax.php", ajax)
	goji.Get("/assets/*", http.FileServer(FS(false)))
	goji.Serve()
}
