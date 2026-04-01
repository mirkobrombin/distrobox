package envfile

import (
	"bufio"
	"context"
	"os"
	"strings"

	"github.com/mirkobrombin/go-conf-builder/v2/pkg/source"
)

// Provider implements source.Provider for KEY=VALUE env files.
type Provider struct {
	path string
}

// New returns a new envfile Provider for the given path.
func New(path string) source.Provider {
	return &Provider{path: path}
}

func (p *Provider) Name() string {
	return "envfile:" + p.path
}

func (p *Provider) Load(_ context.Context) (map[string]any, error) {
	f, err := os.Open(p.path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	values := make(map[string]any)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		k, v, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		k = strings.TrimSpace(k)
		v = strings.TrimSpace(v)
		v = strings.Trim(v, `"'`)
		values[k] = v
	}
	return values, scanner.Err()
}
