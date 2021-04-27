// goEnv.go

package goEnv

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
)

// GetGoEnv: run the shell command and return the go environment variable.
// This is better than using 'os.Getenv', as this also checks for the
// correct presence of golang .
func GetGoEnv(env string) (string, error) {
	var (
		stdout,
		stderr bytes.Buffer
		err    error
		reSpce = regexp.MustCompile(`\s$`)
	)
	cmd := []string{"go", "env"}
	cmd = append(cmd, env)
	execCmd := exec.Command(cmd[0], cmd[1:]...)
	execCmd.Stdout = &stdout
	execCmd.Stderr = &stderr
	err = execCmd.Run()
	if err != nil {
		return fmt.Sprintf(
			"stdout:\n%s\nstderr:\n%s\n",
			stdout.String(),
			stderr.String()), err
	}
	out := reSpce.ReplaceAllString(stdout.String(), "")
	if len(out) > 0 {
		return out, nil
	}
	return "", fmt.Errorf("There is no value for: %s\n", env)
}
