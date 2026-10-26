package lookup

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/burmilla/os/pkg/libcompose/config"
)

// EnvfileLookup is a structure that implements the project.EnvironmentLookup interface.
// It holds the path of the file where to lookup environment values.
type EnvfileLookup struct {
	Path string
}

// Lookup creates a string slice of string containing a "docker-friendly" environment string
// in the form of 'key=value'. It gets environment values using a '.env' file in the specified
// path.
func (l *EnvfileLookup) Lookup(key, serviceName string, config *config.ServiceConfig) []string {
	envs, err := parseEnvFile(l.Path)
	if err != nil {
		return []string{}
	}
	for _, env := range envs {
		e := strings.Split(env, "=")
		if e[0] == key {
			return []string{env}
		}
	}
	return []string{}
}

// parseEnvFile reads a file with environment variables enumerated by lines
// (taken from github.com/docker/docker/runconfig/opts, see docker/opts.go).
//
// Environment variable names are only checked for white space, which is what
// docker does; anything else is passed through as is.
func parseEnvFile(filename string) ([]string, error) {
	const whiteSpaces = " \t"

	fh, err := os.Open(filename)
	if err != nil {
		return []string{}, err
	}
	defer fh.Close()

	lines := []string{}
	scanner := bufio.NewScanner(fh)
	for scanner.Scan() {
		// trim the line from all leading whitespace first
		line := strings.TrimLeft(scanner.Text(), whiteSpaces)
		// line is not empty, and not starting with '#'
		if len(line) > 0 && !strings.HasPrefix(line, "#") {
			data := strings.SplitN(line, "=", 2)

			// trim the front of a variable, but nothing else
			variable := strings.TrimLeft(data[0], whiteSpaces)
			if strings.ContainsAny(variable, whiteSpaces) {
				return []string{}, fmt.Errorf("poorly formatted environment: variable '%s' has white spaces", variable)
			}

			if len(data) > 1 {
				// pass the value through, no trimming
				lines = append(lines, fmt.Sprintf("%s=%s", variable, data[1]))
			} else {
				// if only a pass-through variable is given, clean it up.
				lines = append(lines, fmt.Sprintf("%s=%s", strings.TrimSpace(line), os.Getenv(line)))
			}
		}
	}
	return lines, scanner.Err()
}
