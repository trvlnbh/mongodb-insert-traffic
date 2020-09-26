## Run

```bash
$ go build -o ./mtf 
```

First, please create a database named `test_db` in MongoDB.
You must create collections before inserting the data. Collection is created with the command below. For cases 1 and 2, collections are sharded. Therefore, you must enable `test_db` to sharding before run the command below.

```bash
$ mtf setup --type=<case> --mongodb-uri="<MongoDB Connection URI>"
```

The command below inserts data until you exit. The speed of insertion may vary depending on performance and server status. Every minute, the number of inserted data is printed to standard output.

```bash
$ mtf insert --type=<case> --mongodb-uri="<MongoDB Connection URI>"
```





### CoordData (case 1)

Collection Sharding

* "key": {"_id": "hashed"}
* "numInitialChunks": 8

```json
{
    "timestamp": 0000000000000,
    "date": "Date",
    "first_name": "FirstName",
    "last_name": "LastName",
    "coordinate": {
        "latitude": 000.00000000,
        "longitude": 000.00000000
    }
}
```

| key        | type    | description           |
| ---------- | ------- | --------------------- |
| timestamp  | int64   | Millisecond Timestamp |
| date       | string  | RFC3339Nano           |
| first_name | string  | fake.FirstName()      |
| last_name  | string  | fake.LastName()       |
| latitude   | float32 | fake.Latitude()       |
| longitude  | float32 | fake.Longitude()      |





### FavorColorData (case 2)

Collection Sharding

* "key": {"_id": "hashed"}
* "numInitialChunks": 8

```json
{
    "timestamp": 0000000000000,
    "date": "Date",
    "personal_info": {
        "city": "City",
        "last_name": "LastName",
        "gender": "Gender"
    },
    "color": "Color"
}
```

| key       | type   | description           |
| --------- | ------ | --------------------- |
| timestamp | int64  | Millisecond Timestamp |
| date      | string | RFC3339Nano           |
| city      | string | fake.City()           |
| last_name | string | fake.LastName()       |
| gender    | string | fake.Gender()         |
| color     | string | fake.Color()          |





### EmailData (case 3)

No Sharding

```json
{
    "timestamp": 0000000000000,
    "date": "Date",
    "user_name": "UserName",
    "email": "Email",
    "text": "Text"
}
```

| key       | type   | description           |
| --------- | ------ | --------------------- |
| timestamp | int64  | Millisecond Timestamp |
| date      | string | RFC3339Nano           |
| user_name | string | fake.UserName()       |
| email     | string | fake.EmailAddress()   |
| text      | string | fake.Sentence()       |





## Reference

* https://godoc.org/github.com/icrowley/fake