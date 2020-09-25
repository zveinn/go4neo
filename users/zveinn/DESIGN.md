# Design
in here we will make notes on cypher designs.



# example
```
MATCH (p:Person)-[:WORKS_FOR]->(c:Company)
WHERE c.name STARTS WITH "Face"
AND EXISTS {
  MATCH (p)-[:LIKES]->(t:Technology)
  WHERE size((t)<-[:LIKES]-()) >= 3
}
RETURN p.name as p, c.name AS company;
```

# Pattern suggestions from zveinn
```
Quiery := new(go4neo.Query)

Query.MATCH().NODE(tagx, Person{}).REL(tagy, WorksFor{}).NODE(tagz, Company{})

Query.WHERE("tagx.name = ?", $name)
Query.STARTSWITH("tagx.name", $name)

Query.STARTBLOCK()
    Query.MATCH().NODE(tagx, nil).REL(tagc, Likes{}).NODE(tagz, Technology{})
    Query.WHERE().SIZE("tagz >= ?", 3).REL(tagc, nil).ANY()
Query.ENDBLOCK()

Query.RETURN("tagx.name as name").RETURN("tagz.name as companyName")
Query.Commit()


Query.ORDERBY(x, "desc).LIMIT(100).OFFSET(0)
```