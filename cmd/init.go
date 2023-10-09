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
	"github.com/go-git/go-git/v5"
	"github.com/muesli/reflow/wordwrap"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strings"
	"text/template"
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
	repoTextInput = iota
	pathTextInput
	reviewerTextInput
	tokenTextInput
	deployTextInput
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
	reviewer     string
	timeExpiry   string
	accessOutput string
	repo         string
	azureUPN     string
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
	textReplaced   bool
	deployed       string
	accessOutput   string
	tokenSetup     bool
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
				i := m.inputs[repoTextInput].Value()
				m.repo = i
			} else if m.path == "" {
				i := m.inputs[pathTextInput].Value()
				m.path = i
				var err error
				if m.path == "" || m.path == "." {
					ex, err := os.Executable()
					if err != nil {
						panic(err)
					}
					m.path = filepath.Dir(ex)
				}

				m.path = m.path + "/" + m.exampleChoice
				url := "https://github.com/" + m.repo
				_, err = git.PlainClone(m.path, false, &git.CloneOptions{
					URL: url,
				})

				if err != nil {
					fmt.Print(fmt.Sprintf("Error initializing Abbey project, is Git configured locally?\nError: %v", err))
				}
			} else if m.timeExpiry == "" {
				i, ok := m.timeExpiryList.SelectedItem().(item)
				if ok {
					m.timeExpiry = i.title
				}
			} else if m.reviewer == "" {
				i := m.inputs[reviewerTextInput].Value()
				m.reviewer = i
			} else if !m.tokenSetup {
				i := m.inputs[tokenTextInput].Value()
				if strings.ToLower(i) == "yes" || strings.ToLower(i) == "y" {
					m.tokenSetup = true
				}
			} else if m.deployed == "" {
				i := m.inputs[deployTextInput].Value()
				m.deployed = i
				if strings.ToLower(m.deployed) == "yes" || strings.ToLower(m.deployed) == "y" {
					repo, _ := git.PlainOpen(m.path)
					worktree, _ := repo.Worktree()
					_, _ = worktree.Add("main.tf")
					_, _ = worktree.Add("policies/common/common.rego")

					commit, _ := worktree.Commit("Configuring Abbey via CLI", &git.CommitOptions{})

					_, _ = repo.CommitObject(commit)
					_ = repo.Push(&git.PushOptions{})
				}
			}
		}

		// if all fields are filled, replace the text in files
		if m.timeExpiry != "" && m.reviewer != "" && m.repo != "" && m.path != "" && m.exampleChoice != "" {
			mainFile := "main.tf"
			files, err := template.New(mainFile).ParseFiles("cmd/templ/" + m.exampleChoice + "/main.tf")
			if err != nil {
				panic(err)
			}

			tmpl := template.Must(files, err)

			if m.accessOutput == "" {
				m.accessOutput = "github://" + m.repo + "/access.tf"
			}

			tfInput := MainTfTemplate{
				Reviewer:     m.reviewer,
				PolicyBundle: "github://" + m.repo + "/policies",
				AccessOutput: m.accessOutput,
				// hacky fix since Go recognizes this as a template value
				AbbeyEmail: "{{ .data.system.abbey.identities.abbey.email }}",
				AzureUPN:   azureUPN,
			}

			f, err := os.Create(m.path + "/" + mainFile)
			if err != nil {
				fmt.Println("create mainFile failed: ", err)
			}

			err = tmpl.Execute(f, tfInput)
			if err != nil {
				panic(err)
			}

			err = f.Close()
			if err != nil {
				panic(err)
			}

			policyFile := "common.rego"
			tmpl = template.Must(template.New(policyFile).ParseFiles("cmd/templ/" + m.exampleChoice + "/" + policyFile))

			policyTfInput := PolicyTfTemplate{
				TimeExpiry: m.timeExpiry,
			}

			f, err = os.Create(m.path + "/policies/common/" + policyFile)
			if err != nil {
				fmt.Println("create policy file failed: ", err)
			}

			err = tmpl.Execute(f, policyTfInput)
			if err != nil {
				panic(err)
			}

			err = f.Close()
			if err != nil {
				panic(err)
			}

			m.textReplaced = true
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
		m.inputs[repoTextInput], cmd = m.inputs[repoTextInput].Update(msg)
	} else if m.path == "" {
		m.inputs[pathTextInput], cmd = m.inputs[pathTextInput].Update(msg)
	} else if m.timeExpiry == "" {
		m.timeExpiryList, cmd = m.timeExpiryList.Update(msg)
	} else if m.reviewer == "" {
		m.inputs[reviewerTextInput], cmd = m.inputs[reviewerTextInput].Update(msg)
	} else if !m.tokenSetup {
		m.inputs[tokenTextInput], cmd = m.inputs[tokenTextInput].Update(msg)
	} else if m.deployed == "" {
		m.inputs[deployTextInput], cmd = m.inputs[deployTextInput].Update(msg)
	}
	return m, cmd
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
		output += docStyle.Render(fmt.Sprintf("The repo you created must be public."))
		output += docStyle.Render(fmt.Sprintf(
			"Once you've done that, enter the name of the repo you created: %s\n\n%s",
			m.inputs[repoTextInput].View(),
			"",
		) + "\n")
		return wordwrap.String(output, m.exampleList.Width())
	} else if m.path == "" {
		m.inputs[pathTextInput].Focus()
		output := docStyle.Render(fmt.Sprintf(
			"Enter in the path the repo will be cloned to, or leave empty to clone into current directory: %s\n\n%s",
			m.inputs[pathTextInput].View(),
			"",
		) + "\n")
		return wordwrap.String(output, m.exampleList.Width())
	} else if m.timeExpiry == "" {
		return docStyle.Render(m.timeExpiryList.View())
	} else if m.reviewer == "" {
		m.inputs[reviewerTextInput].Focus()
		output := docStyle.Render(fmt.Sprintf(
			"What's the email address you used for Abbey?\n\n%s\n\n%s",
			m.inputs[reviewerTextInput].View(),
			"",
		) + "\n")
		return wordwrap.String(output, m.exampleList.Width())
	} else if !m.tokenSetup {
		m.inputs[tokenTextInput].Focus()
		output := docStyle.Render(fmt.Sprintf("Before you deploy Abbey, you must set up Abbey tokens in Github."))
		output += docStyle.Render(fmt.Sprintf("Follow instructions at https://docs.abbey.io/product/deploying-your-grant-kit, and return here to confirm once done."))
		output += docStyle.Render(fmt.Sprintf(
			"Have the Abbey tokens been configured? [Yes | No]\n\n%s\n\n%s",
			m.inputs[tokenTextInput].View(),
			"",
		) + "\n")
		return wordwrap.String(output, m.exampleList.Width())

	} else if m.deployed == "" {
		m.inputs[deployTextInput].Focus()
		output := docStyle.Render(fmt.Sprintf(
			"Confirm deployment to Github? [Yes | No]\n\n%s\n\n%s",
			m.inputs[deployTextInput].View(),
			"",
		) + "\n")
		return wordwrap.String(output, m.exampleList.Width())

	} else {
		output := wordwrap.String(docStyle.Render(fmt.Sprintf("Thanks for setting up Abbey! Press ESC or ctrl-c to exit")), m.exampleList.Width())
		output += wordwrap.String(docStyle.Render(fmt.Sprintf("Repo name is %s!", m.repo)), m.exampleList.Width())
		output += wordwrap.String(docStyle.Render(fmt.Sprintf("Path is %s!", m.path)), m.exampleList.Width())
		output += wordwrap.String(docStyle.Render(fmt.Sprintf("Time expiry is %s!", m.timeExpiry)), m.exampleList.Width())
		output += wordwrap.String(docStyle.Render(fmt.Sprintf("Reviewer email address is %s!", m.reviewer)), m.exampleList.Width())
		return output
	}
}

