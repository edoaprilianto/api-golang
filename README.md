# Golang Food Data REST API 🚀 

🔥 Golang Rest Api with basic JWT Authentication

## Development Setup
1. rename .example.env to .env
2. Running "docker-compose up"
Noted : before running this command you must have **docker** and **docker-compose** installed in your system.

3. if you have error , please restart container or you can follow script below :
   * docker-compose down
   * docker-compose up
   


## SERVICE DOCKER
1. API -> http://localhost:9090
2. Admin Mysql UI -> http://localhost:8282
3. OpenApi (Swagger) -> http://localhost/

## Technology
- Language (golang)
- Database (mysql)

### Libraries
- Router (gorilla/mux)
- Server (net/http)
- JWT (dgrijalva/go-jwt)<!-- - Password Encryption (bcrypt) -->
- Database ORM (gorm) 
- Live Reload (cosmtrek/air)


## API Documentation

### Authentication
> **POST** ``/auth/login``

Login with Phone and password.

##### Body

```json
{
    "Phone": "12345",
    "Password": "VQXe",
}
```

#### Output

```json
{
    "User": {
        "ID": 1,
        "CreatedAt": "2022-04-24T09:36:26.192Z",
        "UpdatedAt": "2022-04-24T09:36:26.192Z",
        "DeletedAt": null,
        "Name": "edo aprilianto",
        "Phone": "12345",
        "Role": "Admin",
        "Password": "VQXe"
    },
    "Token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiZWRvIGFwcmlsaWFudG8iLCJwaG9uZSI6IjEyMzQ1Iiwicm9sZSI6IkFkbWluIiwiaWQiOjEsImV4cCI6MTY1MTM5Nzk0Mn0.yWlIllgzO3xUp2Vw-ivovZ3-ExfsYnxTb9xBm2diq3I"
}
```

> **POST** ``/auth/signup``

Create a new user in the database.

##### Body

```json
{
    "Name": "BURKHI",
    "Phone": "1234",
    "Role": "Admin",
}
```

#### Output

```json
{
    "User": {
        "ID": 1,
        "CreatedAt": "2022-04-24T09:36:26.192Z",
        "UpdatedAt": "2022-04-24T09:36:26.192Z",
        "DeletedAt": null,
        "Name": "BURKHI",
        "Phone": "12345",
        "Role": "Admin",
        "Password": "VQXe"
    },
    "Token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiZWRvIGFwcmlsaWFudG8iLCJwaG9uZSI6IiIsInJvbGUiOiJBZG1pbiIsImlkIjoxLCJleHAiOjE2NTEzOTc3ODZ9.UIu9bmcVvtIFkwsZ_cbSI3FUqFY7osPGjRhm4tZibLA"
}
```

> **GET** ``/auth/profile``

##### Body
using **jwt** as ``Authorization``

get all data private claims by **JWT**

### OUTPUT
```json
{
    "name": "BURKHI",
    "phone": "12345",
    "role": "Admin",
    "id": 1,
    "exp": 1651397942
}
```


### JSON Manipulation

All endpoints are protected, must send valid **jwt** as ``Authorization`` header with each request.

> **GET** &nbsp; ``fetch/resources``

Get All resources from https://stein.efishery.com/v1/storages/5e1edf521073e315924ceab4/list
and new field "Price_USD". 
iam using https://api.exchangerate.host/latest?symbols=USD for API converted


#### Output

```json
[
    {
    "uuid": "a46247af-eb81-48c8-98cc-f26f5829a135",
    "komoditas": "dambaba",
    "area_provinsi": "BALI",
    "area_kota": "BULELENG",
    "size": "70",
    "price": "500000",
    "tgl_parsed": "2022-03-10 12:46:54.378",
    "timestamp": "1646891214",
    "price_usd": 539918.5
    }
    {
    "uuid": "9070ff73-9d87-4d10-8cc4-38cb845285ce",
    "komoditas": "dambaba",
    "area_provinsi": "JAWA BARAT",
    "area_kota": "CIREBON",
    "size": "80",
    "price": "434343",
    "tgl_parsed": "2022-03-10T12:58:47.641+0700",
    "timestamp": "1646891927",
    "price_usd": 469019.642091
    }
]
```

> **GET** &nbsp; ``/fetch/aggregate``

Get All resources from https://stein.efishery.com/v1/storages/5e1edf521073e315924ceab4/list
and aggreagate data by Area_propinsi and showing Total , Avg , Max , Min, Median

#### Output

```json
[ 
    {
    "dd": {
        "area_provinsi": "dd",
        "Price": {
            "total": 1000000,
            "avg": 1000000,
            "min": 1000000,
            "max": 1000000,
            "median": 1000000
        }
    },
    "jawa timur tengah": {
        "area_provinsi": "jawa timur tengah",
        "Price": {
            "total": 15000,
            "avg": 15000,
            "min": 15000,
            "max": 15000,
            "median": 15000
        }
    }
]
```


> **GET** &nbsp; ``/fetch/profile``

```json
{
    "name": "BURKHI",
    "phone": "081209109",
    "role": "admin",
    "id": 18,
    "exp": 1651249570
}
```

## Context
![Alt text](Context.png?raw=true "Context")


NOTED : IF U HAVE TROUBLE "CORS" OR FAILED TO FETCH WHEN USING OPEN API, PLEASE DISABLE CORS IN YOUR BROWSER.

Thank u and i hope you happy :)