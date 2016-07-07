package main

import (
	"fmt"
	"strings"

	"gopkg.in/alecthomas/kingpin.v2"

	redis "gopkg.in/redis.v4"
)

var (
	server       = kingpin.Flag("server", "IP or DNS of redis server (e.g. 'redis.nitroplatformdev.com')").Default("localhost").String()
	port         = kingpin.Flag("port", "redis port to connect to").Default("6379").Int()
	redisCommand = kingpin.Flag("command", "redis command").Default("delete").String()
	doKeys       = kingpin.Flag("keys", "key or keys to do stuff with").Default("default").String()
	verbose      = kingpin.Flag("verbose", "Produce verbose output").Short('v').Default("false").Bool()
)

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

// Set the value of keysToSlice
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
func collectKeys() []string {
	client := Connect()
	var returnKeys = []string{}
	iter := client.Scan(0, "user.profile.id.16411652.v1", 0).Iterator()
	for iter.Next() {
		// fmt.Println(iter.Val())
		returnKeys = append(returnKeys, iter.Val())
	}
	if err := iter.Err(); err != nil {
		panic(err)
	}
	// fmt.Println(returnKeys)
	fmt.Println(len(returnKeys))
	return returnKeys
}

func ExampleClient_Scan() {
	client := Connect()
	var n int
	for {
		var cursor int64
		var keyss []string
		var err error
		cursor, keyss, err = client.Scan(cursor, "account.contacts.16408696", 10).Result()
		if err != nil {
			panic(err)
		}
		n += len(keyss)
		if cursor == 0 {
			break
		}
	}

	fmt.Printf("found %d keys\n", n)
	// Output: found 33 keys
}

func main() {
	kingpin.UsageTemplate(kingpin.CompactUsageTemplate).Version("0.0.1").Author("Benjamin Rizkowsky")
	kingpin.CommandLine.Help = "A simple redis tool."
	kingpin.Parse()
	// collectKeys()
	ExampleClient_Scan()
	// keys1 := keysToSlice("test1,test2,test3")
	// keys2 := keysToSlice(*doKeys)
	// for i := range keys2 {
	// 	fmt.Println(keys2[i])
	// }
	// for i := range keys1 {
	// 	fmt.Println(keys1[i])
	// }

}
