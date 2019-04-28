package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	_ "github.com/go-sql-driver/mysql"
)

type Book struct {
	Title  string //标题
	Img    string //图片地址
	Author string //作者
	Sell   string //价格
	Url    string //访问链接
}

func ExampleScrape() {

	db, err := sql.Open("mysql", "root:123456@tcp(119.23.227.157:3306)/wenda")
	if err != nil {
		log.Fatal(err)
	}

	for i := 321450693; i > 300000000; i-- {
		var questionId int64
		res, err := http.Get("https://www.zhihu.com/question/" + strconv.Itoa(i))
		if err != nil || res.StatusCode != 200 {
			continue
		}
		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			log.Fatal(err)
		}

		doc.Find(".QuestionHeader .QuestionHeader-content .QuestionHeader-main").Each(func(i int, s *goquery.Selection) {
			questionTitle := s.Find(".QuestionHeader-title").Text()
			questionContent := s.Find(".QuestionHeader-detail").Text()
			questionContent = questionContent[0 : len(questionContent)-12]

			stmt, err := db.Prepare("INSERT INTO question (title,content,create_time,user_id,comment_count,status) VALUES(?,?,?,?,?,?)")
			if err != nil {
				log.Fatal(err)
			}
			res, err := stmt.Exec(questionTitle, questionContent, time.Now(), 1, 2, 0)
			if err != nil {
				log.Fatal(err)
			}
			questionId, err = res.LastInsertId()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(questionId)
		})

		// doc.Find(".ListShortcut .List .List-header").Each(func(i int, s *goquery.Selection) {
		// 	commentCount := s.Find("h4").Text()
		// 	fmt.Println("commentCount : ", commentCount)

		// })
		doc.Find(".ListShortcut .List .List-item ").Each(func(i int, s *goquery.Selection) {
			head_url, _ := s.Find("a img").Attr("src")
			author := s.Find(".AuthorInfo-head").Text()
			stmt, err := db.Prepare("INSERT INTO user(name,username,password,salt,head_url) VALUES(?,?,?,?,?)")
			if err != nil {
				log.Fatal(err)
			}
			res, err := stmt.Exec(author, "", "", "", head_url)
			if err != nil {
				log.Fatal(err)
			}
			user_id, err := res.LastInsertId()
			if err != nil {
				log.Fatal(err)
			}

			voters := s.Find(".Voters").Text()
			voters = strings.Split(voters, " ")[0]
			content, _ := s.Find(".RichContent-inner").Html()
			createTime := s.Find(".ContentItem-time").Text()
			createTime = strings.Split(createTime, " ")[1]

			tmt, err := db.Prepare("INSERT INTO comment(user_id,entity_id,entity_type,content,create_time,status) VALUES(?,?,?,?,?,?)")
			if err != nil {
				log.Fatal(err)
			}
			_, err = tmt.Exec(user_id, questionId, 1, content, time.Now(), 0)
			if err != nil {
				log.Fatal(err)
			}
			// questionId, err := re.LastInsertId()
			// if err != nil {
			// 	log.Fatal(err)
			// }
		})

	}

}

func main() {
	ExampleScrape()
}
