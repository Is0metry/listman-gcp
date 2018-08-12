package main

import (
	"html/template"
	"log"
	"net/http"
	//change to "lists" if running in developer server (I don't know why either).
	"github.com/is0metry/listman-gcp/lists"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

var templates = template.Must(template.ParseFiles("static/list.html"))

//viewHandler handles incoming request to view a list
func viewHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	keytxt := r.URL.Path[len("/view/"):]
	key := new(datastore.Key)
	k, err := datastore.DecodeKey(keytxt)
	if err != nil {
		log.Println("ERROR")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	key = k
	lst, err := lists.GetList(ctx, key)
	if err == datastore.ErrNoSuchEntity || err == datastore.ErrInvalidKey {
		http.Redirect(w, r, "/", http.StatusFound)
	} else if err != nil {
		http.Error(w, "Error Getting List: "+err.Error(), http.StatusInternalServerError)
	}
	templates.ExecuteTemplate(w, "list.html", lst)
}
func addHandler(w http.ResponseWriter, r *http.Request) {
	var keytxt string
	if r.Method == "POST" {
		ctx := appengine.NewContext(r)
		keytxt = r.URL.Path[len("/add/"):]

		text := r.FormValue("newItem")
		key, err := datastore.DecodeKey(keytxt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)

		}
		if err := lists.AddItem(ctx, text, key); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
	http.Redirect(w, r, "/view/"+keytxt, http.StatusFound)
}
func rootHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	key, err := lists.GetRoot(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	log.Println(key)
	http.Redirect(w, r, "/view/"+key, http.StatusFound)
}
func deleteHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	key, err := datastore.DecodeKey(r.URL.Path[len("/delete/"):])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	parent, err := lists.DeleteItem(ctx, key)
	http.Redirect(w, r, "/view/"+parent, http.StatusFound)
}
func main() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/add/", addHandler)
	http.HandleFunc("/delete/", deleteHandler)
	appengine.Main()
}
