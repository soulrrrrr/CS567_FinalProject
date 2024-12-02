package handlers

import (
	"encoding/json"
	"net/http"
)

func GetPolicyHandler(w http.ResponseWriter, r *http.Request) {
	policyList := GetPolicyFromDB(true)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(policyList)
}
