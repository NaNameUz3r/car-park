### Car park service

This is a learning project from the skillsmart coding school.

#### Current state:

Poor proof of concept of a multi-user (management) service with CRUD for multi-models implementing a kind of car-park for contractor companies.

We can CRUD Vehicles, and fetch their rides with geotracks and some reports.

####

Check this out with fixture dump:

```
docker-compose up
```
And log in with "admin2" username, and "qwerty" password.
There is, actually admin1 with same password. Explore :D

The dump contains a bunch of rides for 2023-2024 years for vehicles with ID's  25184,25185 and 25186

Currently, to be able to view ride routes on the map you should have here.com api token and pass it in docker-compose as env "HERE_API_KEY"