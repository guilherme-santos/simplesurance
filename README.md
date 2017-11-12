# Simplesurance Challenge

Using only the standard library, create a Go HTTP server that on each request responds with a counter of the total number of requests that it has received during the previous 60 seconds (moving window). The server should continue to the return the correct numbers after restarting it, by persisting data to a file.

# Test
To run all unit tests just type:
```
make test
```

# Running
If you have Go installed in your machine, type:
```
make run
```
To run you need some envvars exported, by default our Makefile is exporting `.env` file. If you want change some configuration you can export `SIMPLEINSURANCE_API_PORT` and `SIMPLEINSURANCE_API_COUNTER_FILENAME`.

If you don't have Go installed, you can build a docker image, to do that just type:
```
make docker-image
```

To run this new docker image, type:
```
docker run --rm -p 8080:8080 -v `pwd`/dbdata:/dbdata --env-file .env guilherme-santos/simplesurance
```
Remember of update other parameters based on your configuration file filename and port
