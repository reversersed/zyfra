package storage

import (
	"context"
	"errors"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type db struct {
	sessions *mongo.Collection
	users    *mongo.Collection
}

func New(storage *mongo.Database) *db {
	db := &db{
		sessions: storage.Collection("sessions"),
		users:    storage.Collection("users"),
	}
	defer db.seed()
	return db
}
func (d *db) seed() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	documents, _ := d.users.CountDocuments(ctx, bson.D{})
	if documents > 0 {
		log.Printf("there are %d documents in database, seed canceled", documents)
		return
	}

	pass, _ := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)
	user := &UserModel{
		Id:       primitive.NewObjectID(),
		Login:    "admin",
		Password: pass,
	}

	_, err := d.users.InsertOne(context.Background(), user)
	if err != nil {
		log.Fatalf("cant seed users: %v", err)
	}
}
func (d *db) LoginUser(login, password string) (string, error) {
	reply := d.users.FindOne(context.Background(), bson.M{"login": login})
	var user UserModel
	err := reply.Decode(&user)
	if reply.Err() != nil || err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword(user.Password, []byte(password))
	if err != nil {
		return "", err
	}

	return d.createSession(), nil
}
func (d *db) createSession() string {
	session := UserSession{Id: primitive.NewObjectID(), Expiration: time.Now().UTC().Add(time.Minute)}
	_, err := d.sessions.InsertOne(context.Background(), session)
	if err != nil {
		log.Fatal(err)
	}
	return session.Id.Hex()
}
func (d *db) GetSession(session string) (*UserSession, error) {
	id, err := primitive.ObjectIDFromHex(session)
	if err != nil {
		return nil, err
	}
	reply := d.sessions.FindOne(context.Background(), bson.M{"_id": id})
	var user UserSession
	err = reply.Decode(&user)
	if reply.Err() != nil || err != nil {
		return nil, err
	}
	return &user, nil
}
func (d *db) DeleteSession(session string) error {
	id, err := primitive.ObjectIDFromHex(session)
	if err != nil {
		return err
	}
	result, err := d.sessions.DeleteOne(context.Background(), bson.M{"_id": id})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("session not found")
	}
	return nil
}
