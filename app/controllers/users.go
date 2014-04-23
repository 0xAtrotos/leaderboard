package controllers

import (
	"fmt"
	"github.com/revel/revel"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"os"
)

func (c App) CreateUser(dbname, collection, user, email, pass string) revel.Result {

	// connect to DB server(s)

	if uri == "" {
		fmt.Println("no connection string provided")
		os.Exit(1)
	}
	session, err := mgo.Dial(uri)
	if err != nil {
		panic(err)
	}

	defer session.Close()

	// select DB and Collection
	d := session.DB(dbname).C(collection)

	// TODO: Use encryption through crypto package to hash passwords

	// Query to see if user already exists in collection
	var doc User
	var results []User
	err = d.Find(bson.M{"username": user}).Sort("-timestamp").All(&results)

	if err != nil {
		panic(err)
	} else {
		if len(results) == 0 {
			//do DB operations
			doc = User{Id: bson.NewObjectId(), Username: user, Email: email, Password: pass}
			err = d.Insert(doc)
			if err != nil {
				panic(err)
			}
		} else {
			return c.RenderJson("Error: User already exists")
		}
	}

	return c.RenderJson(doc)

}

func (c App) Auth(dbname, collection, user, pass string) revel.Result {

	// connect to DB server(s)

	if uri == "" {
		fmt.Println("no connection string provided")
		os.Exit(1)
	}
	session, err := mgo.Dial(uri)
	if err != nil {
		panic(err)
	}

	defer session.Close()

	// select DB and Collection
	d := session.DB(dbname).C(collection)

	// TODO: Use encryption through crypto package to hash passwords

	// Query to authenticate

	var results []User

	err = d.Find(bson.M{"username": user, "password": pass}).Sort("-timestamp").All(&results)

	if err != nil {
		panic(err)
	} else {
		//fmt.Println(err)
		c.Session["user"] = user
		c.Flash.Success("Welcome, " + user)
	}

	return c.RenderJson(results)
}

func (c App) Logout() revel.Result {

	c.Session["user"] = ""

	session := c.Session["user"]

	return c.RenderJson(session)
}

func (c App) SaveUserStat(statName, username string) {

}

func (c App) GetUserStats(username string) {

}