package main

import (
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type Students struct {
	Count int `json:"count"`
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())

	e.Any("/*", func(c echo.Context) error {
		return c.JSON(http.StatusOK, Students{Count: studentsCount()})
	})

	e.Logger.Fatal(e.Start(":" + port()))
}

func port() (port string) {
	port = "80"

	val, present := os.LookupEnv("PORT")
	if present {
		tPort, err := strconv.Atoi(val)
		if err == nil {
			port = strconv.Itoa(tPort)
		}
	}

	return
}

var httpClient = http.Client{
	Timeout: time.Second * 5,
}

var rgx = regexp.MustCompile(`Aktuálny počet podpisov: (\d+)`)

func studentsCount() int {
	res, err := httpClient.Get("https://www.zanasufiit.sk/wp-json/wp/v2/posts/120")
	if err != nil {
		return -1
	}

	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return -1
	}
	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return -1
	}
	text := gjson.Get(string(bodyBytes), "excerpt.rendered")

	rs := rgx.FindStringSubmatch(text.String())
	count, err := strconv.Atoi(rs[1])
	if err != nil {
		return -1
	}

	return count
}
