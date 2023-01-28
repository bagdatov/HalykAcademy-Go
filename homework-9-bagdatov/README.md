# homework-9-starter

Условия для домашки:
1) заменить структуру хранения данных о пользователях и крипто кошельках на БД (mysql, postgresql)
2) `*` сделать все это через слои - https://github.com/bxcodec/go-clean-arch

--------

# How implemented:

## PostgreSQL: installment (Linux)
```bash
# download
sudo apt update
sudo apt -y install postgresql

# set up
sudo -i -u postgres
psql
\q

# log in
sudo -u postgres psql -c "ALTER USER postgres PASSWORD 'postgres';"

# create database
sudo -u postgres psql -c "CREATE DATABASE postgres;"

# start server 
# default port is 5432
sudo service postgresql start
sudo -u postgres psql

# create tables
\c postgres
\i init.pgsql 
```

There is an implementation of clean architecture, however I am not sure if I understood correctly.<br>
Here is the database [schema](init.pgsql). <br>
Business logic [here](domain/interfaces.go). <br>
Also I've prepared simple bash [script](script.sh) for requests.<br>

--------

Доступные сущности

Необходимо написать веб сервис, используя стандартный веб сервер из пакета
net/http который умеет выполнять следующие действия

* Для каждого запроса проверять basic авторизацию (заголовок Authorization). 
  Если юзера с таким логином/паролем нет, то выдать ошибку 403.
* Каждый ответ должен содержать заголовок "execution", который содержит
  время обработки запроса.
* У каждого юзера есть свои криптокошельки, которые он может создавать/уничтожать.
* Каждый кошелек имеет метод старта майнинга и остановки. 
* Должны быть реализованы следующие методы:

1. GET /app/user/$id - метод возвращает username пользователя и имена всех 
   криптокошельков у этого юзера по ID. Любой авторизованный пользователь
   может запросить информацию о других пользователях.
2. POST /app/user/$id - метод, регистрирующий нового юзера. В POST параметрах
   передается "username" и "password". Если с данным ID или с данным username
   уже существует пользователь - вернуть код 400. Регистрировать других может
   любой авторизованный пользователь.
3. GET /app/wallet/$name - метод, возвращающий количество денег на кошельке.
   Не забываем валидировать авторизацию через заголовок и проверять какой
   именно пользователь запрашивает. Если же кошелек не относится к авторизованному
   пользователю, то вернуть ошибку 404.
4. POST /app/wallet/$name - метод, создающий новый крипто кошелек для авторизованного
   пользователя. Если кошелек с таким именем существует - вернуть ошибку 400.   
4. OPTION /app/wallet/$name/start - метод, запускающий майнинг крипты.
5. OPTION /app/wallet/$name/stop - метод, останавливающий майнинг крипты.

P.S. нельзя использовать глобальные переменные

Задание со звездочкой: Реализовать задание с использованием сервера fasthttp.

Удачи:)