# examples

example london query:

    g.V('LondonR').As('source').In('follows').As('target').All()


getting members that belong to both meetups:

    g.V("docker-london").As("source").In("follows").Has("follows", "kubernetes-london").As("target").All()


getting intersections:

    var scFollows = g.V('docker-london').As('source').In('follows').Has("follows", "kubernetes-london").As('target').All()
    var klFollows = g.V('kubernetes-london').As('source').In('follows').Has("follows", "docker-london").As('target').All()

    scFollows.Intersect(klFollows)

getting meetups:

     g.V('meetup').As('source').In('kind').As('target').All()

function example:


```javascript
function getFollowers(name, intersection) {
	return g.V(name).As("source").In("follows").Has("follows", intersection).As("target").All()
}

var scFollows = getFollowers("kubernetes-london", "docker-london")
var klFollows = getFollowers("docker-london", "kubernetes-london")

scFollows.Intersect(klFollows)
```


```javascript

function getFollowers(name, intersection) {
	return g.V(name).As("source").In("follows").Has("follows", intersection).As("target").All()
}

function prepareIntersections(x, y) {
   xFollows =  getFollowers(x, y)
   yFollows = getFollowers(y, x)
   
   return xFollows, yFollows
}

var x, y = prepareIntersections("kubernetes-london", "docker-london")

x.Intersect(y)

```

http://www.meetup.com/meetup_api/docs/2/members/

https://secure.meetup.com/meetup_api/console/?path=/2/members

https://github.com/google/cayley/wiki/Cayley-Go-API-(as-a-Library-in-your-own-projects)

https://github.com/google/cayley/wiki/Client-APIs

http://google-opensource.blogspot.co.uk/2014/06/cayley-graphs-in-go.html