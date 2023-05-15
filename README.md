# etzba

**etzba** is a performance tests tool for sql and api servers. the tool use several configuraiton files to run load tests to measure the service performance.

the word etzba used in biblical times as a measurement unit and the english translation from hebrew to etzba is finger. 

### prepare cli tool

this command will build the cli tool:

```sh
make go-build
```

you may move the file to `/usr/local/bin` or run it from `cli` dir as follow:

``` sh
cd cli/

./etz sql --workers=3 --config=secret.json --helpers=sql.csv --duration=1s
./etz api --workers=3 --config=secret.json --helpers=api.json
```

you may also use the `examples/` directory to run tests, deploy environment for testing, example of helpers file or build the binary.

### prepare sql service helpers file

the helpers file used by etzba to create sql queries and schedule the workers job. 

create a csv file and add it to the command arg `--helpers=file.csv` as follow:

```
command,table,constraint,columnref,values
SELECT,results,avg_duration BETWEEN 13.0 AND 15.0,,
SELECT,results,min_duration BETWEEN 50.0 AND 60.0,,
SELECT,results,total BETWEEN 100 AND 200,,
```

if you run specific queries that don't require values, constraint or columnref, you should leave empty (e.g. `,`) in the csv

### prepare api service helpers file

the helpers file will list an array of api requests that can be sent to a service url during the load test, similar to this:
```json
[
  {
    "method": "GET",
    "url": "http://localhost:8080/v1/results"
  },
  {
    "method": "POST",
    "url": "http://localhost:8080/v1/results",
    "payload": "[{\"type\":\"api\",\"job_duration\":65.65,\"avg_duration\":12.32,\"min_duration\":56.32,\"med_duration\":31.14,\"max_duration\":99.9,\"total\":10},{\"type\":\"api\",\"job_duration\":45.45,\"avg_duration\":11.12,\"min_duration\":49.19,\"med_duration\":32.34,\"max_duration\":90.91,\"total\":21}]"
  }
]
```

### create authentication file

for api, you can define authentication by token and bearer or api key:

```json
{
  "apiAuth": {
    "method": "Bearer",
    "token": "XVlBzgbaiCMRAjWwhTHctcuAxhxKQFDa"
  }
}
```

and to authenticate for sql server, create the following json file:
``` json
{
  "sqlAuth": {
    "host": "localhost",
    "port": 5432,
    "database": "etzba",
    "user": "etzba",
    "password": "Pass1234"
  }
}
```

### execute etz command

now we are ready to use the helpers and secret files:

```sh
etz sql --workers=3 --config=secret.json --helpers=sql.csv --duration=1s
```

