package templater

import "testing"

func TestTemplate(t *testing.T) {

	s, err := Template("templates/test.txt", map[string]string{"test": "success"})
	if err != nil {
		t.Error(err)
		return
	}
	if s != "{{asd}}success{}{{}success" {
		t.Errorf("Expected value to be success, but got %s instead", s)
	}
}