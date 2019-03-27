// github-org-users fetches the list of ids and names for users
// in an organization. This requires an access token with permission
// read:org.
package main // import "github.com/psanford/github-org-users"

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"github.com/google/go-github/v24/github"
	"golang.org/x/oauth2"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s <org-name> <access-token>\n", os.Args[0])
		os.Exit(1)
	}

	orgName := os.Args[1]
	accessToken := os.Args[2]

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	out := csv.NewWriter(os.Stdout)
	out.Write([]string{"login", "name"})

	var opt github.ListMembersOptions
	for {
		members, resp, err := client.Organizations.ListMembers(ctx, orgName, &opt)
		if err != nil {
			log.Fatalf("Get org error: %s", err)
		}

		for _, u := range members {
			// fetch the full user record so we have the user's pretty name
			u, _, err = client.Users.Get(ctx, u.GetLogin())
			if err != nil {
				log.Fatalf("User fetch error: %s", err)
			}
			out.Write([]string{u.GetLogin(), u.GetName()})
			out.Flush()
		}

		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
}
