package main

import (
	"log"

	"github.com/neo4j/neo4j-go-driver/neo4j"
)

type DB interface {
	Connect(uri, user, pass string, encrypted bool) error
	NewSession() (neo4j.Session, error)
}

var D DB

type NeoDB struct {
	Driver neo4j.Driver
}

func (db *NeoDB) Connect(uri, user, pass string, encrypted bool) (err error) {
	db.Driver, err = neo4j.NewDriver(uri, neo4j.BasicAuth(user, pass, ""), func(c *neo4j.Config) {
		c.Encrypted = encrypted
	})
	D = db
	return
}

//TODO: might have to rewrite the return value at some point ...
func (db *NeoDB) NewSession() (s neo4j.Session, err error) {
	s, err = db.Driver.Session(neo4j.AccessModeWrite)
	if err != nil {
		panic(err)
	}
	return
}

func NewBuilder() *Builder {
	sess, err := D.NewSession()
	if err != nil {
		panic(err)
	}
	b := new(Builder)
	b.Session = sess
	return b
}

type Builder struct {
	Query   string
	Session neo4j.Session
}

func (b *Builder) NODE(tag string, object string, keyValues map[string]interface{}) *Builder {
	b.Query += "(" + tag + ":" + object

	if len(keyValues) > 0 {
		b.Query += " {"
		for i, v := range keyValues {
			b.Query += i + ":'" + v.(string) + "'"
		}
		b.Query += "}"
	}

	b.Query += ") "
	return b
}
func (b *Builder) MATCH() *Builder {
	b.Query = "MATCH "
	return b
}
func (b *Builder) CREATE() *Builder {
	b.Query = "CREATE "
	return b
}
func (b *Builder) WHERE(field, operator, value string) *Builder {
	b.Query += "WHERE " + field + " " + operator + " '" + value + "' "
	return b
}
func (b *Builder) WITH(tag string) *Builder {
	b.Query += "WITH " + tag + " "
	return b
}
func (b *Builder) ORDERBY(tag, direction string) *Builder {
	b.Query += "ORDER BY " + tag + " " + direction + " "
	return b
}
func (b *Builder) RETURN(tag string, fields []string) *Builder {
	b.Query += "RETURN "
	if len(fields) < 1 {
		b.Query += tag + " "
	} else {
		for _, v := range fields {
			b.Query += tag + "." + v + " "
		}
	}
	return b
}

func (b *Builder) COLLECT(tag string, as string) *Builder {
	b.Query += "WITH collect(" + tag + ") as " + as + " "
	return b
}

func (b *Builder) Log() *Builder {
	log.Println(b.Query)
	return b
}

type Node struct {
	ID     int64
	Labels []string
	Props  map[string]interface{}
}
type Relationship struct {
	ID      int64
	StartID int64
	EndID   int64
	Type    string
	Props   map[string]interface{}
}
type Results struct {
	Relationships map[int64]Relationship
	Nodes         map[int64]Node
	Keys          map[string]interface{}
	String        string
}

func (b *Builder) RUN() (*Results, neo4j.Result, error) {
	res, err := b.Session.Run(b.Query, nil)
	if err != nil {
		return nil, res, err
	}
	var record interface{}
	if res.Next() {
		record = res.Record().GetByIndex(0)
	}
	finalRes := new(Results)
	finalRes.Nodes = make(map[int64]Node)
	finalRes.Relationships = make(map[int64]Relationship)

	x, ok := record.([]interface{})
	if ok {
		for _, v := range x {
			n, ok := v.(neo4j.Node)
			if ok {
				finalRes.Nodes[n.Id()] = Node{
					ID:     n.Id(),
					Labels: n.Labels(),
					Props:  n.Props(),
				}
			}
			r, ok := v.(neo4j.Relationship)
			if ok {
				finalRes.Relationships[r.Id()] = Relationship{
					ID:      r.Id(),
					StartID: r.StartId(),
					EndID:   r.EndId(),
					Type:    r.Type(),
					Props:   r.Props(),
				}
			}
		}
		return finalRes, res, err
	}

	xx, ok := record.(string)
	if ok {
		finalRes.String = xx
		return finalRes, res, err
	}

	return finalRes, res, err
}
