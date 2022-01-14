package mgoc

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Connect ...
func Connect() error {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017/yumi"))
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return err
	}

	type UserInfo struct {
		ID     primitive.ObjectID `bson:"_id"`
		Name   string
		Age    uint16
		Weight uint32
	}
	query := make(primitive.M)
	query["name"] = "xxx"
	cursor, err := client.Database("yumi").Collection("admin").Find(ctx, query)
	if err != nil {
		fmt.Println(err.Error())
	}

	uis := []UserInfo{}
	fmt.Println(cursor.All(ctx, &uis))
	fmt.Println(uis)

	ret := client.Database("yumi").Collection("admin").FindOne(ctx, primitive.M{"_id": uis[0].ID})
	if err != nil {
		fmt.Println(err.Error(), 1)
	}

	uis2 := make(map[string]interface{})
	fmt.Println(ret.Decode(&uis2), 2)
	fmt.Println(uis2)

	// fmt.Println(client.Database("yumi").CreateCollection(ctx, "media"))
	return nil
}

func TestConnect(t *testing.T) {
	err := Connect()
	if err != nil {
		t.Log(err)
	}
}

type ClientInfoList struct {
	ID                 string `json:"id" bson:"_id"`
	TenancyID          string `json:"tenancy_id" bson:"tenancy_id"`
	GroupID            string `json:"group" bson:"group"`
	Mac                string `json:"mac" bson:"mac"`
	HostName           string `json:"host_name" bson:"hostname"`
	ClientInfoListOS   `bson:"osinfo"`
	OSType             string `json:"os_type" bson:"ostype"`
	Status             int    `json:"status" bson:"status"`
	ClientInfoListUser `bson:"user"`
	GroupPath          string `json:"group_path" bson:"-"`
}

type ClientInfoListOS struct {
	OSName    string `json:"os_name" bson:"name"`
	OSVersion string `json:"os_version" bson:"ver"`
	OSArch    string `json:"os_arch" bson:"arch"`
}

type ClientInfoListUser struct {
	PersonName string `json:"person_name" bson:"person"`
}

func TestSetMax(t *testing.T) {
	// cli, err := New("mongodb://localhost:27017/client")
	cli, err := New("mongodb://root:Admin_123@10.34.4.89:27017")
	if err != nil {
		t.Error(err)
		return
	}

	cil := ClientInfoList{}
	coll := cli.Database("client").Collection("clientinfo")
	if err := coll.FindOne(context.Background(), primitive.M{"_id": "E176657B-F2FB-4EFA-34EA-3D129E198E74"}).Decode(&cil); err != nil {
		t.Error(err)
		return
	} else {
		t.Log(cil)
	}

	bs, err := json.Marshal(&cil)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(string(bs))
	return
}
