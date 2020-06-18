package main

import (
	"context"
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
    "google.golang.org/grpc"
    "github.com/unrolled/secure"

	pb "github.com/trailofbits/not-going-anywhere/internal/friends"
)

type UserCtx struct {
    Search string
    Users []*pb.Person
}

type PostCtx struct {
	Posts []*pb.Post
}

const (
	serverAddress = "localhost:5000"
	defun         = "trailofbits"
	address       = ":5080"
)

var (
	key          = []byte("xxzoL3R9zA25mztvbm9AWwBdCEqiVvgj")
	sessionStore = sessions.NewCookieStore(key)
)

func cachingHandler(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate, max-age=0")
        next.ServeHTTP(w, r)
    })
}

func main() {
    router  := mux.NewRouter()

    secureMiddleware := secure.New(secure.Options{
        STSSeconds:            31536000,
        STSIncludeSubdomains:  true,
        STSPreload:            true,
        FrameDeny:             true,
        ContentTypeNosniff:    true,
        BrowserXssFilter:      true,
        ContentSecurityPolicy: "script-src $NONCE",
    })

    router.Handle("/", cachingHandler(http.HandlerFunc(indexPage)))
    router.Handle("/posts/add", http.HandlerFunc(addPosts)).Methods("POST")
    router.Handle("/posts", http.HandlerFunc(listPosts)).Methods("GET")
    router.Handle("/posts/{person}", http.HandlerFunc(listPersonPosts)).Methods("GET")
    router.Handle("/friends", http.HandlerFunc(listFriends)).Methods("GET")
    router.Handle("/friends/add/{person}", http.HandlerFunc(addFriend)).Methods("GET")
    router.Handle("/friends/unfriend/{person}", http.HandlerFunc(unfriendFriend)).Methods("GET")
    router.Handle("/register", http.HandlerFunc(registerUser))
    router.Handle("/login", secureMiddleware.Handler(http.HandlerFunc(loginUser)))
	//router.Handle("/logout", http.HandlerFunc(logoutUser))
	http.ListenAndServe(address, router)
}

func indexPage(w http.ResponseWriter, r *http.Request) {
	log.Print("recieved index request")
	session, _ := sessionStore.Get(r, "not-going-anywhere")
	next := r.FormValue("returnURL")

	if next == "" {
		next = "/posts"
	}

	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Redirect(w, r, "/login", 303)
	} else {
		http.Redirect(w, r, next, 303)
	}
}

func addPosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte("<h1>Working</h1>"))
}

func listPosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte("<h1>Posts</h1>"))
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	conn, err := grpc.Dial(serverAddress, grpc.WithInsecure(), grpc.WithBlock())
	session, _ := sessionStore.Get(r, "not-going-anywhere")

	if auth := session.Values["authenticated"].(bool); !auth {
		http.Redirect(w, r, "/login", 303)
    }

	if err != nil {
		log.Fatalf("could not connect to server; make sure you have 'friends_server_net' running: %v", err)
	}

	defer conn.Close()

	client := pb.NewNotGoingAnywhereClient(conn)

	resposts, err := client.GetAllPosts(ctx, &pb.Empty{})

	if err != nil {
		log.Print("failed to load posts")
	}

	posts := PostCtx{Posts: resposts.GetPosts()}

	tmpl, err := template.ParseFiles("templates/main.html")
	if err != nil {
		log.Print("login template failed")
	}
	tmpl.Execute(w, posts)
}

func listPersonPosts(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	personId := vars["person"]

	log.Printf("recieved requestion for person %s", personId)

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte("<h1>Working</h1>"))
}

