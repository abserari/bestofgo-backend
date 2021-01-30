package main

import (
	"context"
	"log"
	"os"
	"sync"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/shurcooL/githubv4"

	"github.com/google/go-github/v33/github"
	"golang.org/x/oauth2"
)

const token = "test"

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
	clientv4 := githubv4.NewClient(tc)

	rawHeroList, err := LoadFile("./public/list-hof.json")
	if err != nil {
		log.Println(err)
		return
	}

	listHero, err := UnmarshalListHeroes(rawHeroList)
	if err != nil {
		log.Println(err)
		return
	}

	var jsonHero = &HeroFile{}
	jsonHero.Date = time.Now().Format(time.RFC3339)

	for _, v := range listHero {
		hero, err := getHero(ctx, clientv4, v.Login)
		if err != nil {
			log.Println("jump this man")
			continue
		}

		hero.Projects = v.Projects
		jsonHero.Heroes = append(jsonHero.Heroes, *hero)
	}

	dataHero, err := json.Marshal(jsonHero)
	if err != nil {
		log.Println("hero file marshal failed")
		return
	}

	heroFp, err := os.OpenFile("./public/hof.json", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Fatal(err)
	}

	defer heroFp.Close()
	_, err = heroFp.Write(dataHero)
	if err != nil {
		log.Fatal(err)
	}

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

	var wg sync.WaitGroup

	for _, v := range listPro {
		wg.Add(1)
		go func(v OriginProject) {
			log.Println("handle project", v.Repo)
			defer wg.Done()
			// get projects down
			project, err := getProject(context.Background(), client, v.Org, v.Repo)
			if err != nil {
				log.Println(err, "jump this project")
				return
			}

			// tags
			var Code2Tags = make(map[string]string)
			for _, code := range v.Tags {
				Code2Tags[code] = code
			}
			for _, t := range project.Tags {
				Code2Tags[t] = t
			}
			for k, v := range Code2Tags {
				jsonfile.Tags = append(jsonfile.Tags, Tag{Code: k, Name: v})
			}

			// projects
			project.Tags = append(project.Tags, v.Tags...)
			jsonfile.Projects = append(jsonfile.Projects, *project)
		}(v)
	}

	wg.Wait()
	log.Println("handle down")

	data, err := json.Marshal(jsonfile)
	if err != nil {
		log.Println(err)
	}
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
