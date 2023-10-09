/*
Copyright © 2023 ABBEY LABS
*/
package cmd

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
	"github.com/spf13/cobra"
	"log"

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
	repo = iota
	path
	reviewer
)

const (
	quickstart   = "https://github.com/abbeylabs/abbey-starter-kit-quickstart"
	googleGroups = "https://github.com/abbeylabs/google-workspace-starter-kit"
	gcpIdentity  = "https://github.com/abbeylabs/abbey-starter-kit-gcp-identity"
	okta         = "https://github.com/abbeylabs/abbey-starter-kit-okta"
	azure        = "https://github.com/abbeylabs/abbey-starter-kit-azure"
	snowflake    = "https://github.com/abbeylabs/abbey-starter-kit-snowflake"
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
	exampleList    list.Model
	timeExpiryList list.Model
	exampleChoice  string
	repo           string
	inputs         []textinput.Model
	path           string
	timeExpiry     string
	reviewer       string
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" || msg.Type == tea.KeyEsc {
			return m, tea.Quit
		}
		if msg.String() == "enter" {
			if m.exampleChoice == "" {
				i, ok := m.exampleList.SelectedItem().(item)
				if ok {
					m.exampleChoice = i.title
				}
			} else if m.repo == "" {
				i := m.inputs[repo].Value()
				m.repo = i
			} else if m.path == "" {
				i := m.inputs[path].Value()
				m.path = i
				return m, pullExampleRepo(m.path, m.repo)
			} else if m.timeExpiry == "" {
				i, ok := m.timeExpiryList.SelectedItem().(item)
				if ok {
					m.timeExpiry = i.title
				}
			} else if m.reviewer == "" {
				i := m.inputs[reviewer].Value()
				m.reviewer = i
			}
		}

	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.exampleList.SetSize(msg.Width-h, msg.Height-v)
		m.timeExpiryList.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	if m.exampleChoice == "" {
		m.exampleList, cmd = m.exampleList.Update(msg)
	} else if m.repo == "" {
		m.inputs[repo], cmd = m.inputs[repo].Update(msg)
	} else if m.path == "" {
		m.inputs[path], cmd = m.inputs[path].Update(msg)
	} else if m.timeExpiry == "" {
		m.timeExpiryList, cmd = m.timeExpiryList.Update(msg)
	} else if m.reviewer == "" {
		m.inputs[reviewer], cmd = m.inputs[reviewer].Update(msg)
	}
	return m, cmd
}

type statusMsg int

func pullExampleRepo(path string, repoName string) tea.Cmd {
	return func() tea.Msg {
		var err error
		if path == "" {
			path, err = os.Getwd()
			if err != nil {
				log.Println(err)
			}
		}

		//_, err = git.PlainClone(path, false, &git.CloneOptions{
		//	URL:      "https://github.com/" + repoName,
		//	Progress: os.Stdout,
		//})
		//
		//if err != nil {
		//	fmt.Println("Error initializing Abbey project, is Git configured locally?")
		//}

		return statusMsg(200)
	}
}

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
		case "gcp-groups":
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
		output += docStyle.Render(fmt.Sprintf(
			"Once you've done that, enter the name of the repo you created: %s\n\n%s",
			m.inputs[repo].View(),
			"(ctrl-c to quit)",
		) + "\n")
		return wordwrap.String(output, m.exampleList.Width())
	} else if m.path == "" {
		m.inputs[path].Focus()
		output := docStyle.Render(fmt.Sprintf(
			"Enter in the path the repo will be cloned to, or leave empty to clone into current directory: %s\n\n%s",
			m.inputs[path].View(),
			"(ctrl-c to quit)",
		) + "\n")
		return wordwrap.String(output, m.exampleList.Width())
	} else if m.timeExpiry == "" {
		return docStyle.Render(m.timeExpiryList.View())
	} else if m.reviewer == "" {
		m.inputs[reviewer].Focus()
		output := docStyle.Render(fmt.Sprintf(
			"What's the email address you used for Abbey?\n\n%s\n\n%s",
			m.inputs[reviewer].View(),
			"(ctrl-c to quit)",
		) + "\n")
		return wordwrap.String(output, m.exampleList.Width())
	} else {
		output := wordwrap.String(docStyle.Render(fmt.Sprintf("Thanks for setting up Abbey! Press ESC or ctrl-c to exit")), m.exampleList.Width())
		output += wordwrap.String(docStyle.Render(fmt.Sprintf("Repo name is %s!", m.repo)), m.exampleList.Width())
		output += wordwrap.String(docStyle.Render(fmt.Sprintf("Time expiry is %s!", m.timeExpiry)), m.exampleList.Width())
		output += wordwrap.String(docStyle.Render(fmt.Sprintf("Reviewer email address is %s!", m.reviewer)), m.exampleList.Width())
		return output
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

		inputs := make([]textinput.Model, 3)
		inputs[repo] = textinput.New()
		inputs[repo].Placeholder = "github-username/repo-name"
		inputs[repo].Focus()
		inputs[repo].CharLimit = 20
		inputs[repo].Width = 30
		inputs[repo].Prompt = ""

		inputs[path] = textinput.New()
		inputs[path].Placeholder = "/Users/alice/Documents"
		inputs[path].CharLimit = 50
		inputs[path].Width = 50
		inputs[path].Prompt = ""

		inputs[reviewer] = textinput.New()
		inputs[reviewer].Placeholder = "alice@example.com"
		inputs[reviewer].CharLimit = 30
		inputs[reviewer].Width = 40
		inputs[reviewer].Prompt = ""

		timeExpiryOptions := []list.Item{
			item{title: "5m", desc: "5 minutes"},
			item{title: "1hr", desc: "1 hour"},
			item{title: "1d", desc: "1 day"},
			item{title: "7d", desc: "7 days"},
		}

		m := model{
			exampleList:    list.New(items, list.NewDefaultDelegate(), 0, 0),
			timeExpiryList: list.New(timeExpiryOptions, list.NewDefaultDelegate(), 0, 0),
			inputs:         inputs}

		m.timeExpiryList.Title = "Time Expiry Options"
		m.exampleList.Title = "Abbey Examples"

		p := tea.NewProgram(m, tea.WithAltScreen())

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
