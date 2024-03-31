package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"vincent-h-lee/web-crawler/internal/util"
)

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	defer func() {
		if errD := recover(); errD != nil {
			fmt.Println("Exception occurred at ", util.RecoverExceptionDetails(util.FunctionName()), " and recovered in respondWithJSON function, Error Info: ", errD)
		}
	}()
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	defer func() {
		if errD := recover(); errD != nil {
			fmt.Println("Exception occurred at ", util.RecoverExceptionDetails(util.FunctionName()), " and recovered in respondWithError function, Error Info: ", errD)
		}
	}()
	respondWithJSON(w, code, map[string]interface{}{"success": false, "errors": map[string]string{"reason": message}})
}
