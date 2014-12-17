package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Resume struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Place string `json:"place"`
	Jobs  []Job  `json:"jobs"`
	Educ  string `json:"education"`
	Link  string `json:"link"`
}

type Job struct {
	Id       int    `json:"id"`
	Role     string `json:"role"`
	Compagny string `json:"compagny"`
}

const (
	URL = "http://www.indeed.com"
)

var (
	Resumes = []*Resume{}
	At      = 0
)

func main() {
	for i := 0; i <= 10; i++ {
		err := fetch()
		if err != nil {
			break
		}
	}
	Phase(Resumes)
}

func fetch() error {
	from := "http://www.indeed.com/resumes/sales/in-montreal-QC?co=CA&sort=date&start=" + strconv.Itoa(At)
	doc, err := goquery.NewDocument(from)
	if err != nil {
		return err
	}
	doc.Find(".sre").Each(func(i int, s *goquery.Selection) {
		r := &Resume{}
		r.Id = i + 1
		r.Place = s.Find(".app_name span").Text()
		r.Name = s.Find(".app_link").Text()
		r.Jobs = func() []Job {
			js := []Job{}
			j := Job{}
			s.Find(".experience").Each(func(i int, s *goquery.Selection) {
				d := strings.Split(s.Text(), "-")
				j.Id = i + 1
				if len(d) > 1 {
					j.Role = d[0]
					j.Compagny = d[1]
				}
				js = append(js, j)
			})
			return js
		}()
		r.Educ = s.Find(".education").Text()
		link, _ := s.Find(".app_name a").Attr("href")
		r.Link = URL + link
		Resumes = append(Resumes, r)
	})
	At = len(Resumes)
	fmt.Println(At)
	return nil
}

func Phase(r []*Resume) {
	data, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		fmt.Println(err)
	}
	err = ioutil.WriteFile("data.json", data, 0x777)
	if err != nil {
		fmt.Println(err)
	}
}