func listFriends(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/html")
    search := r.FormValue("search")

    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()
    conn, err := grpc.Dial(serverAddress, grpc.WithInsecure(), grpc.WithBlock())

    if err != nil {
        log.Fatalf("could not connect to server; make sure you have 'friends_server_net' running: %v", err)
    }

    defer conn.Close()

    client := pb.NewNotGoingAnywhereClient(conn)
    var people []*pb.Person

    if search == "" {
        rpeople, _ := client.GetPeople(ctx, &pb.Empty{})
        people = rpeople.People
    } else {
        rpeople, _ := client.GetPerson(ctx, &pb.PersonRequest{Uname: search})
        people = rpeople.People
    }

    users := UserCtx{Search: search, Users: people}

    tmpl, err := template.ParseFiles("templates/users.html")
    if err != nil {
        log.Print("list template failed")
    }
    tmpl.Execute(w, users)
}

func addFriend(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	personId := vars["person"]
	log.Printf("recieved requestion for person %s", personId)

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte("<h1>Working</h1>"))
}

func unfriendFriend(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	personId := vars["person"]
	log.Printf("recieved requestion for person %s", personId)

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte("<h1>Working</h1>"))
}

func registerUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	if r.Method == "POST" {
		err := r.ParseForm()

		if err != nil {
			log.Print("fatal parsing error")
		}

		username := r.Form.Get("username")

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		conn, err := grpc.Dial(serverAddress, grpc.WithInsecure(), grpc.WithBlock())

		if err != nil {
			log.Fatalf("could not connect to server; make sure you have 'friends_server_net' running: %v", err)
		}

		defer conn.Close()

		client := pb.NewNotGoingAnywhereClient(conn)

		// check if a user  exists  prior to registering them
		rpeople, err := client.GetPerson(ctx, &pb.PersonRequest{Uname: username})

		if err == nil {
			for _, value := range rpeople.People {
				if value.GetUname() == username {
					w.Write([]byte(`<h1>User Exists</h1>`))
					return
				}
			}
		}

		rperson, err := client.RegisterPerson(ctx, &pb.RegisterRequest{Uname: username})

		if err != nil {
			log.Fatalf("could not register user: %n", err)
		}
		session, _ := sessionStore.Get(r, "not-going-anywhere")
		session.Values["authenticated"] = true
		session.Values["username"] = rperson.GetUname()
		session.Values["uid"] = rperson.GetId()
		session.Save(r, w)

		next := r.FormValue("returnURL")

		if next == "" {
			next = "/posts"
		}

		http.Redirect(w, r, next, 303)
	} else {
		tmpl, err := template.ParseFiles("templates/register.html")
		if err != nil {
			log.Print("login template failed")
		}
		tmpl.Execute(w, nil)
	}
}

func loginUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	if r.Method == "POST" {
		err := r.ParseForm()

		if err != nil {
			log.Print("fatal parsing error")
		}

		username := r.Form.Get("username")

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		conn, err := grpc.Dial(serverAddress, grpc.WithInsecure(), grpc.WithBlock())

		if err != nil {
			log.Fatalf("could not connect to server; make sure you have 'friends_server_net' running: %v", err)
		}

		defer conn.Close()

		client := pb.NewNotGoingAnywhereClient(conn)

		// check if a user  exists  prior to registering them
		rpeople, err := client.GetPerson(ctx, &pb.PersonRequest{Uname: username})

		var rperson *pb.Person

		if err == nil {
			for _, value := range rpeople.People {
				if value.GetUname() == username {
					rperson = value
				}
			}
		}

		if rperson == nil {
			http.Redirect(w, r, "/login", 303)
		}

		session, _ := sessionStore.Get(r, "not-going-anywhere")
		session.Values["authenticated"] = true
		session.Values["username"] = rperson.GetUname()
		session.Values["uid"] = rperson.GetId()
		session.Save(r, w)

		next := r.FormValue("returnURL")

		if next == "" {
			next = "/posts"
		}

		http.Redirect(w, r, next, 303)
	} else {
		tmpl, err := template.ParseFiles("templates/login.html")
		if err != nil {
			log.Print("login template failed")
		}
		tmpl.Execute(w, nil)
	}
}
