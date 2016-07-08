package main

import (
	"fmt"
	"strings"

	"gopkg.in/alecthomas/kingpin.v2"

	redis "gopkg.in/redis.v4"
)

var (
	server       = kingpin.Flag("server", "IP or DNS of redis server (e.g. 'redis.example.com')").Default("localhost").String()
	port         = kingpin.Flag("port", "redis port to connect to").Default("6379").Int()
	redisCommand = kingpin.Flag("command", "redis command").Default("nothing").String()
	doKeys       = kingpin.Flag("keys", "key or keys to do stuff with").Default("default").String()
	verbose      = kingpin.Flag("verbose", "Produce verbose output").Short('v').Default("false").Bool()
)

// stringInSlice a rudimentary contains function
func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// Connect creates redis client connection
func Connect() *redis.Client {
	addr := fmt.Sprint(*server, ":", *port)
	// fmt.Printf("connection string is %s\n", addr)
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	pong, err := client.Ping().Result()
	fmt.Println(pong, err)
	return client
}

// add keys to a slice
func keysToSlice(s string) []string {
	var sKeys = []string{}
	if strings.Contains(s, ",") == true {
		sKeys = strings.Split(s, ",")
	} else {
		sKeys = []string{s}
	}
	return sKeys
}

// collectKeys collects keys matching patterns
func collectKeys(sKeys []string) []string {
	client := Connect()
	var returnKeys = []string{}
	iter := client.Scan(0, "", 0).Iterator()
	for iter.Next() {
		// fmt.Println(iter.Val())
		if stringInSlice(iter.Val(), sKeys) == true {
			returnKeys = append(returnKeys, iter.Val())
		}
	}
	if err := iter.Err(); err != nil {
		panic(err)
	}
	fmt.Println(returnKeys)
	fmt.Println("Found", len(returnKeys), "keys")
	return returnKeys
}

// deleteKeys just deletes keys and returns bool if it deleted anything
func deleteKeys(sKeys []string) bool {
	client := Connect()
	for key := range sKeys {
		client.Del(sKeys[key])
		fmt.Println(key)
		return true
	}
	return false
}

// deleteKeys just deletes keys and returns bool if it deleted anything
func createKeys(sKeys []string) bool {
	client := Connect()
	for key := range sKeys {
		client.Set(sKeys[key], "test", 0)
		fmt.Println(key)
		return true
	}
	return false
}

func main() {
	kingpin.UsageTemplate(kingpin.CompactUsageTemplate).Version("0.0.1").Author("Benjamin Rizkowsky")
	kingpin.CommandLine.Help = "A simple redis tool."
	kingpin.Parse()
	fmt.Println("Command:", *redisCommand)
	switch {
	case *redisCommand == "delete":
		deleteKeys(collectKeys(keysToSlice(*doKeys)))
		fmt.Println("test:", *redisCommand)
	case *redisCommand == "test":
		createKeys(keysToSlice(*doKeys))
		collectKeys(keysToSlice(*doKeys))
	case *redisCommand == "nothing":
		fmt.Println("doing nothing.")
	}

}
