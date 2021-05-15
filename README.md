# Marvel Character APIs

Marvel Character sevice which contains APIs and one job to handle Marvel's character data 

## branches

- [master](https://github.com/nmmugia/marvel-character) (stable)

## step to run with docker
1. Install redis on your local
2. copy .env.example file into .env, and filled the configuration with the intended value
3. for an easy installation, kindly install docker first and then run an executable file:
```bash
$ docker-compose up
```
## step to run without docker
1. Install redis, [go](https://golang.org/doc/install) (should be using v1.12.12>)
2. copy .env.example file into .env, and filled the configuration with the intended value
3. run the program using go command:
```bash
$ go mod download
$ go build -o main
$ go run main.go

```

## utilities
- Message(errx models.Errorx, data interface{}) (res map[string]interface{}) Message function is to build response per standard
- Response(w http.ResponseWriter, data interface{}, errx models.Errorx) Response function is to encode response
- StringToMD5(text string) (result string) StringToMD5 function is to convert/hash string to MD5 format
- StringToInt(value string, def int64) (result int64) StringToInt to casting string to integer, first param would be value, and the second one is default returned value
- ParseFromString(value string) (res time.Time, err error) Parse time to string
- HitMarvelsEndpoint(method string, path string, params string) (result models.MarvelsResult, err error) HitMarvelsEndpoint func is to hit marvel's endpoint based on method(GET, POST, etc), path, and params(should be "&key=value" format)
- IsInteger(v string) (result bool) to check whether the string is actually numeric

## externals

- [Robfig Cron](github.com/robfig/cron) To handle cron job easier 
- [GodotENV](github.com/joho/godotenv) Load ENV helper
- [SQLx](github.com/jmoiron/sqlx) Extension of standard sql library (usage in here is to use JSONText type)
- [Gorilla Mux](github.com/gorilla/mux) Make Routing Easier
- [Go Redis](github.com/go-redis/redis) Library to interfacing with redis
