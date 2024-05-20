## Pharmacy

Сборка проекта
```
mkdir build
cd build
go build ../cmd/pharmacy/main.go
```

Пример запуска

```
env PGCONN="host=127.0.0.1 port=5432 database=pharmacy user=gopher password=pass" ./main
```