// initCmd represents the create command
var initCmd = &cobra.Command{
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

		inputs := make([]textinput.Model, 7)
		inputs[repoTextInput] = textinput.New()
		inputs[repoTextInput].Placeholder = "github-username/repo-name"
		inputs[repoTextInput].Focus()
		inputs[repoTextInput].CharLimit = 80
		inputs[repoTextInput].Width = 80
		inputs[repoTextInput].Prompt = ""

		inputs[pathTextInput] = textinput.New()
		inputs[pathTextInput].Placeholder = "/Users/alice/Documents"
		inputs[pathTextInput].CharLimit = 80
		inputs[pathTextInput].Width = 80
		inputs[pathTextInput].Prompt = ""

		inputs[reviewerTextInput] = textinput.New()
		inputs[reviewerTextInput].Placeholder = "alice@example.com"
		inputs[reviewerTextInput].CharLimit = 80
		inputs[reviewerTextInput].Width = 80
		inputs[reviewerTextInput].Prompt = ""

		timeExpiryOptions := []list.Item{
			item{title: "5m", desc: "5 minutes"},
			item{title: "1hr", desc: "1 hour"},
			item{title: "1d", desc: "1 day"},
			item{title: "7d", desc: "7 days"},
		}

		inputs[tokenTextInput] = textinput.New()
		inputs[tokenTextInput].Placeholder = "Yes"
		inputs[tokenTextInput].CharLimit = 80
		inputs[tokenTextInput].Width = 80
		inputs[tokenTextInput].Prompt = ""

		inputs[deployTextInput] = textinput.New()
		inputs[deployTextInput].Placeholder = "Yes"
		inputs[deployTextInput].CharLimit = 80
		inputs[deployTextInput].Width = 80
		inputs[deployTextInput].Prompt = ""

		m := model{
			exampleList:    list.New(items, list.NewDefaultDelegate(), 0, 0),
			timeExpiryList: list.New(timeExpiryOptions, list.NewDefaultDelegate(), 0, 0),
			inputs:         inputs,
			textReplaced:   false,
			reviewer:       reviewer,
			timeExpiry:     timeExpiry,
			repo:           repo,
			accessOutput:   accessOutput,
		}

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
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	initCmd.Flags().StringVarP(&reviewer, "reviewer", "r", "", "Abbey email address of the reviewer")
	initCmd.Flags().StringVarP(&repo, "repo", "", "", "Git repo name in the format github-username/repo-name")
	initCmd.Flags().StringVarP(&timeExpiry, "timeExpiry", "t", "", "Time expiry of permissions")
	initCmd.Flags().StringVarP(&accessOutput, "accessOutput", "a", "", "Location of Access Output")
	initCmd.Flags().StringVarP(&azureUPN, "azureUPN", "", "", "Azure UPN [for use with Azure example]")

	err := initCmd.MarkFlagRequired("example")
	if err != nil {
		return
	}

}
