package handlers

import (
	"github.com/Askalag/go-around/store"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"path/filepath"
)

type Files struct {
	log *log.Logger
	store store.Storage
}

func NewFiles(l *log.Logger, s store.Storage) *Files {
	return &Files{l, s}
}

func (f *Files) ServeHTTP(rw http.ResponseWriter, req *http.Request)  {
	vars := mux.Vars(req)
	id := vars["id"]
	fn := vars["filename"]

	f.log.Println("Hanfle POST", "id :", id, "fileName :", fn)
	f.saveFile(id, fn, rw, req)
}

func (f *Files) saveFile(id, path string, rw http.ResponseWriter, req *http.Request) {
	f.log.Println("Save file ", "id: ", id, "path: ", path)

	fp := filepath.Join(id, path)
	err := f.store.Save(fp, req.Body)
	if err != nil {
		f.log.Println("Unable to save file", "error: ", err)
		http.Error(rw, "Unable to save file", http.StatusInternalServerError)
	}
}
