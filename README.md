### Car park service

This is a learning project from the skillsmart coding school.

#### Current state:

Gin server at 8888 port, with one Vehicle model stored in PostgreSQL database.

##### View

View with vehicles table is accessible with GET method at `:8888/view/vehicles` URI.

There is no ui for adding new records service tight now, but basic api already implemented:


##### Api

- POST at `:8888/api/vehicles` with json like:

```
{
"description": "Very old cheap car",
"price": 150,
"mileage": 99999,
"manufactured": 1998
}
```

Will add new record in database.

- PUT `"/vehicles/:id` with same json structure will update record with specified ID.
- DELETE `"/vehicles/:id` will delete record with specified ID.


### Check this out:

1. `git clone https://github.com/NaNameUz3r/car-park.git`
2. `cd car-park && docker-compose up`
3. Go to `http://localhost:8888/view/vehicles` in your browser.
4. Play around with api, put some records, for example with curl:

```
curl --location --request POST 'localhost:8888/api/vehicles/' \
     --header 'Content-Type: application/json' \
     --data-raw '{
        "description": "Very old gazel, but it can carry some corpses",
        "price": 777,
        "mileage": 5555,
        "manufactured": 444
        }'
```
