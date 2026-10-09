/*
 * Copyright (C) 2024 STD DEEPANSHU <https://deepanshu.in>
 * STD BOTS - Telegram: @STD_DEEPANSHU, @STDBOTS
 *
 * This file is part of FeedVackBot.
 *
 * FeedVackBot is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, version 3 of the License.
 */

package database

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	client  *mongo.Client
	db      *mongo.Database
	users   *mongo.Collection
	clones  *mongo.Collection
	tickets *mongo.Collection
}

func NewMongoDB(uri, dbName string) (*MongoDB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	db := client.Database(dbName)
	return &MongoDB{
		client:  client,
		db:      db,
		users:   db.Collection("users"),
		clones:  db.Collection("clones"),
		tickets: db.Collection("tickets"),
	}, nil
}

func (db *MongoDB) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return db.client.Disconnect(ctx)
}

func (db *MongoDB) AddUser(id int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	opts := options.Update().SetUpsert(true)
	filter := bson.M{"_id": id}
	update := bson.M{"$setOnInsert": bson.M{
		"_id":       id,
		"joined_at": time.Now(),
		"notif":     true,
		"ban_status": BanStatus{IsBanned: false},
	}}

	_, err := db.users.UpdateOne(ctx, filter, update, opts)
	return err
}

func (db *MongoDB) IsUserExist(id int64) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	count, err := db.users.CountDocuments(ctx, bson.M{"_id": id})
	return count > 0, err
}

func (db *MongoDB) GetUser(id int64) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user User
	err := db.users.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (db *MongoDB) TotalUsersCount() (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return db.users.CountDocuments(ctx, bson.M{})
}

func (db *MongoDB) GetAllUsers() ([]User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := db.users.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []User
	if err = cursor.All(ctx, &users); err != nil {
		return nil, err
	}
	return users, nil
}

func (db *MongoDB) DeleteUser(id int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := db.users.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (db *MongoDB) BanUser(id int64, duration int, reason string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"ban_status.is_banned": true,
			"ban_status.duration":  duration,
			"ban_status.reason":    reason,
			"ban_status.banned_at": time.Now(),
		},
	}
	_, err := db.users.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

func (db *MongoDB) UnbanUser(id int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"ban_status.is_banned": false,
		},
	}
	_, err := db.users.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

func (db *MongoDB) GetBanStatus(id int64) (*BanStatus, error) {
	user, err := db.GetUser(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}
	return &user.BanStatus, nil
}

func (db *MongoDB) GetNotif(id int64) (bool, error) {
	user, err := db.GetUser(id)
	if err != nil {
		return false, err
	}
	if user == nil {
		return false, fmt.Errorf("user not found")
	}
	return user.Notif, nil
}

func (db *MongoDB) SetNotif(id int64, notif bool) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := db.users.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"notif": notif}})
	return err
}

func (db *MongoDB) AddClone(clone *Clone) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	opts := options.Update().SetUpsert(true)
	_, err := db.clones.UpdateOne(ctx, bson.M{"bot_token": clone.BotToken}, bson.M{"$set": clone}, opts)
	return err
}

func (db *MongoDB) GetClone(botToken string) (*Clone, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var clone Clone
	err := db.clones.FindOne(ctx, bson.M{"bot_token": botToken}).Decode(&clone)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &clone, nil
}

func (db *MongoDB) GetAllClones() ([]Clone, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := db.clones.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var clones []Clone
	if err = cursor.All(ctx, &clones); err != nil {
		return nil, err
	}
	return clones, nil
}

func (db *MongoDB) DeleteClone(botToken string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := db.clones.DeleteOne(ctx, bson.M{"bot_token": botToken})
	return err
}

func (db *MongoDB) CreateTicket(userID int64) (*Ticket, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	count, err := db.tickets.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	ticketID := fmt.Sprintf("#%04d", count+1)
	ticket := &Ticket{
		TicketID:  ticketID,
		UserID:    userID,
		Status:    "open",
		CreatedAt: time.Now(),
		Messages:  []TicketMessage{},
	}

	_, err = db.tickets.InsertOne(ctx, ticket)
	if err != nil {
		return nil, err
	}
	return ticket, nil
}

func (db *MongoDB) GetTicket(ticketID string) (*Ticket, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var ticket Ticket
	err := db.tickets.FindOne(ctx, bson.M{"ticket_id": ticketID}).Decode(&ticket)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &ticket, nil
}

func (db *MongoDB) GetUserTickets(userID int64) ([]Ticket, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := db.tickets.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tickets []Ticket
	if err = cursor.All(ctx, &tickets); err != nil {
		return nil, err
	}
	return tickets, nil
}

func (db *MongoDB) AddTicketMessage(ticketID string, msg TicketMessage) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{"$push": bson.M{"messages": msg}}
	_, err := db.tickets.UpdateOne(ctx, bson.M{"ticket_id": ticketID}, update)
	return err
}

func (db *MongoDB) CloseTicket(ticketID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{"$set": bson.M{"status": "closed", "closed_at": time.Now()}}
	_, err := db.tickets.UpdateOne(ctx, bson.M{"ticket_id": ticketID}, update)
	return err
}

func (db *MongoDB) RateTicket(ticketID string, rating int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{"$set": bson.M{"rating": rating}}
	_, err := db.tickets.UpdateOne(ctx, bson.M{"ticket_id": ticketID}, update)
	return err
}

func (db *MongoDB) GetStats() (*Stats, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userCount, err := db.users.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	cloneCount, err := db.clones.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	ticketCount, err := db.tickets.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	return &Stats{
		TotalUsers:   userCount,
		TotalClones:  cloneCount,
		TotalTickets: ticketCount,
	}, nil
}
