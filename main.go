package main

import (
	"log"

	"github.com/kubaj/kubeauth/providers"
	"gopkg.in/AlecAivazis/survey.v1"
)

func main() {
	gcloud := providers.GCloudProvider{}
	accs, err := gcloud.ReadAccounts()
	if err != nil {
		log.Fatalln(err)
	}

	account := ""
	survey.AskOne(
		&survey.Select{
			Message: "Choose an account:",
			Options: accs,
		},
		&account,
		nil,
	)
	gcloud.SelectAccount(account)

	projects, err := gcloud.ReadProjects()
	if err != nil {
		log.Fatalln(err)
	}

	project := ""
	survey.AskOne(
		&survey.Select{
			Message: "Choose a project:",
			Options: projects,
		},
		&project,
		nil,
	)

	err = gcloud.SelectProject(project)
	if err != nil {
		log.Fatalln(err)
	}

	clusters, err := gcloud.ReadClusters()
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(clusters)

	cluster := ""
	survey.AskOne(
		&survey.Select{
			Message: "Choose a cluster:",
			Options: clusters,
		},
		&cluster,
		nil,
	)

	err = gcloud.SelectCluster(cluster)
	if err != nil {
		log.Fatalln(err)
	}
}
