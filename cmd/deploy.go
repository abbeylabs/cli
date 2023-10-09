/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
)

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains templates
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		repo, err := git.PlainOpen(path)
		if err != nil {
			fmt.Printf("Error when opening repo")
			return
		}

		worktree, err := repo.Worktree()
		if err != nil {
			return
		}

		_, err = worktree.Add("main.tf")
		if err != nil {
			return
		}

		_, err = worktree.Add("policies/common/common.rego")
		if err != nil {
			fmt.Println("Error when adding common.rego")
			return
		}

		status, err := worktree.Status()
		if err != nil {
			fmt.Println("Error when showing git status")
			return
		}

		fmt.Println(status)

		fmt.Println("git commit -m \"example go-git commit\"")
		commit, err := worktree.Commit("example go-git commit", &git.CommitOptions{})

		obj, err := repo.CommitObject(commit)
		fmt.Println(obj)

		err = repo.Push(&git.PushOptions{})
		if err != nil {
			fmt.Printf("Error when pushing %+v", err)
			return
		}

	},
}

var (
	path string
)

func init() {
	rootCmd.AddCommand(deployCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deployCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deployCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	deployCmd.Flags().StringVarP(&path, "path", "p", "", "Path of github repo")
}
