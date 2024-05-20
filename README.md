## Pharmacy

Возможно придется установить

```
sudo apt-get install libx11-dev
sudo apt-get install libxcursor-dev
sudo apt-get install libgl1-mesa-dev xorg-dev
```

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