package main

import (
    "context"
    "log"
    "os"
    "time"

    "google.golang.org/grpc"
    pb "github.com/trailofbits/not-going-anywhere/internal/friends"
)

const (
    address = "localhost:5000"
    defun = "trailofbits"
)

func main() {
    conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())

    if err != nil {
        log.Fatalf("could not connect to server; make sure you have 'friends_server_net' running: %v", err)
    }

    defer conn.Close()

    client := pb.NewNotGoingAnywhereClient(conn)

    cmd := "people"

    if len(os.Args) > 1 {
        cmd = os.Args[1]
    }

    username := defun
    if len(os.Args) > 2 {
        username = os.Args[2]
    }

    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()

    switch cmd {
        case "A", "addpost":
            log.Print("add post")
            log.Print("first,  get the user's ident object")
            rpeople, err := client.GetPerson(ctx, &pb.PersonRequest{Uname: username})

            if err != nil {
                log.Fatalf("error retrieving user: %v", err)
            }

            var id *pb.Person
            if len(rpeople.People) > 0 {
                id = rpeople.People[0]
            }  else {
                log.Fatal("no matching users found")
            }

            post := ""
            if len(os.Args) > 3 {
                post = os.Args[3]
            }

            result, err := client.AddPost(ctx, &pb.NewPost{User: id, Post: post})
            if err != nil {
                log.Fatalf("couldn't add post: %v", err)
            }

            if result.Status {
                log.Print("post added successfully")
            } else {
                log.Print("post failed to add")
            }
        case "F", "friend":
            log.Print("request friend")

            rpeople, err := client.GetPerson(ctx, &pb.PersonRequest{Uname: username})

            if err != nil {
                log.Fatalf("error retrieving user: %v", err)
            }

            var id0 *pb.Person

            if len(rpeople.People) > 0 {
                id0 = rpeople.People[0]
            }  else {
                log.Fatal("no matching users found")
            }

            username = os.Args[3]

            rpeople, err = client.GetPerson(ctx, &pb.PersonRequest{Uname: username})

            if err != nil {
                log.Fatalf("error retrieving user: %v", err)
            }

            var id1 *pb.Person

            if len(rpeople.People) > 0 {
                id1 = rpeople.People[0]
            }  else {
                log.Fatal("no matching users found")
            }

            rfriendship, err := client.AddFriend(ctx, &pb.Friendship{Orig: id0, Friend: id1})

            if rfriendship.Status {
                log.Print("friendship added")
            } else {
                log.Print("failed to add friendship")
            }
        case "f", "friends":
            log.Print("list friend")
            rpeople, err := client.GetPerson(ctx, &pb.PersonRequest{Uname: username})

            if err != nil {
                log.Fatalf("error retrieving user: %v", err)
            }

            var id *pb.Person
            if len(rpeople.People) > 0 {
                id = rpeople.People[0]
            }  else {
                log.Fatal("no matching users found")
            }

            rpeople1, err := client.GetFriends(ctx, id)

            for idx, value := range rpeople1.People {
                log.Printf("%d.user %v is friends with user %v", idx, id, value)
            }
        case "U", "unfriend":
            log.Print("unfriend")
            rpeople, err := client.GetPerson(ctx, &pb.PersonRequest{Uname: username})

            if err != nil {
                log.Fatalf("error retrieving user: %v", err)
            }

            var id0 *pb.Person

            if len(rpeople.People) > 0 {
                id0 = rpeople.People[0]
            }  else {
                log.Fatal("no matching users found")
            }

            username = os.Args[3]

            rpeople, err = client.GetPerson(ctx, &pb.PersonRequest{Uname: username})

            if err != nil {
                log.Fatalf("error retrieving user: %v", err)
            }

            var id1 *pb.Person

            if len(rpeople.People) > 0 {
                id1 = rpeople.People[0]
            }  else {
                log.Fatal("no matching users found")
            }

            rfriendship, err := client.Unfriend(ctx, &pb.Friendship{Orig: id0, Friend: id1})

            if rfriendship.Status {
                log.Print("unfriended")
            } else {
                log.Print("failed to unfriend")
            }
        case "S", "posts":
            log.Print("get posts")

            var posts *pb.Posts
            log.Printf("here on 85 %d", len(os.Args))
            if len(os.Args) < 3 {
                resposts, err := client.GetAllPosts(ctx, &pb.Empty{})
                if err != nil {
                    log.Fatalf("error retrieving posts: %v", err)
                }

                posts = resposts
            } else {
                rpeople, err := client.GetPerson(ctx, &pb.PersonRequest{Uname: username})
                if err != nil {
                    log.Fatalf("error retrieving user: %v", err)
                }

                var id *pb.Person
                if len(rpeople.People) > 0 {
                    id = rpeople.People[0]
                }  else {
                    log.Fatal("no matching users found")
                }

                resposts, err := client.GetPosts(ctx, id)
                if err != nil {
                    log.Fatalf("error retrieving posts: %v", err)
                }

                posts = resposts
            }

            for idx, value := range posts.Posts {
                log.Printf("idx: %d, post: %v\n", idx, value)
            }
        case "p", "people":
            log.Print("get all users")
            rpeople, err := client.GetPeople(ctx, &pb.Empty{})

            if err != nil {
                log.Fatal("could not access server")
            }

            for idx, value := range rpeople.People {
                log.Printf("idx: %d, value: %v\n", idx, value)
            }
        case "R", "person":
            log.Print("Get all similar users")
            rpeople, err := client.GetPerson(ctx, &pb.PersonRequest{Uname: username})

            if err != nil {
                log.Fatal("could not access server")
            }

            for idx, value := range rpeople.People {
                log.Printf("idx: %d, value: %v\n", idx, value)
            }
        case "r", "register":
            rperson, err := client.RegisterPerson(ctx, &pb.RegisterRequest{Uname: username})

            if err != nil {
                log.Fatalf("could not register user: %v", err)
            }

            log.Printf("registration success! %s, %s", rperson.GetUname(), rperson.GetId())
        default:
            log.Printf("no such command: %s", cmd)
    }
}
