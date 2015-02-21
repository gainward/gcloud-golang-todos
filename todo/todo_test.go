// Copyright 2014 Google Inc. All rights reserved.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to writing, software distributed
// under the License is distributed on a "AS IS" BASIS, WITHOUT WARRANTIES OR
// CONDITIONS OF ANY KIND, either express or implied.
//
// See the License for the specific language governing permissions and
// limitations under the License.

package todo

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/cloud/datastore"
)

const jsonFile = "key.json"

var CTX cloud.Context

func newClient() (*http.Client, error) {
	jsonKey, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		return nil, err
	}
	conf, err := google.JWTConfigFromJSON(
		jsonKey, datastore.ScopeDatastore, datastore.ScopeUserEmail)
	if err != nil {
		return nil, err
	}
	return conf.Client(oauth2.NoContext), nil
}

func TestMain(m *testing.M) {
	hc, err := newClient()
	if err != nil {
		log.Fatalf("Could not create http client: %v", err)

	}
	CTX = cloud.NewContext("gcloud-golang-todos", hc)
	os.Exit(m.Run())
}

func TestSave(t *testing.T) {

}
