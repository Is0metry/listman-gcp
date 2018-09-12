package handlers

// import (
// 	"html/template"
// 	"log"
// 	"net/http"
// 	"os"
// 	//change to "lists" if running in developer server (I don't know why either).
// 	"github.com/is0metry/listman-gcp/lists"

// 	"google.golang.org/appengine"
// 	"google.golang.org/appengine/datastore"
// )

// var wd, _ = os.Getwd()
// var templates = template.Must(template.ParseFiles(wd + "static/list.html"))

// //WebViewHandler handles incoming request to view a list
// func WebViewHandler(w http.ResponseWriter, r *http.Request) {
// 	ctx := appengine.NewContext(r)
// 	keytxt := r.URL.Path[len("/view/"):]
// 	key := new(datastore.Key)
// 	k, err := datastore.DecodeKey(keytxt)
// 	if err != nil {
// 		log.Println("ERROR")
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 	}
// 	key = k
// 	lst, err := lists.GetList(ctx, key)
// 	if err == datastore.ErrNoSuchEntity || err == datastore.ErrInvalidKey {
// 		http.Redirect(w, r, "/", http.StatusFound)
// 	} else if err != nil {
// 		http.Error(w, "Error Getting List: "+err.Error(), http.StatusInternalServerError)
// 	}
// 	templates.ExecuteTemplate(w, "list.html", lst)
// }
// func WebAddHandler(w http.ResponseWriter, r *http.Request) {
// 	var keytxt string
// 	if r.Method == "POST" {
// 		ctx := appengine.NewContext(r)
// 		keytxt = r.URL.Path[len("/add/"):]

// 		text := r.FormValue("newItem")
// 		key, err := datastore.DecodeKey(keytxt)
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)

// 		}
// 		if err := lists.AddItem(ctx, text, key); err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 		}
// 	}
// 	http.Redirect(w, r, "/view/"+keytxt, http.StatusFound)
// }

// func WebDeleteHandler(w http.ResponseWriter, r *http.Request) {
// 	ctx := appengine.NewContext(r)
// 	key, err := datastore.DecodeKey(r.URL.Path[len("/delete/"):])
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 	}
// 	parent, err := lists.DeleteItem(ctx, key)
// 	http.Redirect(w, r, "/view/"+parent, http.StatusFound)
// }
