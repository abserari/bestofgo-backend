package main

import (
	"context"
	"log"
	"os"
	"time"

	jsoniter "github.com/json-iterator/go"

	"github.com/google/go-github/v33/github"
	"golang.org/x/oauth2"
)

const token = "f9bdb40455a4e5f43d2c221620aeddfa7a789706"

var (
	nowTime = time.Now()

	yesterday = nowTime.AddDate(0, 0, -1)

	lastWeek = nowTime.AddDate(0, 0, -7)

	lastMonth = nowTime.AddDate(0, -1, 0)

	lastYear = nowTime.AddDate(-1, 0, 0)

	json = jsoniter.ConfigFastest
)

func main() {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	// clientv4 := githubv4.NewClient(httpClient)

	raw, err := LoadFile("./public/list-projects.json")
	if err != nil {
		log.Println(err)
		return
	}

	listPro, err := UnmarshalListProjects(raw)
	if err != nil {
		log.Println(err)
		return
	}

	var jsonfile = &JSONFile{}
	// date
	jsonfile.Date = time.Now().Format(time.RFC3339)

	for _, v := range listPro {
		// tags
		var Tags2Code = make(map[string]string)
		for _, t := range v.Tags {
			if Tags2Code[t] != "" {
				continue
			}
			Tags2Code[t] = t
			jsonfile.Tags = append(jsonfile.Tags, Tag{Name: t, Code: t})
		}

		// get projects down
		project, err := getProject(context.Background(), client, v.Org, v.Repo)
		if err != nil {
			log.Println(err, "jump this project")
			continue
		}
		project.Tags = append(project.Tags, v.Tags...)

		jsonfile.Projects = append(jsonfile.Projects, *project)
	}

	data, err := json.Marshal(jsonfile)
	fp, err := os.OpenFile("./public/projects.json", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Fatal(err)
	}

	defer fp.Close()
	_, err = fp.Write(data)
	if err != nil {
		log.Fatal(err)
	}
}
