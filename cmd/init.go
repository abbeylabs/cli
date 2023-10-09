/*
Copyright © 2023 ABBEY LABS
*/
package cmd

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"os"
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

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	docStyle          = lipgloss.NewStyle().Margin(1, 2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type model struct {
	exampleList   list.Model
	exampleChoice string
	altscreen     bool
	repo          string
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		if msg.String() == "enter" {
			i, ok := m.exampleList.SelectedItem().(item)
			if ok {
				m.exampleChoice = i.title
			}
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.exampleList.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.exampleList, cmd = m.exampleList.Update(msg)
	return m, cmd
}

//func pullExampleRepo(example string) tea.Cmd {
//
//}

func (m model) View() string {
	if m.exampleChoice == "" {
		return docStyle.Render(m.exampleList.View())
	} else if m.repo == "" {

		var url string
		switch m.exampleChoice {
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
		}

		output := docStyle.Render(fmt.Sprintf("Great! You've chosen %s to get started.", m.exampleChoice))
		output += docStyle.Render(fmt.Sprintf("Navigate to %s to create a repo from the template in your own github account.", url))
		return output
	} else {
		return docStyle.Render(fmt.Sprintf("Thanks for setting up Abbey!"))
	}
}

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes an Abbey example",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		items := []list.Item{
			item{title: "quickstart", desc: "Fastest example to get Abbey running"},
			item{title: "google-groups", desc: "Managing access to groups via Google Workspace"},
			item{title: "gcp-groups", desc: "Managing access to groups via Google Cloud Platform (GCP)"},
			item{title: "snowflake", desc: "Managing access to tables in Snowflake"},
			item{title: "azure", desc: "Managing access to groups in Azure AD"},
		}

		m := model{exampleList: list.New(items, list.NewDefaultDelegate(), 0, 0)}
		m.exampleList.Title = "Abbey Examples"

		p := tea.NewProgram(m)

		if _, err := p.Run(); err != nil {
			fmt.Println("Error running program:", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	//createCmd.Flags().StringVarP(&example, "example", "e", "", "Abbey Example [Name]|[Github-URL] (required)")
	//createCmd.Flags().BoolVarP(&preReqs, "preReqs", "", false, "Are all the pre-requisites i.e. connecting GitHub, setting up Abbey account done?")
	//createCmd.Flags().StringVarP(&example, "user", "u", "", "Abbey user email address")
	//createCmd.Flags().StringVarP(&reviewer, "reviewer", "r", "", "Abbey email address of the reviewer")
	//createCmd.Flags().StringVarP(&repo, "repo", "", "", "Git repo name")
	//createCmd.Flags().StringVarP(&githubUsername, "githubUsername", "g", "", "Github Username")
	//createCmd.Flags().StringVarP(&timeExpiry, "timeExpiry", "t", "", "Time expiry of permissions")
	//createCmd.Flags().StringVarP(&policyBundle, "policyBundle", "p", "", "Location of Policy Bundle")
	//createCmd.Flags().StringVarP(&accessOutput, "accessOutput", "a", "access.tf", "Location of Access Output")
	//createCmd.Flags().StringVarP(&azureUPN, "azureUPN", "", "", "Azure UPN [for use with Azure example]")

	err := createCmd.MarkFlagRequired("example")
	if err != nil {
		return
	}

}
