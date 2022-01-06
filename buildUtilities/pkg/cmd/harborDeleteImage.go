package cmd

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"galasa.dev/buildUtilities/pkg/galasayaml"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	project    string
	repository string
	tag        string

	harborDeleteCmd = &cobra.Command{
		Use:   "deleteimage",
		Short: "Remove a specified image from a Harbor docker registry",
		Long:  "Specify the project, repository and tag to remove the image from the harbor registry",

		Run: executeHarbor,
	}
)

func init() {
	harborDeleteCmd.PersistentFlags().StringVarP(&project, "project", "p", "", "Which project to be interacting with")
	harborDeleteCmd.PersistentFlags().StringVarP(&repository, "repository", "r", "", "Which repository to be interacting with")
	harborDeleteCmd.PersistentFlags().StringVarP(&tag, "tag", "t", "", "Which tag to be interacting with")

	harborCmd.AddCommand(harborDeleteCmd)
}

func executeHarbor(cmd *cobra.Command, args []string) {
	if harborRepository == "" {
		panic("Please provide a Harbor endpoint URL using the --harbor flag")
	}
	if project == "" {
		panic("Please provide a project using the --project flag")
	}
	if repository == "" {
		panic("Please provide a Harbor repository using the --repository flag")
	}
	if tag == "" {
		panic("Please provide a repository tag using the --tag flag")
	}

	client := http.Client{
		Timeout: time.Second * 30,
	}
	req, err := http.NewRequest("DELETE", harborRepository+v2Api+"/projects/"+project+"/repositories/"+repository+"/artifacts/"+tag+"/tags/"+tag, nil)
	if err != nil {
		panic(err)
	}
	setBasicAuth(req)

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode == 200 {
		fmt.Println("OK")
	} else {
		fmt.Printf("FAILED - Repsonse from harbor: %s\n", resp.Status)
	}
}

func setBasicAuth(req *http.Request) {
	if harborCredentials == "" {
		if harborUsername == "" {
			panic("Please provide a username using the --username flag")
		}
		if harborPassword == "" {
			panic("Please provide a password using the --password flag")
		}

		req.SetBasicAuth(harborPassword, harborUsername)
	} else {
		in, err := ioutil.ReadFile(harborCredentials)
		if err != nil {
			panic(err)
		}
		var creds galasayaml.Credentials
		err = yaml.Unmarshal(in, &creds)
		if err != nil {
			panic(err)
		}
		req.SetBasicAuth(creds.Username, creds.Password)
	}
}
