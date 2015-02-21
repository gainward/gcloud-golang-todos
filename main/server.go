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

// This package implements a simple HTTP server providing a REST API to a todo
// handler.
//
// It provides six methods:
//
// 	GET	/todos		Retrieves all the todos.
// 	POST	/todos		Creates a new todo given a title.
//	DELETE	/todos		Deletes all completed todos.
// 	GET	/todos/{todoID}	Retrieves the todo with the given id.
// 	PUT	/todos/{todoID}	Updates the todo with the given id.
//	DELETE	/todos/{todoID}	Deletes the todo with the given id.
//
// Every method below gives more information about every API call, its
// parameters, and its results.
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/GoogleCloudPlatform/gcloud-golang-todos/todo"
	"github.com/gorilla/mux"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"
	"google.golang.org/cloud"
	"google.golang.org/cloud/datastore"
)

const PathPrefix = "/api/todos"
const SlashedPathPrefix = PathPrefix + "/"

func init() {
	r := mux.NewRouter()
	r.HandleFunc(PathPrefix,
		errorHandler(ListTodos)).Methods("GET")
	r.HandleFunc(PathPrefix,
		errorHandler(NewTodo)).Methods("POST")
	r.HandleFunc(PathPrefix,
		errorHandler(DeleteCompletedTodos)).Methods("DELETE")
	r.HandleFunc(SlashedPathPrefix+"{id}",
		errorHandler(GetTodo)).Methods("GET")
	r.HandleFunc(SlashedPathPrefix+"{id}",
		errorHandler(UpdateTodo)).Methods("PUT")
	r.HandleFunc(SlashedPathPrefix+"{id}",
		errorHandler(DeleteTodo)).Methods("DELETE")
	http.Handle("/", r)
	http.HandleFunc("/api", IsApiEnabled)
}

func newCloudContext(c context.Context) context.Context {
	hc := &http.Client{
		Transport: &oauth2.Transport{
			Source: google.AppEngineTokenSource(c,
				datastore.ScopeDatastore,
				datastore.ScopeUserEmail),
			Base: &urlfetch.Transport{Context: c},
		},
	}
	return cloud.NewContext(appengine.AppID(c), hc)
}

// badRequest is handled by setting the status code in the reply to StatusBadRequest.
type badRequest struct{ error }

// notFound is handled by setting the status code in the reply to StatusNotFound.
type notFound struct{ error }

// IsApiEnabled writes an HTTP 200 indicating that the TodoMvc API is enabled for this app.
func IsApiEnabled(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// errorHandler wraps a function returning an error by handling the error and returning a http.Handler.
// If the error is of the one of the types defined above, it is handled as described for every type.
// If the error is of another type, it is considered as an internal error and its message is logged.
func errorHandler(f func(w http.ResponseWriter, r *http.Request, c, cloudC context.Context) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c := appengine.NewContext(r)
		cc := newCloudContext(c)
		err := f(w, r, c, cc)
		if err == nil {
			return
		}
		switch err.(type) {
		case badRequest:
			http.Error(w, err.Error(), http.StatusBadRequest)
		case notFound:
			http.Error(w, "todo not found", http.StatusNotFound)
		default:
			log.Errorf(c, "internal exception %v", err)
			http.Error(w, "oops", http.StatusInternalServerError)
		}
	}
}

// ListTodos handles GET requests on /todos.
// It requires no parameters and returns a list of todos.
//
// Example:
//
//   req: GET /todos/
//   res: 200 [
//          {"id": 1, "title": "Learn Go", "completed": false},
//          {"id": 2, "title": "Buy bread", "completed": true}
//        ]
func ListTodos(w http.ResponseWriter, r *http.Request, c, cc context.Context) error {
	res, err := todo.All(cc)
	if err != nil {
		log.Errorf(c, "ListTodos: %v", err)
		return err
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(res)
}

// NewTodo handles POST requests on /todos.
// The request body must contain a JSON object with a Title field.
// The status code of the response is used to indicate any error.
//
// Examples:
//
//   req: POST /todos/ {"title": ""}
//   res: 400 empty title
//
//   req: POST /todos/ {"title": "Buy bread"}
//   res: 201
func NewTodo(w http.ResponseWriter, r *http.Request, c, cc context.Context) error {
	req := struct{ Title string }{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return badRequest{err}
	}
	t, err := todo.NewTodo(req.Title)
	if err != nil {
		return badRequest{err}
	}
	log.Infof(c, "Saving new todo: %v", t)
	if err = t.Save(cc); err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	return json.NewEncoder(w).Encode(t)
}

// parseID obtains the id variable from the given request url,
// parses the obtained text and returns the result.
func parseID(r *http.Request) (int64, error) {
	txt, ok := mux.Vars(r)["id"]
	if !ok {
		return 0, fmt.Errorf("todo id not found")
	}
	return strconv.ParseInt(txt, 10, 0)
}

// GetTodo handles GET requsts to /todos/{todoID}.
// There's no parameters and it returns a JSON encoded todo.
//
// Examples:
//
//   req: GET /todos/1
//   res: 200 {"id": 1, "title": "Buy bread", "completed": true}
//
//   req: GET /todos/42
//   res: 404 todo not found
func GetTodo(w http.ResponseWriter, r *http.Request, c, cc context.Context) error {
	id, err := parseID(r)
	if err != nil {
		return badRequest{err}
	}
	t, ok := todo.Get(cc, id)
	if !ok {
		return notFound{}
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(t)
}

// UpdateTodo handles PUT requests to /todos/{todoID}.
// The request body must contain a JSON encoded todo.
// The id property is optional, but if it is included, it must match
// with the request path.
//
// Example:
//
//   req: PUT /todos/1 {"id": 1, "title": "Learn Go", "completed": true}
//   res: 200
//
//   req: PUT /todos/2 {"id": 1, "title": "Learn Go", "completed": true}
//   res: 400 inconsistent todo IDs
func UpdateTodo(w http.ResponseWriter, r *http.Request, c, cc context.Context) error {
	id, err := parseID(r)
	if err != nil {
		return badRequest{err}
	}
	var t todo.Todo
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		return badRequest{err}
	}
	if t.ID == 0 {
		t.ID = id
	}
	if t.ID != id {
		return badRequest{fmt.Errorf("inconsistent todo IDs")}
	}
	if _, ok := todo.Get(cc, id); !ok {
		log.Infof(c, "Unable to find todo: %v", t)
		return notFound{}
	}
	if err = t.Save(cc); err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(t)
}

// DeleteTodo handles DELETE requests to /todos/{todoID}.
// Returns a badRequest error if the ID cannot be parsed, and notFound if
// no corresponding todo can be found.
func DeleteTodo(w http.ResponseWriter, r *http.Request, c, cc context.Context) error {
	id, err := parseID(r)
	if err != nil {
		return badRequest{err}
	}
	log.Infof(c, "Trying to delete id %v", id)
	if ok := todo.Delete(cc, id); !ok {
		return notFound{}
	}
	w.WriteHeader(http.StatusNoContent)
	return nil
}

// DeleteCompletedTodos handles DELETE requests to /todos.
// It attempts to delete all todos which have been marked completed and returns
// an error if one occurred.
func DeleteCompletedTodos(w http.ResponseWriter, r *http.Request, c, cc context.Context) error {
	if err := todo.DeleteCompleted(cc); err != nil {
		return err
	}
	w.WriteHeader(http.StatusNoContent)
	return nil
}
