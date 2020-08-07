package logging

import (
	"net/http"
	"strconv"
)

/*LogWriter - a handler to get recent logs*/
func LogWriter(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()
	detailLevel, _ := strconv.Atoi(queryValues.Get("detail"))
	mLogger.WriteLogs(w, detailLevel)
}

/*MemLogWriter - a handler to get the recent memory logs */
func MemLogWriter(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()
	detailLevel, _ := strconv.Atoi(queryValues.Get("detail"))
	mMLogger.WriteLogs(w, detailLevel)
}
