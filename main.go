package main

import (
	"bytes"
	"log"
	"net/http"
	"net/url"

	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
)

// Configuration contains the necessary bits and pieces to connect your application correctly to Thinglink.
type Configuration struct {
	ClientID     string
	ClientSecret string
	UserID       string
	RedirectURI  string
}

var templates = template.Must(template.ParseFiles("assets/index.html"))
var config Configuration

//  Handles the callbacks from the TL Connect process. You can either get a "code" or an "error" in the
//  query parameters, and in case of an error, there's a more descriptive "error_description" as well.
//  It's up to you to decide how to handle them.
func callback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	error := r.URL.Query().Get("error")

	if error != "" {
		desc := r.URL.Query().Get("error_description")
		fmt.Fprintf(w, "Authentication failed.\nReason: %s\nDescription: %s", error, desc)
	} else {
		// We got the code, so we can now exchange the code for an access token.
		// This is not necessary for TL Connect to work - if you don't need the access token, by all means
		// ignore this following code.
		// See https://thinglink.docs.apiary.io/#reference/authentication for further details.
		resp, err := http.PostForm("https://www.thinglink.com/auth/token",
			url.Values{
				"client_id":     {config.ClientID},
				"client_secret": {config.ClientSecret},
				"grant_type":    {"authorization_code"},
				"code":          {code},
				"redirect_uri":  {config.RedirectURI},
			})

		if err != nil {
			fmt.Fprintf(w, "Oops, unable to exchange the code for an access token.")
		} else {
			body, _ := ioutil.ReadAll(resp.Body)
			var prettyJSON bytes.Buffer

			json.Indent(&prettyJSON, body, "", "\t")

			fmt.Fprintf(w,
				`OK, got back the following response. You can close this window now. Normally, you would autoclose this
window, e.g. with window.close() in Javascript,	and store the access_token and refresh_token
for API calls later on.

You may also receive an error in case the token exchange fails; for example, when you have an incorrect
client secret.

%s`,
				prettyJSON.Bytes())
		}
	}

}

// Respond to generic page rendering requests. We just render the index.html template for
// all possible requests.
func handler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "index")
}

// Use Go's built-in templating mechanism to render a response page.
func renderTemplate(w http.ResponseWriter, tmpl string) {
	err := templates.ExecuteTemplate(w, tmpl+".html", config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {

	// Configuration
	config.ClientID = os.Getenv("CLIENTID")
	config.UserID = os.Getenv("USERID")
	config.ClientSecret = os.Getenv("CLIENTSECRET")

	config.RedirectURI = "http://localhost:8000/callback"

	if config.UserID == "" {
		config.UserID = "sample_id"
	}

	// Set up the web server
	http.Handle("/gfx/", http.FileServer(http.Dir("./assets/")))
	http.HandleFunc("/callback", callback)
	http.HandleFunc("/", handler)

	fmt.Println("Server is now running at http://localhost:8000/")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
