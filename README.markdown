Quoted-Printable encoding for Golang
====================================

Quoted-Printable encoding
-------------------------

detail: [http://en.wikipedia.org/wiki/Quoted-printable](http://en.wikipedia.org/wiki/Quoted-printable)

Usage
-----

	package main

	import "fmt"
	import "github.com/googollee/qp.go"

	func main() {
		str := "If you believe that truth=3Dbeauty, then surely =\nmathematics is the most beautiful branch of philosophy."
		expect := "If you believe that truth=beauty, then surely mathematics is the most beautiful branch of philosophy."
		reader := goqp.NewDecoder(bytes.NewBufferString(str))
		got, _ := ioutil.ReadAll(reader)
		if string(got) != expect {
			fmt.Printf("Decode\n\texpect: %s\n\t   got: %s\n", expect, string(got))
		}
	}