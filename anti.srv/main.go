package main

import core "git.ziniao.com/webscraper/go-gin-http"

func main() {
	core.NewApp(&core.AppConfigSt{Host: ":8878"}).RegHandler(Router).Start()
}
