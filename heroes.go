package main

import (
	"context"
	"log"

	"github.com/shurcooL/githubv4"
)

// Graphql query

// user(login: "abserari") {
//     login,
//     name,
//     bio,
//     avatarUrl,
//     websiteUrl,
//     followers {
//       totalCount
//     }
//   }

func getHero(ctx context.Context, client *githubv4.Client, login string) (*Hero, error) {
	var q struct {
		User struct {
			Login      string
			Name       string
			Bio        string
			AvatarUrl  string
			WebsiteUrl string
			Followers  struct {
				TotalCount int
			}
		} `graphql:"user(login: $login)"`
	}
	variables := map[string]interface{}{
		"login": githubv4.String(login),
	}
	err := client.Query(context.Background(), &q, variables)
	if err != nil {
		// Handle error.
		log.Println(err)
		return nil, err
	}

	var hero = &Hero{
		Username:  q.User.Login,
		Avatar:    q.User.AvatarUrl,
		Followers: q.User.Followers.TotalCount,
		Blog:      q.User.WebsiteUrl,
		Name:      q.User.Name,
		// Projects  []string `json:"projects"`
		Bio: q.User.Bio,
		// Npm       string   `json:"npm"`
		// Modules   int      `json:"modules"`
	}
	return hero, nil
}
