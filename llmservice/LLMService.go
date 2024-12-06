// LlmService.go
package llmservice

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

// getGeminiClient initializes and returns a Gemini API client.
// It fetches the API key from the environment variable "API_KEY" and uses it to authenticate requests.
// Returns the initialized client, a context to be used for API requests, and any error that occurred during initialization.
func getGeminiClient() (*genai.Client, context.Context, error) {
	ctx := context.Background()
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return nil, nil, fmt.Errorf("API key must be set in the environment variable API_KEY")
	}

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, nil, err
	}

	return client, ctx, nil
}

// Extracts the concatenated content from the Gemini API response
func extractResponseText(resp *genai.GenerateContentResponse) string {
	var buffer bytes.Buffer
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				_, err := fmt.Fprintf(&buffer, "%s", part)
				if err != nil {
					log.Printf("Error writing content to buffer: %v", err)
				}

			}
		}
	}

	return buffer.String() // Convert the buffer contents to a string
}

// Convert policy into text
func parsePolicyResponse(response string) (string, string) {
	var newPolicyNeeded, proposedPolicy string

	// Regular expression to match "New Policy Needed: Yes" or "No"
	needPolicyRegex := regexp.MustCompile(`(?mi)^New Policy Needed:\s*(Yes|No)`)
	needPolicyMatches := needPolicyRegex.FindStringSubmatch(response)
	if len(needPolicyMatches) > 1 {
		newPolicyNeeded = needPolicyMatches[1] // Capture "Yes" or "No"
	}

	// Regular expression to match "Proposed Policy" followed by its detailed content
	policyRegex := regexp.MustCompile(`(?mi)^Proposed Policy:\s*(.*)`)
	policyMatches := policyRegex.FindStringSubmatch(response)
	if len(policyMatches) > 1 {
		proposedPolicy = policyMatches[1] // Capture everything
	}

	return newPolicyNeeded, proposedPolicy
}

// GenerateNewPolicy determines if a new policy is needed based on the current policies,
// the content of a post, and optionally, user comments.
// Inputs:
// - currentPolicies: An array of strings representing the existing policies.
// - postContent: A string representing the content of a post being considered for policy evaluation.
// - userComment: An optional string pointer representing user comments about the post.
// Outputs:
// - A string representing the newly generated policy if one is needed.
// - An error if any occurs during the process.
// The function first constructs a prompt to query whether a new policy is required and,
// if so, what that policy should entail. It then calls the Gemini API with this prompt
// and processes the response to determine if a new policy should be drafted.
func GenerateNewPolicy(currentPolicies []string, postContent string, userComment *string) (string, error) {
	client, ctx, err := getGeminiClient()
	if err != nil {
		return "", err
	}
	defer client.Close()

	// Combine current policies into a single string
	policiesText := ""
	for _, policy := range currentPolicies {
		policiesText += "- " + policy + "\n"
	}

	// Construct the prompt
	prompt := fmt.Sprintf(
		"Current Policies:\n%s\n\n"+
			"A user has raised a concern about the following post:\n'%s'\n\n", policiesText, postContent)

	if userComment != nil {
		prompt += fmt.Sprintf("User have following concerns regarding the post:\n'%s'\n\n", *userComment)
	}

	prompt += "Based on the above, please:\n" +
		"1. Determine if a new policy is needed to address the issue.\n" +
		"2. If needed, draft a new policy that effectively addresses the concern. You should only provide a concise description of the policy you drafted and no further discussion needed.\n" +
		"Respond in the following format:\n" +
		"New Policy Needed: [Yes/No]\n" +
		"Proposed Policy: [If applicable, provide the policy text]"

	// Call the Gemini API
	model := client.GenerativeModel("gemini-1.5-flash")
	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", err
	}

	// Process the response to extract the new policy
	responseText := extractResponseText(resp)
	// fmt.Println(responseText)
	// assuming the LLM follows the response format
	var newPolicyNeeded string
	var proposedPolicy string
	//fmt.Printf("%s", responseText)
	newPolicyNeeded, proposedPolicy = parsePolicyResponse(responseText)
	//fmt.Printf("%s", newPolicyNeeded)
	//fmt.Printf("%s", proposedPolicy)

	if newPolicyNeeded == "Yes" {
		return proposedPolicy, nil
	} else {
		return "", nil
	}
}

