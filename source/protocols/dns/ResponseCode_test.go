package dns

import "testing"

func TestResponseCode(t *testing.T) {

	t.Run("String(NoError)", func(t *testing.T) {
		if ResponseCode(ResponseCodeNoError).String() != "No Error" {
			t.Errorf("Expected '%s' but got '%s'", "No Error", ResponseCode(ResponseCodeNoError).String())
		}
	})

	t.Run("String(FormatError)", func(t *testing.T) {
		if ResponseCode(ResponseCodeFormatError).String() != "Format Error" {
			t.Errorf("Expected '%s' but got '%s'", "Format Error", ResponseCode(ResponseCodeFormatError).String())
		}
	})

	t.Run("String(ServerFailure)", func(t *testing.T) {
		if ResponseCode(ResponseCodeServerFailure).String() != "Server Failure" {
			t.Errorf("Expected '%s' but got '%s'", "Server Failure", ResponseCode(ResponseCodeServerFailure).String())
		}
	})

	t.Run("String(NonExistDomain)", func(t *testing.T) {
		if ResponseCode(ResponseCodeNonExistDomain).String() != "Non-Existent Domain" {
			t.Errorf("Expected '%s' but got '%s'", "Non-Existent Domain", ResponseCode(ResponseCodeNonExistDomain).String())
		}
	})

	t.Run("String(NotImplemented)", func(t *testing.T) {
		if ResponseCode(ResponseCodeNotImplemented).String() != "Not Implemented" {
			t.Errorf("Expected '%s' but got '%s'", "Not Implemented", ResponseCode(ResponseCodeNotImplemented).String())
		}
	})

	t.Run("String(QueryRefused)", func(t *testing.T) {
		if ResponseCode(ResponseCodeQueryRefused).String() != "Query Refused" {
			t.Errorf("Expected '%s' but got '%s'", "Query Refused", ResponseCode(ResponseCodeQueryRefused).String())
		}
	})

	t.Run("String(BadCookie)", func(t *testing.T) {
		if ResponseCode(ResponseCodeBadCookie).String() != "Bad Cookie" {
			t.Errorf("Expected '%s' but got '%s'", "Bad Cookie", ResponseCode(ResponseCodeBadCookie).String())
		}
	})

	t.Run("String(unknown)", func(t *testing.T) {
		if ResponseCode(100).String() != "" {
			t.Errorf("Expected '%s' but got '%s'", "", ResponseCode(100).String())
		}
	})

}
