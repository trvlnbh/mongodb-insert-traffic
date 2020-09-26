package main

import (
	"log"
	"strings"
	"time"
	"fmt"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
	"github.com/icrowley/fake"
)

type CoordData struct {
	Timestamp	int64	`bson:"timestamp" json:"timestamp"`
	Date		string	`bson:"date" json:"date"`
	FirstName	string	`bson:"first_name" json:"first_name"`
	LastName	string	`bson:"last_name" json:"last_name"`
	Coordinate	Coord	`bson:"coordinate" json:"coordinate"`
}

type Coord struct {
	Latitude	float32	`bson:"latitude" json:"latitude"`
	Longitude	float32 `bson:"longitude" json:"longitude"`
}

type FavorColorData struct {
	Timestamp	int64		`bson:"timestamp" json:"timestamp"`
	Date		string		`bson:"date" json:"date"`
	PsInfo		PersonInfo	`bson:"personal_info" json:"personal_info"`
	Color		string		`bson:"color" json:"color"`
}

type PersonInfo struct {
	City		string	`bson:"city" json:"city"`
	LastName	string	`bson:"last_name" json:"last_name"`
	Gender		string	`bson:"gender" json:"gender"`
}

type EmailData struct {
	Timestamp	int64	`bson:"timestamp" json:"timestamp"`
	Date		string	`bson:"date" json:"date"`
	UserName	string	`bson:"user_name" json:"user_name"`
	Email		string	`bson:"email" json:"email"`
	Text		string	`bson:"text" json:"text"`
}

func (m *MTF) InsertData(tp int64) {
	dataCount := 0
	mutex := &sync.Mutex{}

	go func() {
		tk := time.NewTicker(1 * time.Minute)
		defer tk.Stop()

		for {
			select {
			case <-tk.C:
				fmt.Printf("insert %d per minute\n", dataCount)
				mutex.Lock()
				dataCount = 0
				mutex.Unlock()
			}
		}
	}()

	switch tp {
	case 1:
		for {
			d := generateRandomCoordData()
	
			_, err := m.cl.Database(DB).Collection(CollCoord).InsertOne(m.ctx, d)
			if err != nil {
				log.Fatalln(err)
			}
			mutex.Lock()
			dataCount++
			mutex.Unlock()

			time.Sleep(1 * time.Millisecond)
		}
	case 2:
		for {
			d := generateRandomColorData()

			_, err := m.cl.Database(DB).Collection(CollColor).InsertOne(m.ctx, d)
			if err != nil {
				log.Fatalln(err)
			}
			mutex.Lock()
			dataCount++
			mutex.Unlock()

			time.Sleep(2 * time.Millisecond)
		}
	case 3:
		for {
			d := generateRandomEmailData()

			_, err := m.cl.Database(DB).Collection(CollEmail).InsertOne(m.ctx, d)
			if err != nil {
				log.Fatalln(err)
			}
			mutex.Lock()
			dataCount++
			mutex.Unlock()

			time.Sleep(2 * time.Millisecond)
		}
	}
}

func generateRandomCoordData() CoordData {
	data := CoordData{
		Timestamp:	time.Now().UnixNano() / int64(time.Millisecond),
		Date:		time.Now().UTC().Format(time.RFC3339Nano),
		FirstName:	fake.FirstName(),
		LastName:	fake.LastName(),
	}
	coord := Coord{
		Latitude:	fake.Latitude(),
		Longitude:	fake.Longitude(),
	}
	data.Coordinate = coord

	return data
}

func generateRandomColorData() FavorColorData {
	data := FavorColorData{
		Timestamp:	time.Now().UnixNano() / int64(time.Millisecond),
		Date:		time.Now().UTC().Format(time.RFC3339Nano),
		Color:		fake.Color(),
	}
	pinfo := PersonInfo{
		City:		fake.City(),
		LastName:	fake.LastName(),
		Gender:		fake.Gender(),
	}
	data.PsInfo = pinfo

	return data
}

func generateRandomEmailData() EmailData {
	data := EmailData{
		Timestamp:	time.Now().UnixNano() / int64(time.Millisecond),
		Date:		time.Now().UTC().Format(time.RFC3339Nano),
		UserName:	fake.UserName(),
		Email:		fake.EmailAddress(),
		Text:		fake.Sentence(),
	}

	return data
}

func (m *MTF) SetupCollection(tp int64) {
	switch tp {
	case 1:
		err := m.cl.Database(DB).RunCommand(
			m.ctx,
			bson.D{{"create", CollCoord}},
		).Err()
		if err != nil && !strings.Contains(err.Error(), "already exists") {
			log.Fatalln(err)
		}
	
		err = m.cl.Database(ADMIN).RunCommand(
			m.ctx,
			bson.D{
				{"shardCollection", DB+"."+CollCoord},
				{"key", bson.M{"_id": "hashed"}},
				{"numInitialChunks", 8},
			},
		).Err()
		if err != nil {
			log.Fatalln(err)
		}
	case 2:
		err := m.cl.Database(DB).RunCommand(
			m.ctx,
			bson.D{{"create", CollColor}},
		).Err()
		if err != nil && !strings.Contains(err.Error(), "already exists") {
			log.Fatalln(err)
		}
	
		err = m.cl.Database(ADMIN).RunCommand(
			m.ctx,
			bson.D{
				{"shardCollection", DB+"."+CollColor},
				{"key", bson.M{"_id": "hashed"}},
				{"numInitialChunks", 8},
			},
		).Err()
		if err != nil {
			log.Fatalln(err)
		}
	case 3:
		err := m.cl.Database(DB).RunCommand(
			m.ctx,
			bson.D{{"create", CollEmail}},
		).Err()
		if err != nil && !strings.Contains(err.Error(), "already exists") {
			log.Fatalln(err)
		}
	}
}
