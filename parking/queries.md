# create (p:Person {name:"sven"})-[:USES]->(l:Language {name:"golang"})
# match (p:Person {name:"sven"}) set p.name = "sveinn"
# MATCH (a:Person) WHERE a.name = 'sveinn' CREATE (a)-[:USES]->(l:Language {name:"javascript"})
# MATCH (a:Person) WHERE a.name = 'sveinn' CREATE (a)-[:USES {frequency:"rarely"}]->(l:Language {name:"javascript"}