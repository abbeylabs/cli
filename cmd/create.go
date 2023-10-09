/*
Copyright © 2023 ABBEY LABS
*/
package cmd

import (
	"bufio"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
	"os"
	"strings"
	template "text/template"
)

type MainTfTemplate struct {
	Reviewer     string
	PolicyBundle string
	AccessOutput string
	AbbeyEmail   string
	AzureUPN     string
}

type PolicyTfTemplate struct {
	TimeExpiry string
}

const (
	quickstart        = "https://github.com/abbeylabs/abbey-starter-kit-quickstart"
	googleGroups      = "https://github.com/abbeylabs/google-workspace-starter-kit"
	gcpIdentity       = "https://github.com/abbeylabs/abbey-starter-kit-gcp-identity"
	okta              = "https://github.com/abbeylabs/abbey-starter-kit-okta"
	azure             = "https://github.com/abbeylabs/abbey-starter-kit-azure"
	snowflake         = "https://github.com/abbeylabs/abbey-starter-kit-snowflake"
	supported_options = "Supported Options: quickstart, googleGroups, gcpIdentity, okta, azure, snowflake"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes an Abbey example",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		var url string
		switch example {
		case "quickstart":
			url = quickstart
		case "google-groups":
			url = googleGroups
		case "gcp-identity":
			url = gcpIdentity
		case "okta":
			url = okta
		case "snowflake":
			url = snowflake
		case "azure":
			url = azure
		default:
			fmt.Println("Please pick one of the supported options\n" + supported_options)
			return
		}

		fmt.Println("Initializing Abbey project...")
		path := "/tmp/foo/" + example
		_, err := git.PlainClone(path, false, &git.CloneOptions{
			URL:      url,
			Progress: os.Stdout,
		})
		if err != nil {
			fmt.Println("Error initializing Abbey project, is Git configured locally?")
			return
		}

		inputReader := bufio.NewReader(os.Stdin)

		if !preReqs {
			var isGithubConnected string
			fmt.Println("Have you connected Abbey to your GitHub account yet? [yes | n]")

			isGithubConnected, _ = inputReader.ReadString('\n')
			isGithubConnected = strings.TrimSuffix(isGithubConnected, "\n")

			if isGithubConnected != "yes" {
				fmt.Println("Please connect Abbey to your Github account at https://app.abbey.io/connections")
			}
		}

		if policyBundle == "" {
			// prompt user for policy bundle location
			if repo == "" {
				fmt.Println("What is the name of your repo?")
				repo, _ = inputReader.ReadString('\n')
				repo = strings.TrimSuffix(repo, "\n")
			}

			if githubUsername == "" {
				fmt.Println("What is your github username?")
				githubUsername, _ = inputReader.ReadString('\n')
				githubUsername = strings.TrimSuffix(githubUsername, "\n")
			}

			policyBundle = "github://" + githubUsername + "/" + repo
		}

		if example == "azure" {

		}

		if !writeMainTf(path) {
			fmt.Println("Failed to write main.tf!")
			return
		}

		if !writePolicyRego(path) {
			fmt.Println("Failed to write policy file!")
			return
		}

		fmt.Println("Abbey setup complete!")
	},
}

func writePolicyRego(path string) bool {
	inputReader := bufio.NewReader(os.Stdin)

	policyFile := "common.rego"
	tmpl := template.Must(template.New(policyFile).ParseFiles("templates/abbey-starter-kit-quickstart/" + policyFile))

	if timeExpiry == "" {
		fmt.Println("How long do you want your permissions to last? Valid values: [_m][_h]")
		timeExpiry, _ = inputReader.ReadString('\n')
		timeExpiry = strings.TrimSuffix(timeExpiry, "\n")
	}

	policyTfInput := PolicyTfTemplate{
		TimeExpiry: timeExpiry,
	}

	f, err := os.Create(path + "/policies/common/" + policyFile)
	if err != nil {
		fmt.Println("create policy file: ", err)
		return false
	}

	err = tmpl.Execute(f, policyTfInput)
	if err != nil {
		panic(err)
	}

	err = f.Close()
	if err != nil {
		panic(err)
	}
	return true
}

func writeMainTf(path string) bool {
	mainFile := "main.tf"
	tmpl := template.Must(template.New(mainFile).ParseFiles("templates/abbey-starter-kit-" + example + "/main.tf"))

	tfInput := MainTfTemplate{
		Reviewer:     reviewer,
		PolicyBundle: policyBundle,
		AccessOutput: accessOutput,
		// hacky fix since Go recognizes this as a template value
		AbbeyEmail: "{{ .data.system.abbey.identities.abbey.email }}",
		AzureUPN:   azureUPN,
	}

	f, err := os.Create(path + "/" + mainFile)
	if err != nil {
		fmt.Println("create mainFile: ", err)
		return false
	}

	err = tmpl.Execute(f, tfInput)
	if err != nil {
		panic(err)
	}

	err = f.Close()
	if err != nil {
		panic(err)
	}

	return true
}

var (
	example        string
	reviewer       string
	timeExpiry     string
	policyBundle   string
	accessOutput   string
	repo           string
	githubUsername string
	preReqs        bool
	azureUPN       string
)

func init() {
	rootCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	createCmd.Flags().StringVarP(&example, "example", "e", "", "Abbey Example [Name]|[Github-URL] (required)")
	createCmd.Flags().BoolVarP(&preReqs, "preReqs", "", false, "Are all the pre-requisites i.e. connecting GitHub, setting up Abbey account done?")
	createCmd.Flags().StringVarP(&example, "user", "u", "", "Abbey user email address")
	createCmd.Flags().StringVarP(&reviewer, "reviewer", "r", "", "Abbey email address of the reviewer")
	createCmd.Flags().StringVarP(&repo, "repo", "", "", "Git repo name")
	createCmd.Flags().StringVarP(&githubUsername, "githubUsername", "g", "", "Github Username")
	createCmd.Flags().StringVarP(&timeExpiry, "timeExpiry", "t", "", "Time expiry of permissions")
	createCmd.Flags().StringVarP(&policyBundle, "policyBundle", "p", "", "Location of Policy Bundle")
	createCmd.Flags().StringVarP(&accessOutput, "accessOutput", "a", "access.tf", "Location of Access Output")
	createCmd.Flags().StringVarP(&azureUPN, "azureUPN", "", "", "Azure UPN [for use with Azure example]")

	err := createCmd.MarkFlagRequired("example")
	if err != nil {
		return
	}

}
