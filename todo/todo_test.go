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
	"os"
	"testing"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/cloud"
	"google.golang.org/cloud/datastore"
)

const (
	envProjID     = "GCLOUD_GOLANG_TODOS_PROJECT_ID"
	envPrivateKey = "GCLOUD_GOLANG_TODOS_KEY"
)

var CTX context.Context

func Context(scopes ...string) context.Context {
	key, projID := os.Getenv(envPrivateKey), os.Getenv(envProjID)
	if key == "" || projID == "" {
		log.Fatalf("%v and %v must be set. See CONTRIBUTING.md.",
			envProjID, envPrivateKey)
	}
	jsonKey, err := ioutil.ReadFile(key)
	if err != nil {
		log.Fatalf("Cannot read the JSON key file, err: %v", err)
	}
	conf, err := google.JWTConfigFromJSON(jsonKey, scopes...)
	if err != nil {
		log.Fatal(err)
	}
	return cloud.NewContext(projID, conf.Client(oauth2.NoContext))
}

func TestMain(m *testing.M) {
	CTX = Context(datastore.ScopeDatastore, datastore.ScopeUserEmail)
	os.Exit(m.Run())
}

func TestSave(t *testing.T) {

}
