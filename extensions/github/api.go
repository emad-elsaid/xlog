package github

import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"os"
	"strings"

	"github.com/emad-elsaid/xlog"
	"github.com/emad-elsaid/xlog/extensions/shortcode"
	"github.com/google/go-github/v53/github"
	"golang.org/x/oauth2"
)

var githubTokenPossibleVariables = []string{"GITHUB_TOKEN", "GITHUB_API_TOKEN"}
var tokenNotAvailable = errors.New("Github token env variable not found in any of: " + strings.Join(githubTokenPossibleVariables, ", "))
var perPage = 100

func init() {
	shortcode.RegisterShortCode("github-search-issues", shortcode.ShortCode{Render: seachIssuesShortcode})
}

func seachIssuesShortcode(in xlog.Markdown) template.HTML {
	return template.HTML(issues(context.Background(), string(in)))
}

func token() (string, error) {
	for _, v := range githubTokenPossibleVariables {
		value := os.Getenv(v)
		if len(value) > 0 {
			return value, nil
		}
	}

	return "", tokenNotAvailable
}

func client() (*github.Client, error) {
	token, err := token()
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	return github.NewClient(tc), nil
}

func issues(ctx context.Context, query string) string {
	client, err := client()
	if err != nil {
		return err.Error()
	}

	result, _, err := client.Search.Issues(ctx, query, &github.SearchOptions{
		ListOptions: github.ListOptions{
			PerPage: perPage,
		},
	})
	if err != nil {
		return err.Error()
	}

	if len(result.Issues) == 0 {
		return fmt.Sprintf("No results for query: %s", query)
	}

	issues := "<ul>"
	for _, i := range result.Issues {
		assignee := i.GetUser()

		issues += fmt.Sprintf(`<li>
<span class="icon-text" >
	<figure class="icon image is-24x24 my-0 mx-2">
	  <img src="%s" class="is-rounded">
	</figure>
	<a href="%s">%s</a>
</span>
</li>`, assignee.GetAvatarURL(), i.GetHTMLURL(), i.GetTitle())
	}
	issues += "</ul>"

	return issues
}
