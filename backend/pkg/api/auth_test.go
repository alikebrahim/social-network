package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
	
	"github.com/stretchr/testify/assert"
	"socialNetwork/pkg/domain/auth"
)

func TestUserRegistration(t *testing.T) {
	testCases := []struct {
		name           string
		email          string
		password       string
		expectedStatus int
		expectError    bool
		errorField     string
	}{
		{
			name:           "Valid Registration",
			email:          "test@example.com",
			password:       "password123",
			expectedStatus: http.StatusCreated,
			expectError:    false,
		},
		{
			name:           "Invalid Email",
			email:          "invalid",
			password:       "password123",
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
			errorField:     "email",
		},
		{
			name:           "Password Too Short",
			email:          "test@example.com",
			password:       "pass",
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
			errorField:     "password",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			registrationRequest := auth.RegistrationRequest{
				Email:        tc.email,
				Password:     tc.password,
				First_name:   "Test",
				Last_name:    "User",
				Profile_type: "public",
			}
			
			reqBody, _ := json.Marshal(registrationRequest)
			
			req := httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(string(reqBody)))
			req.Header.Set("Content-Type", "application/json")
			
			w := httptest.NewRecorder()
			
			assert.Equal(t, tc.expectedStatus, w.Code)
			
			if tc.expectError {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "errors")
				
				errors, ok := response["errors"].(map[string]interface{})
				assert.True(t, ok)
				assert.Contains(t, errors, tc.errorField)
			}
		})
	}
}