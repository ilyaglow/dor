package dor

import (
	"log"
	"sync"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

const docsLimit = 100

// MongoStorage implements the Storage interface for MongoDB
type MongoStorage struct {
	sess *mgo.Session // mongodb session
	db   string       // database name
	c    string       // collection name
	wNum int          // number of workers
	size int          // size of batch request
	ret  bool         // data retention
}

// NewMongoStorage bootstraps MongoStorage, creates indexes
//	u is the Mongo URL
//	db is the database name
//	col is the collection name
//	size is the bulk message size
//	w is number of workers
//	ret is the data retention option
func NewMongoStorage(u string, db string, col string, size int, w int, ret bool) (*MongoStorage, error) {
	sess, err := mgo.Dial(u)
	if err != nil {
		return nil, err
	}

	index := mgo.Index{
		Key:        []string{"domain", "rank", "source"},
		Background: true,
		Sparse:     true,
	}

	if ret {
		index.DropDups = true
	}

	err = sess.DB(db).C(col).EnsureIndex(index)
	if err != nil {
		return nil, err
	}

	if ret {
		expindex := mgo.Index{
			Key:         []string{"last_update"},
			Background:  true,
			ExpireAfter: time.Hour * 24 * 7,
		}

		err = sess.DB(db).C(col).EnsureIndex(expindex)
		if err != nil {
			return nil, err
		}
	}

	return &MongoStorage{
		sess: sess,
		wNum: w,
		db:   db,
		c:    col,
		size: size,
	}, nil
}

// Put implements Storage interface method Put
//	s - is the data source
//	t - is the data datetime
func (m *MongoStorage) Put(c <-chan Rank, s string, t time.Time) error {
	var wg sync.WaitGroup
	wg.Add(m.wNum)

	for i := 1; i <= m.wNum; i++ {
		go func() {
			defer wg.Done()
			m.send(c, s, t)
		}()
	}

	wg.Wait()
	return nil
}

// func (m *MongoStorage) getSimple(d string) ([]ExtendedRank, error) {
// }
//
// func (m *MongoStorage) getMajestic(d string) ([]MajesticRank, error) {
// }

// Get implements Storage interface method Get
func (m *MongoStorage) Get(d string, sources ...string) ([]Rank, error) {
	s := m.sess.Copy()
	c := s.DB(m.db).C(m.c)

	var query *mgo.Query
	var ranks []Rank
	e := ExtendedRank{}

	if len(sources) > 0 {
		for i := range sources {
			err := c.Find(bson.M{"domain": d, "source": sources[i]}).Sort("-last_update").One(&e)
			if err != nil {
				log.Println(err)
				continue
			}
			ranks = append(ranks, e)
		}
	} else {
		query = c.Find(bson.M{"domain": d}).Sort("-last_update").Limit(docsLimit)
		items := query.Iter()
		for items.Next(&e) {
			ranks = append(ranks, e)
		}
	}

	return ranks, nil
}

// GetMore implements Storage GetMore function
func (m *MongoStorage) GetMore(d string, lps int, sources ...string) ([]Rank, error) {
	s := m.sess.Copy()
	c := s.DB(m.db).C(m.c)

	var query *mgo.Query
	var ranks []Rank
	e := ExtendedRank{}

	// check if lps is not bigger than allowed
	if lps > docsLimit {
		lps = docsLimit
	}

	if len(sources) > 0 {
		for i := range sources {
			query = c.Find(bson.M{"domain": d, "source": sources[i]}).Sort("-last_update").Limit(lps)
			items := query.Iter()
			for items.Next(&e) {
				ranks = append(ranks, e)
			}
		}
	} else {
		query = c.Find(bson.M{"domain": d}).Sort("-last_update").Limit(docsLimit)
		items := query.Iter()
		for items.Next(&e) {
			ranks = append(ranks, e)
		}
	}

	return ranks, nil
}

func (m *MongoStorage) send(c <-chan Rank, s string, t time.Time) error {
	mc := m.sess.Copy()
	col := mc.DB(m.db).C(m.c)

	bulk := col.Bulk()
	bulk.Unordered()
	i := 0

	for r := range c {
		extr := &ExtendedRank{
			Domain:     r.GetDomain(),
			Rank:       r.GetRank(),
			Source:     s,
			LastUpdate: t,
		}

		bulk.Insert(extr)
		i++

		if i == m.size {
			log.Println("bulk.Run()")
			if _, err := bulk.Run(); err != nil {
				log.Println("mongo storage: failed to run bulk run")
			}

			bulk = col.Bulk()
			bulk.Unordered()
			i = 0
		}
	}

	if i != 0 {
		log.Println("bulk.Run() the tail")
		if _, err := bulk.Run(); err != nil {
			log.Println("mongo storage: failed to run bulk run")
		}
	}

	return nil
}
