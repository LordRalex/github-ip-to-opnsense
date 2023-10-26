package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
)

func main() {
	e := gin.New()
	e.GET("", getIPs)

	listener, err := net.Listen("tcp", "0.0.0.0:8080")
	if err != nil {
		panic(err)
	}

	err = e.RunListener(listener)
	if err != nil {
		panic(err)
	}
}

func getIPs(c *gin.Context) {
	meta, err := getGithubMeta()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	ips := append(meta.Web, meta.Api...)
	ips = append(ips, meta.Git...)
	ips = append(ips, meta.Packages...)
	ips = append(ips, meta.Pages...)
	ips = append(ips, meta.Actions...)
	ips = append(ips, meta.Web...)

	for _, v := range ips {
		ip, _, err := net.ParseCIDR(v)
		if err != nil {
			continue
		}
		if ip.To4() != nil {
			c.Writer.WriteString(v + "\r\n")
		}
	}
}

func getGithubMeta() (data GithubResponse, err error) {
	var response *http.Response
	response, err = http.Get("https://api.github.com/meta")
	if err != nil {
		return
	}
	defer response.Body.Close()

	err = json.NewDecoder(response.Body).Decode(&data)
	return
}

type GithubResponse struct {
	Web      []string
	Api      []string
	Git      []string
	Packages []string
	Pages    []string
	Actions  []string
}
