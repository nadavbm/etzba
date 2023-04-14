# etzba

the word etzba used in biblical times as a measurement unit and the english translation from hebrew to etzba is finger. 

### prepare cli tool

this command will build the cli tool:

```sh
go-build
```

you may move the file to `/usr/local/bin` or run it from `cd cli/` as follow:

``` sh
cd cli/

./etz sql --file "../examples/somefile.csv" --workers 100
```

optionally you can add `--verbose` for logs in stdout:

```sh
./etz sql --file "../examples/somefile.csv" --workers 100 --verbose
```
