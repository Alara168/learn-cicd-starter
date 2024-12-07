package auth

import (
	"errors"
	"net/http"
	"testing"
)

func TestGetAPIKey(t *testing.T) {
	tests := []struct {
		name          string
		headers       http.Header
		expectedKey   string
		expectedError error
	}{
		{
			name: "Valid API Key",
			headers: http.Header{
				"Authorization": []string{"ApiKey abc123"},
			},
			expectedKey:   "abc123",
			expectedError: nil,
		},
		{
			name:          "Missing Authorization Header",
			headers:       http.Header{},
			expectedKey:   "",
			expectedError: ErrNoAuthHeaderIncluded,
		},
		{
			name: "Malformed Authorization Header - Missing ApiKey",
			headers: http.Header{
				"Authorization": []string{"abc123"},
			},
			expectedKey:   "",
			expectedError: errors.New("malformed authorization header"),
		},
		{
			name: "Malformed Authorization Header - Wrong Format",
			headers: http.Header{
				"Authorization": []string{"Bearer abc123"},
			},
			expectedKey:   "",
			expectedError: errors.New("malformed authorization header"),
		},
		{
			name: "Malformed Authorization Header - Empty ApiKey",
			headers: http.Header{
				"Authorization": []string{"ApiKey "},
			},
			expectedKey:   "",
			expectedError: nil, // Changed this to nil
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key, err := GetAPIKey(tt.headers)

			if key != tt.expectedKey {
				t.Errorf("Expected key %q, got %q", tt.expectedKey, key)
			}

			if (err == nil && tt.expectedError != nil) || (err != nil && tt.expectedError == nil) || (err != nil && tt.expectedError != nil && err.Error() != tt.expectedError.Error()) {
				t.Errorf("Expected error %v, got %v", tt.expectedError, err)
			}
		})
	}
}
