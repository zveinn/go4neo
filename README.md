# go4neo
Neo4j interface written in golang.

# Information
 - Schedule: https://zhuffles.dev/blog/classes/
 - Stream: https://www.twitch.tv/zhuffles
 - Youtube: https://www.youtube.com/channel/UCW6eiMiVqYroPX1qiosAbnQ
 - Discord: https://discord.com/invite/r4wxkXd

# Materials
 - Introduction: https://www.youtube.com/watch?v=oRtVdXvtD3o&ab_channel=MicrosoftResearch
 - Installing: https://neo4j.com/download/
 - Getting started: https://neo4j.com/developer/go/
 - Driver: https://github.com/neo4j/neo4j-go-driver
 - Cypher: https://neo4j.com/developer/cypher/intro-cypher/
 - Cypher Manual: https://neo4j.com/docs/cypher-manual/current
 - Operations Manual: https://neo4j.com/docs/operations-manual/current/introduction/

 # Good to know
 - Delete all data: MATCH (n) DETACH DELETE n;

# Examples of Query operators
 - match X {name:dasdasd}
 - Match X with(Y) Z 
 - Match X with(Y) Z Where exists X with(D) A RETURN x,y,z
 - min, max - gets min or max value from data
 - Union - unifies two queries.. 
 - ORDER BY .. etc.
 - AS - sets something to a variable
 - collect - to get a list of results
 - set - used to change fields

# Advanced query
Using a block {} makes variables inside that block accessable .. SEE: t
```
MATCH (p:Person)-[:WORKS_FOR]->(c:Company)
WHERE c.name STARTS WITH "Face"
AND EXISTS {
  MATCH (p)-[:LIKES]->(t:Technology)
  WHERE size((t)<-[:LIKES]-()) >= 3
}
RETURN p.name as p, c.name AS company;
```

# Advanced Query combination with Union
The first query returns results for the next one
```
MATCH (p:Person)-[:LIKES]->(:Technology {type: "Java"})
RETURN p.name AS person, p.birthdate AS dob
ORDER BY dob DESC

UNION

MATCH (p:Person)
WHERE size((p)-[:IS_FRIENDS_WITH]->()) > 1
RETURN p.name AS person, p.birthdate AS dob
ORDER BY dob DESC;
```

# Super advanced 2x sets combined to list
```
CALL {
	MATCH (p:Person)-[:LIKES]->(:Technology {type: "Java"})
	RETURN p

	UNION

	MATCH (p:Person)
	WHERE size((p)-[:IS_FRIENDS_WITH]->()) > 1
	RETURN p
}
RETURN p.name AS person, p.birthdate AS dob
ORDER BY dob DESC;
```
and if you want to show the relationships too ...
```
CALL {
	MATCH (p:Person)-[:LIKES]->(:Technology {type: "Java"})
	RETURN p

	UNION

	MATCH (p:Person)
	WHERE size((p)-[:IS_FRIENDS_WITH]->()) > 1
	RETURN p
}
WITH p,
     [(p)-[:LIKES]->(t) | t.type] AS technologies,
     [(p)-[:IS_FRIENDS_WITH]->(f) | f.name] AS friends

RETURN p.name AS person, p.birthdate AS dob, technologies, friends
ORDER BY dob DESC;
```