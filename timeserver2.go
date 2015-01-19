// CSS 490
// Magdalena Grzmiel
// Assignments #1
// Copyright 2015 Magdalena Grzmiel
// This is a simple example of http server which display the actuall time
// after receiving following request: /time.
// Otherwise, it displayes information "These are not the URLs you're looking for"
// and sets the status code to 404 - Not Found.

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
	"os/exec"
	"bytes"
)


var sessions map[string]string

// logIn is a function which read the user name from request and display that on website
func logIn(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "text/plain")
	name := req.URL.Query().Get("name")
	id := generateUniqueId ()
	sessions[id] = name
	cookie := http.Cookie{Name : "uuid", Value : id, Path: "/"}
	http.SetCookie(res, &cookie)
	fmt.Fprintf(
		res,
		name + id,
	)
}

func generateUniqueId() string{
	cmd := exec.Command("/usr/bin/uuidgen")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	return out.String()
}


// loginForm is a function which displays the login form
func loginForm(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(
		res,
		`<html>
		<body>
		<form action="login">
		  What is your name, Earthling?
		  <input type="text" name="name" size="50">
		  <input type="submit">
		</form>
		</p>
		</body>
		</html>`,
	)
}

// getTime is a function which display the time on the webside
func getTime(res http.ResponseWriter, req *http.Request) {
	now := time.Now().Format("3:04:05 PM")
	res.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(
		res,
		`<doctype html>
        <html>
		<head>
		<style>
		p {font-size: xx-large}
		span.time {color: red}
		</style>
		</head>
		<body>
		<p>The time is now <span class="time">`+now+`</span>.</p>
		</body>
		</html>`,
	)
}

// unknownRoute is a function which display "These are not the URLs you're looking for"
// It also sets the status code to 404
func unknownRoute(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusNotFound)
	res.Header().Set("Content-Type", "html")
	fmt.Fprintf(
		res,
		`<html>
		<body>
		<p>These are not the URLs you're looking for.</p>
		</body>
		</html>`,
	)
}

// main function
// This function is responsible for the whole program, it first read the flags from command
// line and based on the flag it either runs the server or prints the version of the program.
// If the user does not provide the port number, the default one is 8080.
// When hitting /time, the website will display the current time.
// For every other route the website will return 404 status code.
func main() {

	var port int
	var version bool
	sessions = make(map[string]string)

	// parse the flags
	flag.IntVar(&port, "port", 8080, "used port")
	flag.BoolVar(&version, "V", false, "version of the program")
	flag.Parse()

	// if user type -V, the V flag is set up to true
	if version {
		// display the information about the version
		fmt.Println("version 1.1")
		// otherwise run the server
	} else {
		portNr := strconv.Itoa(port)
		http.HandleFunc("/time", getTime)
		http.HandleFunc("/", unknownRoute)
		http.HandleFunc("/index.html", loginForm)
		http.HandleFunc("/login", logIn)
		err := http.ListenAndServe(":"+portNr, nil)
		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	}
}
