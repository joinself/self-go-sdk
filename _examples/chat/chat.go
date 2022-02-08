// Copyright 2020 Self Group Ltd. All Rights Reserved.

package main

import (
	"io/ioutil"
	"log"
	"os"
	"time"

	selfsdk "github.com/joinself/self-go-sdk"
	"github.com/joinself/self-go-sdk/chat"
)

// expects 1 argument - the Self ID you want to authenticate
func main() {
	groups := make(map[string]*chat.Group, 0)

	cfg := selfsdk.Config{
		SelfAppID:           os.Getenv("SELF_APP_ID"),
		SelfAppDeviceSecret: os.Getenv("SELF_APP_DEVICE_SECRET"),
		StorageKey:          "my-secret-crypto-storage-key",
		StorageDir:          "../.storage/",
	}

	if os.Getenv("SELF_ENV") != "" {
		cfg.Environment = os.Getenv("SELF_ENV")
	}

	client, err := selfsdk.New(cfg)
	if err != nil {
		panic(err)
	}

	if len(os.Args) < 2 {
		panic("you must specify a self id as an argument")
	}

	log.Println("sending chat message")

	cs := client.ChatService()

	cs.OnMessage(func(cm *chat.Message) {
		opts := map[string]interface{}{}
		cm.Message("howdy", opts)
		println("chat.message received with " + cm.Body)
		if len(cm.Objects) > 0 {
			for _, o := range cm.Objects {
				c, err := o.GetContent()
				if err != nil {
					println(err.Error())
					continue
				}
				err = ioutil.WriteFile("/tmp/obj", c, 0644)
				if err != nil {
					println(err.Error())
					continue
				}
				println(" - file on /tmp/obj")
			}
		}
		time.Sleep(5 * time.Second)
		println("sending a direct response")
		cm.Respond("tupu")
		time.Sleep(5 * time.Second)
		println("sending a new message to that conversation")
		nm := cm.Message("supu", map[string]interface{}{})
		time.Sleep(5 * time.Second)
		println("editing new message")
		nm.Edit("about to be removed")
		time.Sleep(5 * time.Second)
		println("remove new message")
		nm.Delete()
	})

	cs.OnInvite(func(g *chat.Group) {
		println("you've been invited to " + g.Name)
		g.Join()
		groups[g.GID] = g
		time.Sleep(5 * time.Second)
		g.Message("hey!", map[string]interface{}{})
	})

	cs.OnJoin(func(iss, gid string) {
		if _, ok := groups[gid]; ok {
			groups[gid].Members = append(groups[gid].Members, iss)
		}
	})

	cs.OnLeave(func(iss, gid string) {
		delete(groups, gid)
	})

	var opts map[string]interface{}
	// Public object
	obj = map[string]interface{}{
		"name": "Hello",
		"link": "https://user-images.githubusercontent.com/14011726/94132137-7d4fc100-fe7c-11ea-8512-69f90cb65e48.gif",
		"mime": "image/gif",
	}
	/*
		// Add a private object
		dat, err := os.ReadFile("/tmp/obj.png")
		if err == nil {
			println("attaching local object")
			// Private object
			obj := map[string]interface{}{
				"name": "Test",
				"data": dat,
				"mime": "image/png",
			}
			opts = map[string]interface{}{
				"objects": []map[string]interface{}{obj},
			}
		}
	*/

	cs.Message([]string{os.Args[1]}, "oyoyo!", opts)

	if err != nil {
		log.Fatal("error sending message: ", err)
	}

	time.Sleep(10 * time.Minute)

	err = client.Close()
	if err != nil {
		log.Fatal(err)
	}
}
