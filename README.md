# api-fosdem

"FOSDEM is the Free and Open source Software Developers' European Meeting, an event that is organized by the community for the community, as a non-for-profit." [Github](https://github.com/FOSDEM) - [Website](https://fosdem.org/)

**api-fosdem** attempts to provide a Free and Open Source API exposing the FOSDEM data, provided as XML.
Main goal of the project is to expand the possibility and ecosystem of the applications that now use a scraping of the pentabarf file, and provide other missing features.

Since the website of the FOSDEM never changed since the 2013 the information are available since then.

## Status

At the moment the API is in early stage, hosted in a free dyno on Heroku, and available at this URL:

> https://api-fosdem.herokuapp.com

Downtimes or a slow response dued to the automatic switch off of the dyno has to be expected.

## Endpoints

### /api/v1/speakers

Returns the full list of the speakers of the FOSDEM.
This endpoint is limited to return a maximum of 100 speakers, ordered by ID.
The `count` field can be used to understand how many documents matched the query. To iterate through the results use the parameters `offset` and `limit`.

```json
{
	"count": 2370,
	"data": [{
		"id": 4,
		"slug": "eben_moglen",
		"name": "Eben Moglen",
		"profile_image": "/2013/schedule/speaker/eben_moglen/96fb3c2694c07b996d58d26b1813f5a59d62c97b48cce52bf082ac8542e92d31.jpg",
		"profile_page": "/2013/schedule/speaker/eben_moglen/",
		"bio": "Eben Moglen is Executive Director of the Software Freedom Law Center and Professor of Law and Legal History at Columbia University Law School. He has represented many of the world's leading free software developers. Professor Moglen earned his PhD in History and law degree at Yale University during what he sometimes calls his “long, dark period” in New Haven. After law school he clerked for Judge Edward Weinfeld of the United States District Court in New York City and for Justice Thurgood Marshall of the United States Supreme Court. He has taught at Columbia Law School since 1987 and has held visiting appointments at Harvard University, Tel Aviv University and the University of Virginia. In 2003 he was given the Electronic Frontier Foundation's Pioneer Award for efforts on behalf of freedom in the electronic society. Professor Moglen is admitted to practice in the State of New York and before the United States Supreme Court.",
		"years": [2013]
	}, {
		"id": 6,
		"slug": "fosdem_staff",
		"name": "FOSDEM Staff",
		"profile_image": "/2018/schedule/speaker/fosdem_staff/c743944bdab7dce4a5b7d4696bd7264a8139bd4034a74033223fc9babf1c2d57.png",
		"profile_page": "/2018/schedule/speaker/fosdem_staff/",
		"years": [2013, 2014, 2015, 2016, 2017, 2018]
	}]
}
```

In order to *search* through the speakers the `slug` and `year` parameter can be used.

The `slug` is used in a regex over the slug field of the speaker (splitting the words).
The `year` is used to find a speaker that was present in the specified year. Multiple years can be specified (comma separataed).

#### examples:
- https://api-fosdem.herokuapp.com/api/v1/speakers?slug=oy$&year=2018

will return the speakers with a *slug* ending in "oy" that did a talk in the "2013".

- https://api-fosdem.herokuapp.com/api/v1/speakers?slug=ain&year=2015,2018

will return the speakers with a *slug* containing "ain" that did a talk in the "2015" OR in the "2018".


### /api/v1/speakers/{id}

Returns the details of the specified speaker.

- https://api-fosdem.herokuapp.com/api/v1/speakers/2072

```json
{
	"id": 2072,
	"slug": "francesc_campoy",
	"name": "Francesc Campoy",
	"profile_image": "/2018/schedule/speaker/francesc_campoy/cac4fd830f6d7dd839e1a8cd77ad17c9f5ba9bb39b9c2bc44b05f4568a72a1b6.jpg",
	"profile_page": "/2018/schedule/speaker/francesc_campoy/",
	"bio": "VP of Developer Relations at source{d} - Previously Google and Amadeus\n\nFrancesc Campoy Flores is the VP of Developer Relations at source{d}, a startup applying ML to source code and building the platform for the future of developer tooling. Previously, he worked at Google as a Developer Advocate for Google Cloud Platform and the Go team.\n\nHe’s passionate about programming and programmers, specially Go and gophers. As part of his effort to help those learning he’s given many talks and workshops at conferences like Google I/O, Gophercon(s), GOTO, or OSCON.\n\nWhen he’s not on stage he’s probably coding, writing blog posts, or working on his justforfunc YouTube series where he hacks while cracking bad jokes.",
	"links": [{
		"url": "http://twitter.com/francesc",
		"title": "twitter"
	}, {
		"url": "http://campoy.cat",
		"title": "personal page"
	}, {
		"url": "http://github.com/campoy",
		"title": "github"
	}, {
		"url": "http://justforfunc.com",
		"title": "justforfunc"
	}],
	"years": [2016, 2017, 2018]
}
```

