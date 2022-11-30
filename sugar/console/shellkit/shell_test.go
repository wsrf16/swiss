package shellkit

import (
	"github.com/wsrf16/swiss/sugar/encoding/jsonkit"
	"testing"
)

func TestExecute(t *testing.T) {
	stdout, stderr, err := Execute("whoami")
	if err != nil {
		t.Error(err.Error())
	}
	t.Log(stdout)
	t.Log(stderr)

	strings := []string{"dir", "whoami"}
	batch, err := ExecuteBatch(strings)
	if err != nil {
		t.Error(err.Error())
	}
	t.Log(jsonkit.Marshal(batch))
}
