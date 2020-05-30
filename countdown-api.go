package main

import (
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
	e.Use(middleware.CORS())

	e.Any("/*", func(c echo.Context) error {
		count, err := studentsCount()

		if err != nil {
			return c.String(http.StatusServiceUnavailable, "error fetching data")
		}

		return c.JSON(http.StatusOK, Students{Count: count})
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

func studentsCount() (int, error) {
	res, err := httpClient.Get("https://m.facebook.com/story.php?story_fbid=133952524762076&id=104098081080854")
	if err != nil {
		return -1, err
	}

	if res.StatusCode != http.StatusOK {
		return -1, err
	}
	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return -1, err
	}

	err = res.Body.Close()
	if err != nil {
		return -1, err
	}

	rs := rgx.FindStringSubmatch(string(bodyBytes))
	count, err := strconv.Atoi(rs[1])
	if err != nil {
		return -1, err
	}

	return count, nil
}
