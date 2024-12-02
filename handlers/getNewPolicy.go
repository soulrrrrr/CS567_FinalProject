package handlers

import (
	"encoding/json"
	"net/http"
)

func GetNewPolicyHandler(w http.ResponseWriter, r *http.Request) {
	policyList := GetPolicyFromDB(false)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(policyList)
}
