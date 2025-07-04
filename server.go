package main

import (
	"fmt"
	database "mymodule/db"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func getHostname() (string, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return "", err
	}

	return hostname, nil
}

func setupRouter() *gin.Engine {
	r := gin.Default()

	// set up CORS
	// config := cors.DefaultConfig()
	// frontendOrigin := os.Getenv("FRONTEND")
	// config.AllowOrigins = []string{frontendOrigin}
	// r.Use(cors.New(config))

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong 3")
	})

	// GET /counter
	r.GET("/counter", func(c *gin.Context) {
		hostname, err := getHostname()
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"data":     nil,
				"hostname": hostname,
				"message":  "error in get hostname server",
			})
		}

		// get counter from db
		row := database.Db.QueryRow(`SELECT count FROM counters WHERE id = $1;`, 1)
		counter := 0
		err = row.Scan(&counter)
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"data":     nil,
				"hostname": hostname,
				"message":  "couldn't get counter from db",
			})
			return
		}

		resp := gin.H{
			"data": gin.H{
				"counter": counter,
			},
			"hostname": hostname,
			"status":   "success",
		}
		c.JSON(http.StatusOK, resp)
	})

	// POST /counter
	r.POST("/counter", func(c *gin.Context) {
		hostname, err := getHostname()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// + 1 counter in db
		_, err = database.Db.Exec(`UPDATE counters SET count = count + 1 WHERE id = $1 RETURNING count;`, 1)
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"data":     nil,
				"hostname": hostname,
				"message":  "error in increment counter",
			})
			return
		}

		// get updated counter
		row := database.Db.QueryRow(`SELECT count FROM counters WHERE id = $1;`, 1)
		counter := 0
		err = row.Scan(&counter)
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"data":     nil,
				"hostname": hostname,
				"message":  "couldn't get counter from db",
			})
			return
		}

		// counter := 0
		resp := gin.H{
			"data": gin.H{
				"counter": counter,
			},
			"hostname": hostname,
			"message":  "success",
		}
		c.JSON(http.StatusOK, resp)
	})

	return r
}

func main() {
	r := setupRouter()

	database.ConnectDatabase()

	appUrl := os.Getenv("APP_URL")
	fmt.Printf("Server listening on %s...\n", appUrl)
	r.Run(appUrl)
}
