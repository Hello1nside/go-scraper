package main

import (
	"database/sql"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	_ "github.com/go-sql-driver/mysql"
	//"io/ioutil"
	"bufio"
	"os"
	"log"
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func readFile() {
    file, err := os.Open("urls.txt")
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        //fmt.Println(scanner.Text())
	parseUrl(scanner.Text())
    }

    if err := scanner.Err(); err != nil {
        log.Fatal(err)
    }

}

func parseUrl(url string) {

	fmt.Println("Request: " + url)
	doc, err := goquery.NewDocument(url)
	checkErr(err)

	doc.Find(".single-product-wrapper").Each(func(index int, item *goquery.Selection) {
		productName := item.Find(".product_title").Text()
		productDescription := item.Find(".woocommerce-product-details__short-description").Text()

		imageTag := item.Find("img")
		image, _ := imageTag.Attr("src")

		// fmt.Printf("#%d - Product name: %s \n Description: %s \n Image Link: %s \n URL: %s", index, productName, productDescription, image, url)

		insert(productName, productDescription, image, url)
	})
}

func insert(name, descr, image, url string) {
	db, err := sql.Open("mysql", "root:password@/dbname?charset=utf8")
	checkErr(err)

	// insert
	stmt, err := db.Prepare("INSERT newtechstore SET name=?, description=?, image=?, url=?")
	checkErr(err)

	res, err := stmt.Exec(name, descr, image, url)
	checkErr(err)

	id, err := res.LastInsertId()
	fmt.Println(id)

	db.Close()
}

func main() {
	readFile()
}
