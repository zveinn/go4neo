package main

import (
	"log"
	"os"
	"reflect"

	"github.com/neo4j/neo4j-go-driver/neo4j"
)

type Person struct {
	Name string `json:"name"`
}

var Driver neo4j.Driver

func main() {

	db := &NeoDB{}
	err := db.Connect("bolt://localhost:7687", "neo4j", "test", false)
	if err != nil {
		panic(err)
	}

	values := make(map[string]interface{})
	values["name"] = "Sveinn"

	builder := NewBuilder()

	// res, err := builder.CREATE().NODE("p", "Person", values).RUN()
	// if err != nil {
	// 	panic(err)
	// }
	// log.Println(res)

	data, _, err := builder.MATCH().NODE("p", "Person", values).COLLECT("p").Log().RUN()
	if err != nil {
		panic(err)
	}
	// log.Println(res)
	for _, v := range data.Nodes {
		log.Println(v.ID, v.Labels)
		for ii, iv := range v.Props {
			log.Println(ii, iv)
		}
	}

	os.Exit(1)

	// OLD STUFF

	// Driver, err = Connect("bolt://localhost:7687", "neo4j", "test", false)
	// if err != nil {
	// 	panic(err)
	// }
	// defer Driver.Close()

	// session, err := NewSession()
	// if err != nil {
	// 	panic(err)
	// }
	// defer session.Close()

	// values := make(map[string]interface{})
	// values["name"] = "Sveinn3"
	// values["id"] = "11"
	// values["type"] = "newb"

	// // _, err = QUERY(session, "CREATE", "Player", values, []string{"name", "id"})
	// // if err != nil {
	// // 	log.Println(err)
	// // 	os.Exit(1)
	// // }

	// res, err := QUERY(session, "MATCH", "Player", values, []string{"name", "id", "type"})
	// if err != nil {
	// 	log.Println(err)
	// 	os.Exit(1)
	// }
	// for _, record := range res {
	// 	log.Println(record.Keys(), record.Values())
	// }

	// List(session)
}

func Connect(uri, user, pass string, encrypted bool) (driver neo4j.Driver, err error) {
	driver, err = neo4j.NewDriver(uri, neo4j.BasicAuth(user, pass, ""), func(c *neo4j.Config) {
		c.Encrypted = encrypted
	})
	return
}
func NewSession() (s neo4j.Session, err error) {
	return Driver.Session(neo4j.AccessModeWrite)
}

func QUERY(session neo4j.Session, cypherPrefix, nodeName string, keyValues map[string]interface{}, returnKeys []string) ([]neo4j.Record, error) {

	data := []neo4j.Record{}
	lastIndex := len(keyValues) - 1
	count := 0
	query := cypherPrefix + " (x:" + nodeName
	if len(keyValues) > 0 {
		query += " { "
	}
	for i := range keyValues {
		if count == lastIndex {
			query += i + ": $" + i
		} else {
			query += i + ": $" + i + ", "
		}
		count++
	}
	if len(keyValues) > 0 {
		query += "})"
	} else {
		query += ")"
	}

	if len(returnKeys) > 0 {
		query += " RETURN "
	}
	lastIndex = len(returnKeys) - 1
	count = 0
	for _, v := range returnKeys {

		if count == lastIndex {
			query += "x." + v
		} else {
			query += "x." + v + ", "
		}
		count++

	}

	log.Println(query)
	// _, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
	result, err := session.Run(query, keyValues)
	if err != nil {
		return nil, err
	}

	re, err := result.Summary()
	log.Println(re.ResultAvailableAfter())
	if result.Next() {
		data = append(data, result.Record())
	}

	// return nil, result.Err()
	// })
	// if err != nil {
	// 	return nil, err
	// }

	return data, err
}
func List(session neo4j.Session) {

	result, err := session.Run(
		"MATCH (n: Player) RETURN collect(n)",
		map[string]interface{}{})
	if err != nil {
		os.Exit(1)
	}

	if result.Next() {
		log.Println(reflect.TypeOf(result.Record().GetByIndex(0)))
		log.Println(result.Record().GetByIndex(0))
		list := result.Record().GetByIndex(0).([]interface{})
		for i, v := range list {
			log.Println(i, v.(neo4j.Node).Props())
		}
	}

}
