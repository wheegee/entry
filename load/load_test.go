package load

import (
	"bytes"
	"os"
	"testing"
)

func TestStdout(t *testing.T) {
	old := os.Stdout // keep backup of the real stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	envSlice := []string{"DB_HOST=localhost", "DB_PORT=5432"}
	expectedOutput := "export DB_HOST=localhost\nexport DB_PORT=5432\n"

	Stdout(envSlice)

	w.Close()
	os.Stdout = old // restoring the real stdout
	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	if output != expectedOutput {
		t.Errorf("Expected output %q, got %q", expectedOutput, output)
	}
}