// SimulatePolicy simulates the effects of a new policy from different community roles perspectives
// and provides an overview of the policy impact and recommendations.
// Inputs:
// - currentPolicies: An array of strings representing the existing policies.
// - newPolicy: A string representing a newly proposed policy to be evaluated.
// Outputs:
//   - A map of string keys to string values where keys are roles ("Regular User", "Moderator", "Abuser", "Policy Overview")
//     and values are the responses generated by the Gemini API simulating the impact of the new policy.
//   - An error if any occurs during the process.
//
// The function defines different roles and their corresponding simulation prompts,
// sends these prompts to the Gemini API, and collects the responses.
// It then creates a comprehensive overview prompt incorporating all responses to provide
// a summary and recommendations about the new policy.
func SimulatePolicy(currentPolicies []string, newPolicy string) (map[string]string, error) {
	client, ctx, err := getGeminiClient()
	if err != nil {
		return nil, err
	}
	defer client.Close()

	// Combine current policies into a single string
	policiesText := ""
	for _, policy := range currentPolicies {
		policiesText += "- " + policy + "\n"
	}
	policiesText += "- " + newPolicy + "\n"

	// Define the roles and their prompts
	roles := map[string]string{
		"Regular User": "As a regular user, you see the following new policy:\n'%s'\n\nHow would you react to this policy? Provide your thoughts and feelings precisely in two sentences. Start your response with: The general user",
		"Moderator":    "As a moderator, you are tasked with enforcing the following new policy:\n'%s'\n\nWhat steps would you take to implement this policy effectively?  Provide your thoughts and feelings precisely in two sentences. Start your response with: The moderator",
		"Abuser":       "As someone who wants to test the limits of community guidelines, you encounter the following new policy:\n'%s'\n\nHow would you attempt to circumvent this policy?  Provide your thoughts and feelings precisely in two sentences. Start your response with: The abuser",
	}

	// Store the responses
	responses := make(map[string]string)

	for role, rolePrompt := range roles {
		// Construct the prompt for each role
		prompt := fmt.Sprintf(rolePrompt, newPolicy)

		// Optionally, you can include the current policies
		fullPrompt := fmt.Sprintf(
			"Current Policies:\n%s\n\n%s", policiesText, prompt)

		// Call the Gemini API
		model := client.GenerativeModel("gemini-1.5-flash")
		resp, err := model.GenerateContent(ctx, genai.Text(fullPrompt))
		if err != nil {
			return nil, err
		}

		responses[role] = extractResponseText(resp)
	}

	overviewPrompt := "Based on the following responses from different community roles:\n\n"
	for role, response := range responses {
		overviewPrompt += fmt.Sprintf("%s Response:\n%s\n\n", role, response)
	}
	overviewPrompt += "Provide a summary of the potential impacts of the new policy and offer recommendations for improvements. Provide your thoughts and feelings precisely in four sentences"
	model := client.GenerativeModel("gemini-1.5-flash")
	overviewResp, err := model.GenerateContent(ctx, genai.Text(overviewPrompt))
	if err != nil {
		return nil, err
	}

	responses["Policy Overview"] = extractResponseText(overviewResp)

	return responses, nil
}

////////////////////////////////////////////////////// EXAMPLE USAGE //////////////////////////////////////////////////////////////////

// func main() {
// 	currentPolicies := []string{}

// 	postContent := "I think we should exclude all newcomers from our group activities."

// 	userComment := "This post feels discriminatory and could foster a hostile environment."

// 	// Generate a new policy if needed
// 	newPolicy, err := GenerateNewPolicy(currentPolicies, postContent, &userComment)
// 	if err != nil {
// 		fmt.Printf("Error generating new policy: %v\n", err)
// 		return
// 	}

// 	if newPolicy != "" {
// 		fmt.Printf("New Policy Generated:\n%s\n\n", newPolicy)
// 	} else {
// 		fmt.Println("No new policy needed based on the analysis.")
// 		return
// 	}

// 	// Simulate the new policy
// 	responses, err := SimulatePolicy(currentPolicies, newPolicy)
// 	if err != nil {
// 		fmt.Printf("Error simulating policy: %v\n", err)
// 		return
// 	}

// 	// Display the simulation responses
// 	for role, response := range responses {
// 		fmt.Printf("%s:\n%s\n\n", role, response)
// 	}
// }
