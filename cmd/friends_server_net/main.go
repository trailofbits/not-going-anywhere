package main

import (
    "context"
    "net"
    "log"
    "database/sql"
    "os"
    "fmt"

    // external dependencies
    "google.golang.org/grpc"
    "github.com/google/uuid"
    _ "github.com/mattn/go-sqlite3"

    // our ProtoBuffers internal libraries
    pb "github.com/trailofbits/not-going-anywhere/internal/friends"
)

const (
    port = ":5000"
    defaultdbFile = "./db/friends.db"
)

var (
    databaseHandle *sql.DB
)

type server struct {
    pb.UnimplementedNotGoingAnywhereServer
    dbHnd *sql.DB
}

func (s *server) RegisterPerson(ctx context.Context, in *pb.RegisterRequest) (*pb.Person, error) {
    log.Printf("received: %s", in.GetUname())
    uid, _ := uuid.NewRandom()

    query := fmt.Sprintf(`INSERT INTO people (username, uid) VALUES ("%s", "%s")`, in.GetUname(), uid.String())

    log.Printf("query: %s", query)

    statement, err := s.dbHnd.Prepare(query)

    if err != nil {
        log.Fatalf("error adding user: %v", err)
    }

    statement.Exec()

    return &pb.Person{Id: uid.String(), Uname: in.GetUname()}, nil
}

func (s *server) GetPeople(ctx context.Context, in *pb.Empty) (*pb.People, error) {
    log.Print("received request for people")

    rows, err := s.dbHnd.Query("SELECT username, uid FROM people")

    if err != nil || rows == nil {
        log.Fatalf("error occured: %v %v", err, rows)
    }

    ret := pb.People{}
    var uname string
    var uid string

    for rows.Next() {
        rows.Scan(&uname, &uid)
        log.Printf("found: %s, %s", uname, uid)
        ret.People = append(ret.People, &pb.Person{Id: uid, Uname: uname})
    }

    return &ret, nil
}

func (s *server) GetPerson(ctx context.Context, in *pb.PersonRequest) (*pb.People, error) {
    log.Print("received request for people")

    rows, err := s.dbHnd.Query(fmt.Sprintf(`SELECT username, uid FROM people WHERE username like "%%%s%%"`, in.GetUname()))

    if err != nil || rows == nil {
        log.Fatalf("error occured: %v %v", err, rows)
    }

    ret := pb.People{}
    var uname string
    var uid string

    for rows.Next() {
        rows.Scan(&uname, &uid)
        log.Printf("found: %s, %s", uname, uid)
        ret.People = append(ret.People, &pb.Person{Id: uid, Uname: uname})
    }

    return &ret, nil
}

func (s *server) GetFriends(ctx context.Context, in *pb.Person) (*pb.People, error) {
    log.Print("list user's friends")
    query := `select people.uid, people.username from friends join people on friends.userIDTo = people.uid where friends.userIDFrom = ?`
    rows, err := s.dbHnd.Query(query, in.GetId())

    if err != nil {
        log.Fatalf("error selecting data: %v")
    }

    ret := pb.People{}

    var uname string
    var uid string

    for rows.Next() {
        rows.Scan(&uid, &uname)
        log.Printf("found: %s, %s", uname, uid)
        ret.People = append(ret.People, &pb.Person{Id: uid, Uname: uname})
    }

    return &ret, nil
}

func (s *server) GetPosts(ctx context.Context, in *pb.Person) (*pb.Posts, error) {
    log.Print("received request for people")

    rows, err := s.dbHnd.Query(fmt.Sprintf(`SELECT username, postID, post FROM posts where username = "%s"`, in.GetUname()))

    if err != nil || rows == nil {
        log.Fatalf("error occured: %v %v", err, rows)
    }

    ret := pb.Posts{}
    var uname string
    var pid string
    var text string

    for rows.Next() {
        rows.Scan(&uname, &pid, &text)
        log.Printf("found: %s, %s, %s", uname, pid, text)
        ret.Posts = append(ret.Posts, &pb.Post{Id: pid, Uname: uname, Text: text})
    }

    return &ret, nil
}

