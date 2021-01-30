package main

import (
	"context"
	"log"
	"time"

	"github.com/google/go-github/v33/github"
)

func getProject(ctx context.Context, client *github.Client, org, repo string) (*Project, error) {
	startTime := time.Now()
	log.Println("start now at ", startTime.Sub(startTime))
	var project = &Project{}
	stats, _, err := client.Repositories.Get(ctx, org, repo)
	if err != nil {
		return nil, err
	}
	err = getStarsTrending(ctx, client, project, org, repo)
	if err != nil {
		return nil, err
	}
	log.Println("end network at ", startTime.Sub(time.Now()))

	project.Name = repo
	project.Stars = stats.GetStargazersCount()
	project.FullName = stats.GetFullName()
	project.Description = stats.GetDescription()
	project.OwnerID = stats.GetOwner().GetID()
	project.PushedAt = stats.GetCreatedAt().Format("2006-01-02")
	project.CreatedAt = stats.GetCreatedAt().Format("2006-01-02")
	project.URL = stats.GetHomepage()
	opt := &github.ListContributorsOptions{
		ListOptions: github.ListOptions{PerPage: 10},
	}
	contributors, _, err := client.Repositories.ListContributors(context.Background(), org, repo, opt)
	project.ContributorCount = len(contributors)

	// strings, err := json.Marshal(stats)
	// log.Println(string(strings))

	project.Tags = stats.Topics

	// get all pages of results
	return project, nil
}

func getStarsTrending(ctx context.Context, client *github.Client, project *Project, org, repo string) error {
	var (
		daily   int
		weekly  int
		monthly int
		yearly  int
	)

	var stars []*github.Stargazer
	var page int = 0
	for {
		stargazers, resp, err := client.Activity.ListStargazers(context.Background(), org, repo, &github.ListOptions{Page: page})
		if err != nil {
			log.Println(err)
			return err
		}
		stars = append(stars, stargazers...)
		// todo: make client pool or snap or graphql to solve this ratelimit and slowly call problem
		if resp.NextPage == 0 || resp.NextPage > 10 {
			break
		}
		page = resp.NextPage
	}

	for _, v := range stars {
		if v.StarredAt.Time.After(yesterday) {
			daily++
			weekly++
			monthly++
			yearly++
		} else if v.StarredAt.Time.After(lastWeek) {
			weekly++
			monthly++
			yearly++
		} else if v.StarredAt.Time.After(lastMonth) {
			monthly++
			yearly++
		} else if v.StarredAt.Time.After(lastYear) {
			monthly++
			yearly++
		}
	}
	project.Trends.Daily, project.Trends.Weekly, project.Trends.Monthly, project.Trends.Yearly = daily, weekly, monthly, yearly
	return nil
}
