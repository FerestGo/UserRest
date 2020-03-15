package main

var config Config

func main() {
	config = GetConfig()

	db := DB(config)
	defer db.Close()

	router := getRouter()
	StartServer(router, config)
}
