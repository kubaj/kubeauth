package providers

import (
	"encoding/json"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
)

// KubeConfig ...
type KubeConfig struct{}

// KubeContext define
type KubeContext struct {
	Kind        string `json:"kind"`
	APIVersion  string `json:"apiVersion"`
	Preferences struct {
	} `json:"preferences"`
	Clusters []struct {
		Name    string `json:"name"`
		Cluster struct {
			Server                string `json:"server"`
			InsecureSkipTLSVerify bool   `json:"insecure-skip-tls-verify"`
		} `json:"cluster"`
	} `json:"clusters"`
	Users []struct {
		Name string `json:"name"`
		User struct {
			ClientCertificateData string `json:"client-certificate-data"`
			ClientKeyData         string `json:"client-key-data"`
		} `json:"user"`
	} `json:"users"`
	Contexts []struct {
		Name    string `json:"name"`
		Context struct {
			Cluster string `json:"cluster"`
			User    string `json:"user"`
		} `json:"context"`
	} `json:"contexts"`
	CurrentContext string `json:"current-context"`
}

// ReadContexts fetches all available contexts in kube config.
func (k *KubeConfig) ReadContexts(cluster string) ([]string, error) {
	var contexts KubeContext
	kubectlCmd := exec.Command("kubectl", "config", "view", "-o", "json")
	out, err := kubectlCmd.Output()
	if err != nil {
		return []string{}, err
	}
	err = json.Unmarshal(out, &contexts)
	if err != nil {
		return []string{}, err
	}

	currentContext, err := k.ReadCurrentContext()
	if err != nil {
		return []string{}, err
	}

	items := make([]string, 0)
	for _, context := range contexts.Contexts {
		if context.Context.Cluster == currentContext && context.Name != currentContext {
			items = append(items, strings.TrimSpace(context.Name))
		}
	}
	return items, err
}

// SelectContext sets a selected context
func (k *KubeConfig) SelectContext(context string) error {
	return exec.Command("kubectl", "config", "use-context", context).Run()
}

// ReadCurrentContext returns the current value for the kube config context
func (k *KubeConfig) ReadCurrentContext() (string, error) {
	out, err := exec.Command("kubectl", "config", "current-context").Output()
	if err != nil {
		return "", errors.Wrap(err, "Failed to fetch current context")
	}
	return strings.TrimSuffix(string(out), "\n"), nil
}
