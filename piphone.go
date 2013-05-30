/*
Copyright 2013 Google Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package piphone

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var digits []string

func init() {
	f, err := os.Open("digits.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	s := bufio.NewScanner(f)
	for s.Scan() {
		d := strings.TrimPrefix(s.Text(), " 3.")
		d = strings.Join(strings.Split(d, ""), " ")
		digits = append(digits, d)
	}
	if err := s.Err(); err != nil {
		panic(err)
	}

	http.HandleFunc("/", handler)
}

func handler(w http.ResponseWriter, r *http.Request) {
	switch r.FormValue("Digits") {
	case "6":
		fmt.Fprint(w, gather)
		return
	case "5359":
		fmt.Fprint(w, dialConf)
		return
	}
	o, _ := strconv.Atoi(r.FormValue("offset"))
	if o < 0 || o >= len(digits) {
		fmt.Fprint(w, hangup)
		return
	}
	d := ""
	if o == 0 {
		d = "3 point "
	}
	d += digits[o]
	fmt.Fprintf(w, say, d, o+1)
}

const say = `<?xml version="1.0" encoding="UTF-8"?>
<Response>
    <Gather timeout="0">
        <Say>%v</Say>
    </Gather>
    <Redirect method="GET">http://pi-phone.appspot.com/?offset=%v</Redirect>
</Response>`

const hangup = `<?xml version="1.0" encoding="UTF-8"?>
<Response><Hangup/></Response>`

const gather = `<?xml version="1.0" encoding="UTF-8"?>
<Response>
    <Gather timeout="10"/>
    <Hangup/>
</Response>`

const dialConf = `<?xml version="1.0" encoding="UTF-8"?>
<Response>
    <Dial><Conference>Pi chat</Conference></Dial>
</Response>`
