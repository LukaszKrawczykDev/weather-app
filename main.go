package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

const author = "Łukasz Krawczyk"
const defaultPort = "8080"

const pageHTML = `<!doctype html>
<html lang="pl">
<head>
<meta charset="utf-8">
<title>Pogoda</title>
<style>
 body{font-family:sans-serif;max-width:480px;margin:1.5rem auto}
 select,button{padding:.45rem;margin:.4rem 0}
</style>
</head>
<body>
<h1>Aktualna pogoda</h1>
<form method="post">
 <label>Kraj:
   <select id="country" name="country" onchange="populateCities()">
     <option value="PL">Polska</option>
     <option value="DE">Niemcy</option>
     <option value="GB">Wielka&nbsp;Brytania</option>
   </select>
 </label><br>
 <label>Miasto:
   <select id="city" name="city"></select>
 </label><br>
 <button type="submit">Pokaż pogodę</button>
</form>

{{if .Name}}
<h2>{{.Name}}</h2>
<p style="font-size:1.4rem">{{printf "%.1f" .Main.Temp}} °C — {{(index .Weather 0).Description}}</p>
<img alt="ikona pogody" src="https://openweathermap.org/img/wn/{{(index .Weather 0).Icon}}@2x.png">
{{end}}

<script>
const cities={
 "PL":["Warsaw","Krakow","Gdansk"],
 "DE":["Berlin","Munich","Hamburg"],
 "GB":["London","Manchester","Edinburgh"]
};
function populateCities(){
 const c=document.getElementById("country").value;
 const citySel=document.getElementById("city");
 citySel.innerHTML="";
 (cities[c]||[]).forEach(t=>{
   const o=document.createElement("option");
   o.value=t;o.textContent=t;citySel.appendChild(o);
 });
}
window.addEventListener("DOMContentLoaded",populateCities);
</script>
</body>
</html>`

var tmpl = template.Must(template.New("page").Parse(pageHTML))
var apiKey = os.Getenv("OPENWEATHER_KEY")

type weatherResp struct {
	Main struct{ Temp float64 `json:"temp"` } `json:"main"`
	Weather []struct {
		Description string `json:"description"`
		Icon        string `json:"icon"`
	} `json:"weather"`
	Name string `json:"name"`
}

func main() {
	port := getenv("PORT", defaultPort)
	log.Printf("%s | Autor: %s | Nasłuch na porcie %s", time.Now().Format(time.RFC3339), author, port)
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		_ = tmpl.Execute(w, nil)
		return
	}
	city := r.FormValue("city")
	country := r.FormValue("country")
	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s,%s&units=metric&lang=pl&appid=%s", city, country, apiKey)
	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, "błąd pobierania", http.StatusInternalServerError); return
	}
	defer resp.Body.Close()
	var data weatherResp
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		http.Error(w, "błąd dekodowania", http.StatusInternalServerError); return
	}
	_ = tmpl.Execute(w, data)
}

func getenv(k, d string) string {
	if v := os.Getenv(k); v != "" { return v }
	return d
}