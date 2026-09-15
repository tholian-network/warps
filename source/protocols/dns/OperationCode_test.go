package dns

import "testing"

func TestOperationCode(t *testing.T) {

	t.Run("String(Query)", func(t *testing.T) {
		if OperationCode(OperationCodeQuery).String() != "Query" {
			t.Errorf("Expected '%s' but got '%s'", "Query", OperationCode(OperationCodeQuery).String())
		}
	})

	t.Run("String(Status)", func(t *testing.T) {
		if OperationCode(OperationCodeStatus).String() != "Status" {
			t.Errorf("Expected '%s' but got '%s'", "Status", OperationCode(OperationCodeStatus).String())
		}
	})

	t.Run("String(Notify)", func(t *testing.T) {
		if OperationCode(OperationCodeNotify).String() != "Notify" {
			t.Errorf("Expected '%s' but got '%s'", "Notify", OperationCode(OperationCodeNotify).String())
		}
	})

	t.Run("String(Update)", func(t *testing.T) {
		if OperationCode(OperationCodeUpdate).String() != "Update" {
			t.Errorf("Expected '%s' but got '%s'", "Update", OperationCode(OperationCodeUpdate).String())
		}
	})

	t.Run("String(StatefulOperation)", func(t *testing.T) {
		if OperationCode(OperationCodeStatefulOperation).String() != "Stateful Operation" {
			t.Errorf("Expected '%s' but got '%s'", "Stateful Operation", OperationCode(OperationCodeStatefulOperation).String())
		}
	})

	t.Run("String(unknown)", func(t *testing.T) {
		if OperationCode(15).String() != "" {
			t.Errorf("Expected '%s' but got '%s'", "", OperationCode(15).String())
		}
	})

}
