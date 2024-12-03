package handlers

import (
	"567_final/llmservice"
	"encoding/json"
	"fmt"
	"net/http"
)

// postSimulation

type SimulationResponse struct {
	Results []SimulationResult `json:"results"`
}

type SimulationResult struct {
	Policy string            `json:"policy"`
	Result map[string]string `json:"result"`
}

func GetSimulationHandler(w http.ResponseWriter, r *http.Request) {

	// get current policies
	currPolicyList := GetPolicyFromDB(true)
	var currPolicies []string
	for _, policy := range currPolicyList {
		currPolicies = append(currPolicies, policy.PolicyDescription)
	}

	// get new policies
	newPolicyList := GetPolicyFromDB(false)
	var newPolicies []string
	for _, policy := range newPolicyList {
		newPolicies = append(newPolicies, policy.PolicyDescription)
	}

	// TODO: connect with LLM and get feedback
	var results []SimulationResult

	for _, newPolicy := range newPolicies {
		var result SimulationResult
		result.Policy = newPolicy
		result.Result = make(map[string]string)

		// Simulate the new policy
		responses, err := llmservice.SimulatePolicy(currPolicies, newPolicy)
		if err != nil {
			fmt.Printf("Error simulating policy: %v\n", err)
			return
		}

		// Display the simulation responses
		for role, response := range responses {
			result.Result[role] = response
		}
		results = append(results, result)
	}

	// Return the response as JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(results); err != nil {
		http.Error(w, "Failed to generate response", http.StatusInternalServerError)
		return
	}

}
