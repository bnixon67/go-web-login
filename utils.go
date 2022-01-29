/*
   Copyright 2022 Bill Nixon

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
package weblogin

import (
	"crypto/rand"
	"encoding/base64"
)

// GenerateRandomString returns n bytes encoded in URL friendly base64.
func GenerateRandomString(n int) (string, error) {
	// buffer to store n bytes
	b := make([]byte, n)

	// get b random bytes
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	// convert to URL friendly base64
	return base64.URLEncoding.EncodeToString(b), err
}