func (s *server) GetAllPosts(ctx context.Context, in *pb.Empty) (*pb.Posts, error) {
    log.Print("received request for people")

    rows, err := s.dbHnd.Query(`SELECT username, postID, post FROM posts`)

    if err != nil || rows == nil {
        log.Fatalf("error occured: %v %v", err, rows)
    }

    ret := pb.Posts{}
    var uname string
    var pid string
    var text string

    for rows.Next() {
        rows.Scan(&uname, &pid, &text)
        log.Printf("found: %s, %s, %s", uname, pid, text)
        ret.Posts = append(ret.Posts, &pb.Post{Id: pid, Uname: uname, Text: text})
    }

    return &ret, nil
}

func (s *server) AddPost(ctx context.Context, in *pb.NewPost) (*pb.TrueForSuccess, error) {
    log.Printf("recieved: post %s for user %v", in.GetPost(), in.GetUser())
    ret := pb.TrueForSuccess{Status: false}

    post := in.GetPost()
    username := in.User.GetUname()
    pid, _ := uuid.NewRandom()

    query := fmt.Sprintf(`INSERT INTO posts (username, post, postID) VALUES ("%s", "%s", "%s")`, username, post, pid.String())

    log.Printf("query: %s", query)

    statement, err := s.dbHnd.Prepare(query)

    if err != nil {
        return &ret, err
    }

    statement.Exec()

    ret.Status = true

    return &ret, nil
}

func (s *server) AddFriend(ctx context.Context, in *pb.Friendship) (*pb.TrueForSuccess, error) {
    originalUser := in.GetOrig()
    friendUser := in.GetFriend()

    statement, err := s.dbHnd.Prepare(`INSERT INTO friends (userIDFrom, userIDTo) VALUES (?, ?)`)
    if err  != nil {
        log.Fatalf("error preparing friendship statement")
    }

    statement.Exec(originalUser.GetId(), friendUser.GetId())

    return &pb.TrueForSuccess{Status: true}, nil
}

func (s *server) Unfriend(ctx context.Context, in *pb.Friendship) (*pb.TrueForSuccess, error) {
    originalUser := in.GetOrig()
    friendUser := in.GetFriend()

    statement, err := s.dbHnd.Prepare(`DELETE FROM friends WHERE userIDFrom = ? and userIDTo = ?`)
    if err  != nil {
        log.Fatalf("error preparing friendship statement")
    }

    statement.Exec(originalUser.GetId(), friendUser.GetId())

    return &pb.TrueForSuccess{Status: true}, nil
}

func main() {
    netserv, err := net.Listen("tcp", port)
    if err != nil {
        log.Fatalf("failed to bind to port %d: %v", port, err)
    }

    dbFile := defaultdbFile

    if len(os.Args) > 1 {
        dbFile = os.Args[1]
    }

    databaseHandle, err := sql.Open("sqlite3", dbFile)
    statement, err := databaseHandle.Prepare("CREATE TABLE IF NOT EXISTS people (id INTEGER PRIMARY KEY, username TEXT UNIQUE, uid TEXT)")
    if err != nil {
        log.Fatalf("database error: %v", err)
        panic(err)
    }
    statement.Exec()

    statement, err = databaseHandle.Prepare("CREATE TABLE IF NOT EXISTS friends (id INTEGER PRIMARY KEY, userIDFrom TEXT, userIDTo TEXT)")
    statement.Exec()

    statement, err = databaseHandle.Prepare("CREATE TABLE IF NOT EXISTS posts  (id INTEGER PRIMARY KEY, username TEXT, post TEXT, postID TEXT)")
    statement.Exec()

    defer databaseHandle.Close()
    grpcserv := grpc.NewServer()
    pb.RegisterNotGoingAnywhereServer(grpcserv, &server{dbHnd: databaseHandle})

    if err := grpcserv.Serve(netserv); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}
