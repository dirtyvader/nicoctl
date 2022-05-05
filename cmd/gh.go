/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"context"
	"fmt"

	"github.com/shurcooL/githubv4"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

// ghCmd represents the gh command
var ghCmd = &cobra.Command{
	Use:   "gh",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		token := viper.GetString("github.access_token")
		src := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)
		httpClient := oauth2.NewClient(context.Background(), src)
		client := githubv4.NewClient(httpClient)

		var query struct {
			Viewer struct {
				Login        githubv4.String
				CreatedAt    githubv4.DateTime
				AvatarUrl    githubv4.URI
				Repositories struct {
					Edges []struct {
						Node struct {
							Id        githubv4.ID
							Name      githubv4.String
							SshUrl    githubv4.String
							UpdatedAt githubv4.DateTime
							Languages struct {
								Nodes []struct {
									Name  githubv4.String
									Color githubv4.String
								}
							} `graphql:"languages(first: 5)"`
						}
					}
				} `graphql:"repositories(first: 5)"`
			}
		}

		err := client.Query(context.Background(), &query, nil)
		if err != nil {
			fmt.Printf("Error : %v\n", err)
		}

		fmt.Println("       Login:", query.Viewer.Login)
		fmt.Println("   CreatedAt:", query.Viewer.CreatedAt)
		fmt.Println("   AvatarUrl:", query.Viewer.AvatarUrl.URL)
		fmt.Println("Repositories:")
		for i := 0; i < len(query.Viewer.Repositories.Edges); i++ {
			currentEdge := query.Viewer.Repositories.Edges[i]
			fmt.Println("########################################")
			fmt.Printf("\t-        ID: %s\n", currentEdge.Node.Id)
			fmt.Printf("\t-      Name: %s\n", currentEdge.Node.Name)
			fmt.Printf("\t-    SshUrl: %s\n", currentEdge.Node.SshUrl)
			fmt.Printf("\t- UpdatedAt: %v\n", currentEdge.Node.UpdatedAt)
			fmt.Printf("\t- Languages:\n")
			for j := 0; j < len(currentEdge.Node.Languages.Nodes); j++ {
				currentNode := currentEdge.Node.Languages.Nodes[j]
				fmt.Printf("\t\t-  Name: %s\n", currentNode.Name)
				fmt.Printf("\t\t- Color: %s\n\n", currentNode.Color)
			}
			fmt.Println("########################################")
		}

	},
}

func init() {
	rootCmd.AddCommand(ghCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// ghCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// ghCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
