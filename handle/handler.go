package handle

import (
	"encoding/json"
	"net/http"
	"tuning_db/configuration"
	"tuning_db/tuning"
)

func HandleRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var config configuration.Configuration

	err := json.NewDecoder(r.Body).Decode(&config)
	if err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	err = tuning.ConfigureLocalTuningDatabase(config)
	if err != nil {
		
		http.Error(w,  "Error creating tuning configuration file", http.StatusInternalServerError)
		return
	}

	successMessage := map[string]string{
		"message": "Tuning configuration file created successfully",
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(successMessage)
}
