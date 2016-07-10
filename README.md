# Sunny day, 4 PM

> I’m an algorithm! Making ice cream everyday. When you favourite an image I learn
> from it to generate the next one. Learning day by day from your favourites!
> :icecream:
> [Sunny day, 4 PM](https://twitter.com/sunnyday4pm)

## About

I’m a three-minute-a-day artist. Every day I will go on Twitter and get the most
favourited status from my account. From those I will create a new image.
Learning day by day — a new image posted every day at 4 PM. Creating our perfect
ice cream.  I’m interested in knowing how an ice cream makes you feel. What are
its dynamics? How do colours and shapes make it an ice cream? The program uses a
parameterised circle, rectangle and triangle. It learns and evolves their
placement, size, rotation and color by using a genetic algorithm. Using basic
shapes help us study the relationship between the ice cream’s components: cone,
scoop and chocolate flake. The first instalment is fully random. The program
then uses previous ice creams to create a new one each day. It assumes that the
more Twitter favourite an ice cream gets, the more its characteristics are
important. It picks parents. The first parent is the ice cream with the most
favourites. The second parent is randomly picked from all the ice creams that
received favourites. The features of the parents are then mixed at random (genes
crossover) to create the new ice cream DNA. Each shape has a random color picked
from a subset of all possible colours. These subsets evolve over generations.

[Read more on Medium](https://medium.com/@fhoehl/sunny-day-4-pm-16efbc33b4e7#.51ojusbsg)

## Quick setup

You will need:

* Go
* Python 2.7
* Redis

```bash
pip install -r requirements.txt
make
./bin/make/icecreamd
```

### With Docker

```bash
docker-compose build
docker-compose run botmachine icecreamd
```

Point your browser to: **localhost:2003**

You will see a gallery of ice cream. Reload the page to create a new one. Click
once or multiple times on an image to favourite it.

## Tools!

After building the bin folder will cointain several tools:

* **icecreamd**: a web server serving the ice cream web app
* **makeicecream**: a tool to generate new ice cream or get an SVG from the
  database
* **tweeticecream**: a tool to tweet!
* **syncicecream**: a tool that will query Twitter for the favourite count and
  update the database for every ice cream.
* **freezer2dot**: a tool to generate a Graphviz visualisation of the database.
* **freezer2json**: a tool to dump the database as a JSON file.

In order to tweet one must fill the .env file with the following environment
variables:

* TWITTER_CONSUMER_KEY
* TWITTER_CONSUMER_SECRET
* TWITTER_ACCESS_TOKEN
* TWITTER_ACCESS_TOKEN_SECRET

## License

Copyright (c) 2016 François Hoehl

[MIT License](http://en.wikipedia.org/wiki/MIT_License)

