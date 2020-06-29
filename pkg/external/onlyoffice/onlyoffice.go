package onlyoffice

import "net/http"

func Edit() {

}

func View() {

}

func Convert() {

}

func Download() {

}

func Upload(rw http.ResponseWriter, req *http.Request) {

}

func Track(rw http.ResponseWriter, req *http.Request) {

}

func processSave(downloadUri, body, fileName, userAddress string, respw http.ResponseWriter) {

}

func processForceSave(downloadUri, body, fileName, userAddress string, respw http.ResponseWriter) {

}

func checkJwt(req *http.Request) (body interface{}, fileName, userAddress string) {

}
