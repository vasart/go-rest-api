package db

import (
	"fmt"

	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2"

	"github.com/vasart/go-rest-api/model"
)

type userCollection struct {
	ID			bson.ObjectId `bson:"_id,omitempty"`
	Username	string
	Password	string
}

func userCollectionIndex() mgo.Index {
	return mgo.Index{
		Key:		[]string{"username"},
		Unique:		true,
		DropDups:	true,
		Background:	true,
		Sparse:		true,
	}
}

func fromUser(u *model.User) *userCollection {
	return &userCollection{
		Username: u.Username,
		Password: u.Password,
	}
}

func(u *userCollection) toUser() *model.User {
	return &model.User{
		ID:			u.ID.Hex(),
		Username:	u.Username,
		Password:	u.Password,
	}
}

type UserMgoRepository struct {
	collection *mgo.Collection
}

func NewUserService(session *Session, dbName string, collectionName string) *UserMgoRepository {
 	collection := session.GetCollection(dbName, collectionName)
 	collection.EnsureIndex(userCollectionIndex())
 	return &UserMgoRepository{collection}
}

func(repo *UserMgoRepository) CreateUser(u *model.User) error {
	userModel := fromUser(u)
	return repo.collection.Insert(&userModel)
}

func(repo *UserMgoRepository) GetByUsername(username string) (*model.User, error) {
	collection := userCollection{}
	err := repo.collection.Find(bson.M{"username": username}).One(&collection)
	return collection.toUser(), err
}

func(repo *UserMgoRepository) Login(c model.Credentials) (u *model.User, err error) {
	collection := userCollection{}
	err = repo.collection.Find(bson.M{"username": c.Username}).One(&collection)
	if err != nil {
		return nil, err
	}

	u = collection.toUser()

	if u.CheckPassword(c.Password) {
		return u, nil
	}

	return nil, fmt.Errorf("password incorrect")
}
