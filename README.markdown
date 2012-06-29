Some encoding implement for Go
==============================

- Quoted-Printable encoding
- Ignore reader

Quoted-Printable encoding
-------------------------

detail: [http://en.wikipedia.org/wiki/Quoted-printable](http://en.wikipedia.org/wiki/Quoted-printable)

Usage
-----

- Quoted-Printable encoding

		func main() {
		    str := "If you believe that truth=3Dbeauty, then surely =\nmathematics is the most beautiful branch of philosophy."
		    expect := "If you believe that truth=beauty, then surely mathematics is the most beautiful branch of philosophy."
		    reader := encodingex.NewDecoder(bytes.NewBufferString(str))
		    got, _ := ioutil.ReadAll(reader)
		    if string(got) != expect {
		        fmt.Printf("Decode\n\texpect: %s\n\t   got: %s\n", expect, string(got))
		    }
		}

- Ignore reader

		func main() {
			str := "\r\n1231\tjkl\rrewq\nerq"	
			expect := "1231\tjklrewqerq"	
			reader := NewIgnoreReader(bytes.NewBufferString(str), []byte("\r\n"))
			got, _ := ioutil.ReadAll(reader)
			if string(got) != expect {
				fmt.Printff("Expect: %s, got: %s", expect, string(got))
			}
		}
