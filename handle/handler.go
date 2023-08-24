package handle

import (
	"encoding/json"
	"fmt"
	"net/http"
	"tuning_db/configuration"
	"tuning_db/response"
	"tuning_db/tuning"
)

func HandleRequest(w http.ResponseWriter, r *http.Request) {
	var config configuration.Configuration
	err := json.NewDecoder(r.Body).Decode(&config)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	tuningConfig, err := tuning.CalculateMySQLTuning(config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	responseMessage := response.GenerateResponseMessage(tuningConfig)

	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(w, responseMessage)
}
