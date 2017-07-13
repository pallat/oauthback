package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// Handler
func callback(c echo.Context) error {
	var v map[string]interface{}
	m := map[string]interface{}{}
	m["request"] = c.Request().RequestURI
	err := c.Bind(&v)
	m["error"] = err
	m["payload"] = v
	return c.JSON(http.StatusOK, m)
}

func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Routes
	e.GET("/", callback)
	e.GET("/login", handleGoogleLogin)
	e.GET("/oauth/callback", handleGoogleCallback)

	// Start server
	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))
	// e.Logger.Fatal(e.Start(":1323"))
}

var (
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  "https://oauthback.herokuapp.com/oauth/callback",
		ClientID:     os.Getenv("googlekey"),
		ClientSecret: os.Getenv("googlesecret"),
		Scopes: []string{
			"https://www.googleapis.com/auth/urlshortener",
			"https://www.googleapis.com/auth/userinfo.profile",
			"https://www.googleapis.com/auth/userinfo.email",
		},
		Endpoint: google.Endpoint,
	}
	// Some random string, random for each request
	oauthStateString = "random"
)

func handleGoogleCallback(c echo.Context) error {
	state := c.QueryParam("state")
	if state != oauthStateString {
		fmt.Printf("invalid oauth state, expected '%s', got '%s'\n", oauthStateString, state)
		return c.Redirect(http.StatusTemporaryRedirect, "/")
	}

	code := c.QueryParam("code")
	token, err := googleOauthConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		fmt.Printf("Code exchange failed with '%s'\n", err)
		return c.Redirect(http.StatusTemporaryRedirect, "/")
	}

	fmt.Println("accessToken", token.AccessToken)

	// response, err := http.Get("https://www.googleapis.com/oauth2/v1/tokeninfo?access_token=" + token.AccessToken)
	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	fmt.Println(string(contents))
	return c.String(200, string(contents)+`
*** validate token: "https://www.googleapis.com/oauth2/v2/userinfo?access_token="`+token.AccessToken)
}

func handleGoogleLogin(c echo.Context) error {
	url := googleOauthConfig.AuthCodeURL(oauthStateString)
	return c.Redirect(http.StatusTemporaryRedirect, url)
}
