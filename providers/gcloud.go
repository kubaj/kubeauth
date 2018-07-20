package providers

import (
	"encoding/json"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
)

// GCloudAccount represents single Google Cloud account
type GCloudAccount struct {
	Account string `json:"account"`
	Status  string `json:"status"`
}

// GCloudProject represents Google Cloud project
type GCloudProject struct {
	Name      string `json:"name"`
	ProjectID string `json:"projectId"`
}

// GCloudProvider represents Google Cloud, there can be multiple providers in the future (AWS, Azure, ...)
type GCloudProvider struct {
	ProjectID string
}

// ReadAccounts fetches all available Google Cloud accounts
func (g *GCloudProvider) ReadAccounts() ([]string, error) {
	cmd := exec.Command("gcloud", "auth", "list", "--format=json")

	var accounts []GCloudAccount
	out, err := cmd.Output()
	if err != nil {
		return []string{}, err
	}

	err = json.Unmarshal(out, &accounts)
	if err != nil {
		return []string{}, err
	}

	items := make([]string, len(accounts))
	for i, account := range accounts {
		items[i] = account.Account
	}
	return items, err
}

// SelectAccount selects Google Cloud account with provided email address
func (g *GCloudProvider) SelectAccount(account string) error {
	return exec.Command("gcloud", "config", "set", "account", account).Run()
}

// ReadProjects fetches all available Google Cloud projects
func (g *GCloudProvider) ReadProjects() ([]string, error) {
	cmd := exec.Command("gcloud", "projects", "list", "--format=json")

	var projects []GCloudProject
	out, err := cmd.Output()
	if err != nil {
		return []string{}, err
	}

	err = json.Unmarshal(out, &projects)
	if err != nil {
		return []string{}, err
	}

	items := make([]string, len(projects))
	for i, project := range projects {
		items[i] = project.ProjectID
	}
	return items, err
}

// SelectProject selects provided project as an active
func (g *GCloudProvider) SelectProject(project string) error {
	g.ProjectID = project
	return exec.Command("gcloud", "config", "set", "project", project).Run()
}

// ReadClusters fetches all avaiable Google Cloud Kubernetes clusters
func (g *GCloudProvider) ReadClusters() ([]string, error) {
	cmd := exec.Command("gcloud", "container", "clusters", "list", "--format=value(name,zone)")

	out, err := cmd.Output()
	if err != nil {
		return []string{}, errors.Wrap(err, "Failed to fetch clusters")
	}

	return strings.Split(string(out), "\n"), nil
}

// SelectCluster fetches credentials of the cluster from Google Cloud. Arg cluster is in format 'cluster_name\tcluster_zone'
func (g *GCloudProvider) SelectCluster(cluster string) error {
	clusterID := strings.Split(cluster, "\t")[0]
	zone := strings.Split(cluster, "\t")[1]
	return exec.Command("gcloud", "container", "clusters", "get-credentials", clusterID, "--zone", zone, "--project", g.ProjectID).Run()
}

// ReadNamespaces attempts to find any namespaces associated to
// a given cluster and make them available for selection.
func (g *GCloudProvider) ReadNamespaces(cluster string) ([]string, error) {
	kubeConfig := KubeConfig{}
	contexts, err := kubeConfig.ReadContexts(cluster)
	return contexts, err
}

// SelectContext selects a particular kube context
func (g *GCloudProvider) SelectContext(context string) error {
	kubeConfig := KubeConfig{}
	return kubeConfig.SelectContext(context)
}
