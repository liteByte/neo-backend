package templater

import "testing"

func TestTemplate(t *testing.T) {

	s, err := Template("templates/test.txt", map[string]string{"test": "success", "{test": "asd"})
	if err != nil {
		t.Error(err)
		return
	}
	if s != "{{asd}}success{}{{}success" {
		t.Errorf("Expected value to be '{{asd}}success{}{{}success', but got '%s' instead", s)
	}
}