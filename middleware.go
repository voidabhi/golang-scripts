package main

import (
	"github.com/codegangsta/negroni"
	"github.com/gorilla/context"
	"github.com/unrolled/render"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
	"net/http"
	"os"
)

type key int

const db key = 0

func GetDb(r *http.Request) *mgo.Database {
	if rv := context.Get(r, db); rv != nil {
		return rv.(*mgo.Database)
	}
	return nil
}

func SetDb(r *http.Request, val *mgo.Database) {
	context.Set(r, db, val)
}

type Site struct {
	Env  string
	Id   bson.ObjectId `bson:"_id"`
	Name string        `bson:"name"`
}

func getSite(db *mgo.Database) *Site {
	var aSite Site
	sites := db.C("COLLECTION__NAME")

	err := sites.Find(bson.M{"name": "some site name"}).One(&aSite)
	if err != nil {
		log.Print(err)
	}
	aSite.Env = os.Getenv("APP_ENV")
	return &aSite
}

func MongoMiddleware() negroni.HandlerFunc {
	database := os.Getenv("DB_NAME")
	session, err := mgo.Dial("127.0.0.1:27017")

	if err != nil {
		panic(err)
	}

	return negroni.HandlerFunc(func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		reqSession := session.Clone()
		defer reqSession.Close()
		db := reqSession.DB(database)
		SetDb(r, db)
		next(rw, r)
	})
}

func main() {
	r := render.New(render.Options{
		Layout: "index",
		Delims: render.Delims{"[[", "]]"},
	})

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		db := GetDb(req)
		site := getSite(db)
		r.HTML(w, http.StatusOK, "home", site)
	})

	n := negroni.Classic()
	n.Use(MongoMiddleware())
	n.UseHandler(mux)
	n.Run(":3200")
}
