package main

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"syscall"

	"time"

	"github.com/gin-gonic/contrib/cors"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/urfave/cli/v2"
	"golang.org/x/term"

	"github.com/golbeng-original/chomakers-web/apis"
	"github.com/golbeng-original/chomakers-web/models"
)

// Allow-Control-Allow-Origin을 해결하기 위한 방법
// 커스텀한 미들웨어

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, Set-Cookie")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		/*
			if c.Request.Method == "OPTIONS" {
				c.AbortWithStatus(204)
				return
			}
		*/

		c.Next()
	}
}

func serveMiddleware(localFileSystem static.ServeFileSystem) gin.HandlerFunc {
	fileserver := http.FileServer(http.Dir("./view"))

	fileRegex, _ := regexp.Compile("^/(manager|potofolio|essay|about|edit|new)/?.?")

	return func(c *gin.Context) {

		if c.Request.URL.Path == "/" || fileRegex.MatchString(c.Request.URL.Path) {

			fmt.Println("matched!!! [" + c.Request.URL.Path + "]")

			c.Request.URL.Path = "/"
			fileserver.ServeHTTP(c.Writer, c.Request)
			c.Abort()
			return
		}

		c.Next()
	}
}

// Jwt Token을 확인하는 영역
// Refresh Token도 갱신
// Access Token 갱신
func vertifyTokenMiddleware(repoConfigure *models.RepositoryConfigure) gin.HandlerFunc {

	return func(c *gin.Context) {

		if c.Request.URL.Path == "/api/login" {
			c.Next()
			return
		}

		if c.Request.Method == "GET" || c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}

		isAuthentication, err := apis.CheckAuthentication(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, apis.FailedResponsePreset(err.Error()))
			c.Abort()
			return
		}

		if !isAuthentication {
			c.JSON(http.StatusUnauthorized, apis.FailedResponsePreset(""))
			c.Abort()
		}

		c.Next()
	}
}

func Setup(repoConfigure *models.RepositoryConfigure, imagePath string) *gin.Engine {

	router := gin.Default()

	//router.LoadHTMLFiles("", "manager", "potofolio", "essay", "about")

	//localFileSystem := static.LocalFile("./assets/images", true)
	//router.Use(serveMiddleware(localFileSystem))
	//router.Use(static.Serve("/assets", static.LocalFile("./assets", true)))

	//router.Use(static.Serve("/images", localFileSystem))

	// Cors 허용 여부
	//router.Use(corsMiddleware())

	corHandler := cors.New(cors.Config{
		//AllowAllOrigins: true,
		AllowedOrigins:   []string{"http://chomakers.com", "http://www.chomakers.com"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Origin", "Cookie"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
		//AllowOriginFunc: func(origin string) bool {
		//	fmt.Println(origin)
		//	return true
		//},
	})
	router.Use(corHandler)

	if repoConfigure.IsCheckAuthorize {
		router.Use(vertifyTokenMiddleware(repoConfigure))
	}

	router.Static("/images", imagePath)

	// api 등록 구간
	api := router.Group("api")

	apis.LoginApis(api, repoConfigure)
	apis.PotofolioApis(api, repoConfigure)
	apis.EssayApis(api, repoConfigure)
	apis.AboutApis(api, repoConfigure)

	return router
}

func RunWebServer(port string) {
	dbConnection := models.DBConnection{}
	dbConnection.Open("./assets/data.db")
	defer dbConnection.Close()

	repositoryConfigure := &models.RepositoryConfigure{}
	repositoryConfigure.Init(&dbConnection)

	Setup(repositoryConfigure, "./assets/images").Run(port)
}

func CreateUser(username, password string) error {
	dbConnection := models.DBConnection{}
	dbConnection.Open("./assets/data.db")
	defer dbConnection.Close()

	userRepository := &models.UserRespository{DBConnect: &dbConnection}
	err := userRepository.CreateTable()
	if err != nil {
		return err
	}

	err = userRepository.AddUser(username, password)
	if err != nil {
		return err
	}

	return nil
}

func main() {

	port := ":8081"

	app := &cli.App{
		Name:     "qudghweb",
		Version:  "v0.0.1",
		Compiled: time.Now(),
		Commands: []*cli.Command{
			{
				Name: "create-super",
				Action: func(c *cli.Context) error {

					var username string
					fmt.Print("username :")
					fmt.Scan(&username)

					fmt.Print("passsword :")
					password, _ := term.ReadPassword(int(syscall.Stdin))
					fmt.Println()

					fmt.Print("password Confirm :")
					passwordConfirm, _ := term.ReadPassword(int(syscall.Stdin))
					fmt.Println()

					if string(password) != string(passwordConfirm) {
						fmt.Println("password not same")
						return fmt.Errorf("password not same")
					}

					err := CreateUser(username, string(password))
					if err != nil {
						fmt.Println(err.Error())
						return err
					}

					fmt.Println("add user success")

					return nil
				},
			},
		},
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:  "port",
				Usage: "server port",
				Value: 8081,
			},
		},
	}

	app.Action = func(c *cli.Context) error {
		port = fmt.Sprintf(":%d", c.Int("port"))
		RunWebServer(port)
		return nil
	}

	app.Run(os.Args)
}
