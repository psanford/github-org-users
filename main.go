// github-org-users fetches the list of ids and names for users
// in an organization. This requires an access token with permission
// read:org.
package main // import "github.com/psanford/github-org-users"

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/google/go-github/v24/github"
	"golang.org/x/oauth2"
)

var outputFormat = flag.String("format", "csv", "Output format (csv|json)")
var fetchFullUser = flag.Bool("fetch-full-user", false, "Fetch full user record")

func main() {
	flag.Parse()

	accessToken := os.Getenv("GITHUB_API_KEY")
	if accessToken == "" {
		fmt.Fprintf(os.Stderr, "Must set GITHUB_API_KEY in environment variable")
	}

	args := flag.Args()

	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "Usage: %s <org-name>\n", os.Args[0])
		os.Exit(1)
	}

	orgName := args[0]

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	header := []string{"login"}
	if *fetchFullUser {
		header = []string{"login", "name"}
	}

	csvOut := csv.NewWriter(os.Stdout)
	if *outputFormat == "csv" {
		csvOut.Write(header)
	}
	jsonOut := json.NewEncoder(os.Stdout)

	var opt github.ListMembersOptions
	for {
		members, resp, err := client.Organizations.ListMembers(ctx, orgName, &opt)
		if err != nil {
			log.Fatalf("Get org error: %s", err)
		}

		for _, u := range members {
			// fetch the full user record so we have the user's pretty name
			if *fetchFullUser {
				u, _, err = client.Users.Get(ctx, u.GetLogin())
				if err != nil {
					log.Fatalf("User fetch error: %s", err)
				}
			}

			if *outputFormat == "csv" {
				if len(header) == 1 {
					csvOut.Write([]string{u.GetLogin()})
				} else {
					csvOut.Write([]string{u.GetLogin(), u.GetName()})
				}
			} else {
				jsonOut.Encode(u)
			}
		}

		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
}
