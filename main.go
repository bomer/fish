package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"os/exec"
)

var (
	Name    string
	Players map[string]*Player //Map of all players, index is their cookie
)

type Player struct {
	name      string
	fishcount int
	gold      int
	rodlevel  int
	fishes    map[string]int
}

func (p *Player) addFish(fish string) {
	p.fishcount++
	p.fishes[fish]++
}

func (p *Player) init() {
	p.fishes = map[string]int{"MEGA": 0, "OK": 0, "LAME": 0}
}

func getUuid() string {
	out, err := exec.Command("uuidgen").Output()
	if err != nil {
		fmt.Println(err.Error())
	}
	o := fmt.Sprintf("%s", out)
	fmt.Printf("New user, gernerationg UUID - %s", out)
	return o
}

func handleCookie(w *http.ResponseWriter, r *http.Request) *Player {
	cookie, err := r.Cookie("fishuuid")
	var newuuid string
	if err != nil {
		fmt.Println(err.Error())
	}
	if cookie == nil {
		fmt.Printf("NO COOKIE!", cookie)
		newuuid = getUuid()
		cookie1 := &http.Cookie{Name: "fishuuid", Value: newuuid, HttpOnly: false}
		http.SetCookie(*w, cookie1)
	} else {
		newuuid = cookie.Value

	}
	if Players[newuuid] == nil {
		Players[newuuid] = &Player{}
		Players[newuuid].init()
	}
	//If the player has just added a name, set it.
	n := r.URL.Query().Get("name")
	fmt.Printf("%s", n)
	if len(n) > 0 {
		Players[newuuid].name = n
	}

	if !(len(Players[newuuid].name) > 0) {
		fmt.Fprint(*w, "<html><body><h1>New Player!</h1> <form action='/'> Please enter your Name : <input name='name'></input><button type='submit'>Start!</button> </form></body></html")
	} else {
		return Players[newuuid]
	}
	return Players[newuuid]

}

//Before Fish, shows the scoreboard + Instructions
func beforefish(w http.ResponseWriter, r *http.Request) {

	player := handleCookie(&w, r)
	if len(player.name) > 0 {
		fmt.Fprint(w, "<html><body>")
		fmt.Fprintf(w, "Hi there! %s", player.name)
		fmt.Fprint(w, "You are about to fish... BRACE YOURSELF! <br><br> <button type='button'><a href='/fish'>FISH</a></button>")
		fmt.Fprintf(w, "Caught:! %+v", player.fishes)
		fmt.Fprintf(w, "Total Caught: %d", player.fishcount)

		fmt.Fprintf(w, "<br><h2>Scoreboard</h2>")
		for _, v := range Players {
			if len(v.name) > 0 {
				fmt.Fprintf(w, "Player [%s] = %d\n <br>", v.name, v.fishcount)
			}
		}
		fmt.Fprintf(w, "")
	}

}
func fish(w http.ResponseWriter, r *http.Request) {

	//Todo - Show the map ?
	//Todo - Select a square ?
	player := handleCookie(&w, r)

	//Show result
	randumNum := rand.Intn(100)
	fmt.Fprint(w, "<html><body>")
	if randumNum > 80 {
		player.addFish("MEGA")
		fmt.Fprintf(w, "YOU RECIEVED THE MEGA FISH!!")
	} else if randumNum > 40 {
		player.addFish("OK")
		fmt.Fprintf(w, "YOU RECIEVED THE OK FISH!!")
	} else if randumNum >= 15 {
		player.addFish("LAME")
		fmt.Fprintf(w, "YOU RECIEVED THE LAME FISH!!")
	}
	fmt.Fprintf(w, "<br><br>")
	beforefish(w, r)

}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there. Hit '/data' on this server to get the latest results in JSON.")
}

// //Test code at beggining
// func manualFix() {
// 	fmt.Println("Please enter your name: ")
// 	fmt.Scanln(&name)

// 	fmt.Printf("Hi there %s\n", name)

// 	for {
// 		beforefish()
// 		fish()
// 	}
// }

func main() {
	fmt.Printf("Hello!\n")
	fmt.Println("Go Fish!! ( .  c )< ~~~     <><")

	Players = make(map[string]*Player)

	http.HandleFunc("/", beforefish)
	http.HandleFunc("/fish", fish)

	http.ListenAndServe(":8000", nil)

}
