package basic

import (
	"github.com/astaxie/beego"
	"log"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/go-redis/redis"
	"strconv"
)

type connectionsCache struct {
	RedisClient *redis.Client
	init bool
}

var cCache = new(connectionsCache)

func GetRedisClient() *redis.Client {
	if cCache.init == false {
		db := getRedisDbNum()

		client := redis.NewClient(&redis.Options{
			Addr:     beego.AppConfig.String("redis_addr"),
			Password: beego.AppConfig.String("redis_password"),
			DB:       db,
		})

		_, err := client.Ping().Result()
		if err != nil {
			log.Fatal(err)
			panic("Redis is unavailable")
		}

		cCache.RedisClient = client
		cCache.init = true
	}

	return cCache.RedisClient
}

func getRedisDbNum() int {
	db, err := strconv.Atoi(beego.AppConfig.String("redis_db"))
	if err != nil {
		log.Fatal(err)
		panic("Database number is wrong in the config file")
	}

	return db
}
