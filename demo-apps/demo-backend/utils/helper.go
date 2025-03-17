package utils

import "fmt"

// extractRoles extracts and validates the 'roles' claim from the provided claims map.
// It returns a slice of strings representing the roles, or an error if the extraction fails.
func ExtractRoles(claims map[string]interface{}) ([]string, error) {
	// Extract the 'roles' claim and assert it as a slice of interfaces
	rolesRaw, ok := claims["roles"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("roles not found or of incorrect type")
	}

	// Initialize the roles slice with the correct capacity
	roles := make([]string, 0, len(rolesRaw))

	// Iterate through the raw roles and convert them to strings
	for i, role := range rolesRaw {
		roleStr, ok := role.(string)
		if !ok {
			return nil, fmt.Errorf("role at position %d is not a string", i)
		}
		roles = append(roles, roleStr)
	}

	return roles, nil
}
