package jsonapi

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/MFarkha/my-mailinglist-microservice/mdb"
)

func setJsonHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
}

func fromJson[T any](body io.Reader, target T) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(body)
	json.Unmarshal(buf.Bytes(), &target)
}

func returnJson[T any](w http.ResponseWriter, withData func() (T, error)) error {
	setJsonHeaders(w)
	data, serverErr := withData()
	if serverErr != nil {
		w.WriteHeader(500)
		serverErrJson, err := json.Marshal(&serverErr)
		if err != nil {
			log.Printf("error from json.marshal of serverErr: %v", err)
			return err
		}
		w.Write(serverErrJson)
		return nil
	}
	dataJson, err := json.Marshal(&data)
	if err != nil {
		log.Printf("error from json.marshal of data: %v", err)
		w.WriteHeader(500)
		return err
	}
	w.Write(dataJson)
	return nil
}

func returnErr(w http.ResponseWriter, err error, code int) {
	returnJson(w, func() (interface{}, error) {
		errorMessage := struct {
			Err string
		}{
			Err: err.Error(),
		}
		w.WriteHeader(code)
		return errorMessage, nil
	})
}

func CreateEmail(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			return
		}
		emailEntry := mdb.EmailEntry{}
		fromJson(req.Body, &emailEntry)

		if err := mdb.CreateEmailEntry(db, emailEntry.Email); err != nil {
			returnErr(w, err, 400)
			return
		}
		returnJson(w, func() (interface{}, error) {
			log.Printf("JSON CreateEmail: %v\n", emailEntry.Email)
			return mdb.GetEmailEntry(db, emailEntry.Email)
		})
	})
}

func GetEmail(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			return
		}
		emailEntry := mdb.EmailEntry{}
		fromJson(req.Body, &emailEntry)

		returnJson(w, func() (interface{}, error) {
			log.Printf("JSON GetEmail: %v\n", emailEntry.Email)
			return mdb.GetEmailEntry(db, emailEntry.Email)
		})
	})
}

func GetEmailBatch(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			return
		}
		queryOptions := mdb.GetEmailBatchQueryParams{}
		fromJson(req.Body, &queryOptions)
		if queryOptions.Count <= 0 || queryOptions.Page <= 0 {
			returnErr(w, errors.New("page and count should be set and >0"), 400)
			return
		}
		returnJson(w, func() (interface{}, error) {
			log.Printf("JSON GetEmailBatch: %v\n", queryOptions)
			return mdb.GetEmailBatch(db, queryOptions)
		})

	})
}

func UpdateEmail(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Method != "PUT" {
			return
		}
		emailEntry := mdb.EmailEntry{}
		fromJson(req.Body, &emailEntry)

		if err := mdb.UpdateEmailEntry(db, &emailEntry); err != nil {
			returnErr(w, err, 400)
			return
		}
		returnJson(w, func() (interface{}, error) {
			log.Printf("JSON UpdateEmail: %v\n", emailEntry.Email)
			return mdb.GetEmailEntry(db, emailEntry.Email)
		})
	})
}

func DeleteEmail(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			return
		}
		emailEntry := mdb.EmailEntry{}
		fromJson(req.Body, &emailEntry)

		if err := mdb.DeleteEmailEntry(db, emailEntry.Email); err != nil {
			returnErr(w, err, 400)
			return
		}
		returnJson(w, func() (interface{}, error) {
			log.Printf("JSON DeleteEmail: %v\n", emailEntry.Email)
			return mdb.GetEmailEntry(db, emailEntry.Email)
		})
	})
}

func Serve(db *sql.DB, bind string) {
	http.Handle("/email/create", CreateEmail(db))
	http.Handle("/email/get", GetEmail(db))
	http.Handle("/email/get_batch", GetEmailBatch(db))
	http.Handle("/email/update", UpdateEmail(db))
	http.Handle("/email/delete", DeleteEmail(db))
	log.Printf("JSON API server is listening on %s\n", bind)
	err := http.ListenAndServe(bind, nil)
	if err != nil {
		log.Fatalf("JSON server failure to bind: %v, error: %v", bind, err)
	}
}
