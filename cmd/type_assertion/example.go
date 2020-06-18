package main

import (
	"fmt"
	"time"
)

type User struct {
	username string
	dob      time.Time
}

func (u *User) isFriend(u2 User) bool { return true }

type Engage interface {
	Friend(interface{}) error
	Unfriend(interface{}) error
	DM(User, string) error
}

func (u *User) Friend(user interface{}) error { return nil }

func (u *User) Unfriend(user interface{}) error { return nil }
func (u *User) DM(u2 User, msg string) error    { return nil }

type Frenemy struct {
	username    string
	hatredLevel int
}

func main() {
	a := User{username: "ayy"}
	b := User{username: "bee"}
	fe := Frenemy{username: "theWooorst"}

	fmt.Println("a friends with b?", a.isFriend(b))

	// but this happens all the time when intaking user data
	// we aren't positive what it'll come in as, so we call it "interface{}" and typecast to what it should be
	// this will wreck you -- always check if !ok
	theWorst, ok := fe.(User)

	if !ok {
		fmt.Println("just give up on this friendship")
		return
	}

	if !theWorst.isFriend(a) {
		// suggest a weird coffee catch up date
	}

}
