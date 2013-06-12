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
	"os"
	"strings"

	"github.com/nf/twilio"
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

	twilio.Handle("/", handler)
}

func handler(c twilio.Context) {
	o := c.IntValue("Offset")
	if o < 0 || o >= len(digits) {
		c.Hangup()
		return
	}
	d := ""
	if o == 0 {
		d = "3 point "
	}
	d += digits[o]
	c.Responsef(`
		<Say>%v</Say>
		<Redirect method="GET">http://pi-phone.appspot.com/?Offset=%v</Redirect>
	`, d, o+1)
}
