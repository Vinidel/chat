package main
import (
	"github.com/stretchr/objx"
 "log"
 "net/http"
 "text/template"
 "path/filepath"
 "sync"
 "flag"
 "trace"
 "os"
 "github.com/stretchr/gomniauth"
 "github.com/stretchr/gomniauth/providers/google"
)

// templ represents a single template
type templateHandler struct {
	once sync.Once
	filename string
	templ *template.Template
}

 // ServeHTTP handles the HTTP request.
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates",t.filename)))

		data := map[string]interface{}{
			"Host": r.Host,
		}

		if authCookie, err := r.Cookie("auth"); err == nil {
			data["UserData"] = objx.MustFromBase64(authCookie.Value)
		}

		t.templ.Execute(w, data)
	})
}

func main() {
	addr := flag.String("addr", ":8080", "The addr of the  application.")
	flag.Parse() // parse the flags
	// r := newRoom(UseAuthAvatar)
	r := newRoom(UseGravatarAvatar)

// setup gomniauth
gomniauth.SetSecurityKey("viniciusisawesome")
gomniauth.WithProviders(
	google.New("2313323659-h6478p0e0vvinuduoloduc091v0jc8tp.apps.googleusercontent.com", "pWdgRbEvUZRULULJogOxC1Qd",
		"http://vinichat.com/auth/callback/google"),
)

	r.tracer = trace.New(os.Stdout)
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/room", r)
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHander)
	http.HandleFunc("/logout", func(w http.ResponseWriter, r  *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:   "auth",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		})
		w.Header().Set("Location", "/chat")
		w.WriteHeader(http.StatusTemporaryRedirect)
	}) 
	// get the room going
	go r.run()
	// start the web server
	log.Println("Starting web server on", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}