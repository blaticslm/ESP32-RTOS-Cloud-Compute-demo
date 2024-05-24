package main

//TODO: import sufficient packages
import (
	"log"
	"math"
	"strconv"

	// "net/http"
	"context"
	"fmt"

	// "github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/ioutil"
)

const (
	caFilePath = "rds-combined-ca-bundle.pem"

	// Option 1: Deploy on your local machine for testing
	MongoDB_URL = "mongodb://localhost:27017"

	// Option 2: Deploy on the AWS cloud machine with AWS DocumentDB
	/*
		MongoDB_URL = "mongodb://ciscoumich:3dprinting@docdbciscoproject.c5e2vtkjgjac.us-east-2.docdb.amazonaws.com:27017/?ssl=true&ssl_ca_certs=rds-combined-ca-bundle.pem&retryWrites=false"
	*/

)

type Machine struct {
	_id         int     `bson:"_id"`
	Job_order   int     `bson:"Job_order"`
	Machine_id  int     `bson:"Machine_ID"`
	Job_id      int     `bson:"Job_ID"`
	Layer       int     `bson:"Layer"`
	X_acc       float64 `bson:"X_acc"`
	X_input_pos float64 `bson:"X_input_pos"`
	X_act_pos   float64 `bson:"X_act_pos"`
	Y_acc       float64 `bson:"Y_acc"`
	Y_input_pos float64 `bson:"Y_input_pos"`
	Y_act_pos   float64 `bson:"Y_act_pos"`
	TimeDiff    int     `bson: "TimeDiff"`
	isPrint     bool    `bson: "IsPrint"`
	Observer    bool    `bson:"Observer"`
}

type Machine_group struct {
	Group_data []interface{} `bson:"group_data"`
}

var client *mongo.Client = nil
var connect_err error = nil

func getCustomTLSConfig(caFile string) (*tls.Config, error) {
	tlsConfig := new(tls.Config)
	certs, err := ioutil.ReadFile(caFile)

	if err != nil {
		return tlsConfig, err
	}

	tlsConfig.RootCAs = x509.NewCertPool()
	ok := tlsConfig.RootCAs.AppendCertsFromPEM(certs)

	if !ok {
		return tlsConfig, errors.New("failed parsing pem file")
	}

	return tlsConfig, nil
}

func connectToDB() (*mongo.Client, error) {
	// Option 1: Connect to the local Mongo DB
	clientOptions := options.Client().ApplyURI(MongoDB_URL)

	// Option 2: Connect to the AWS documentDB
	/*
		tlsConfig, err := getCustomTLSConfig(caFilePath)
		if err != nil {
			log.Panic("Failed getting TLS configuration: ", err)
			return nil, err
		}

		clientOptions := options.Client().ApplyURI(MongoDB_URL).SetTLSConfig(tlsConfig)

	*/

	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Panic(err) //I dont want to have os.exit(1)
		CloseClientDB(client)
		return nil, err
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Panic(err)
		CloseClientDB(client)
		return nil, err
	}

	return client, nil
}

//closing the connection to DB
func CloseClientDB(client *mongo.Client) {

	err := client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	// TODO optional you can log your closed MongoDB client
	fmt.Println("Connection to MongoDB closed.")
}

// Rest of the methods below are the key functions!

/*
	This function is connecting MongoDB and store the connection state to a local variables.
	Once we successfully connect, the we can reuse the state varible to do DB operation
	without creating the new connection session.

	@brief
		Connect to MongoDB

	@param
		none

	@return
		nil: Connect successfully
		connect_err: errors that cause connect unsuccessfully
*/
func testConnectToMongoDB() error {
	if client == nil {
		client, connect_err = connectToDB()

		if connect_err != nil {
			log.Panic(connect_err)
			connect_err = nil
			return connect_err
		}
	}

	return nil
}

/*
	This function is to return a number in the Machine<X> database such that the collection
	associate with this number in this database is empty. Each time we will upload the data
	to the empty collection to avoid messing up the existing collection.

	@brief
		Obtain the empty collection from specific database

	@param
		machine_id: the specific database id for Machine<machine_id>

	@return
		result, nil: 		obtain the empty collection number succesfully
		math.MinInt32, err: errors that fail to obtain the empty collection number
*/
func getEmptyCollectionFromMongoDB(machine_id int) (int, error) {
	err := testConnectToMongoDB()

	if err != nil {
		return math.MinInt32, err
	}

	machine_name := fmt.Sprint("machine", machine_id)

	machineDatabase := client.Database(machine_name)
	machineCollection, err := machineDatabase.ListCollectionNames(context.TODO(), bson.D{})

	if err != nil {
		return math.MinInt32, err
	}

	result := 0

	for i := 0; i < len(machineCollection); i++ {
		cur_job := machineCollection[i]
		convert, err := strconv.Atoi(cur_job[3:])

		if err != nil {
			return math.MinInt32, err
		}

		temp := math.Max(float64(result), float64(convert))
		result = int(temp)

	}

	result = result + 1
	fmt.Println(machine_name)
	fmt.Println(result)

	return result, nil

}

/*
	This function is to handle the group data from HTTP post requests.
	THe machine_id and job_id are determining which database and collection
	the data should go.

	@brief
		Save machine data array at once

	@param
		input_data: 	the array that contains many machine struct data
		machine_id: 	the specific database id for Machine<machine_id>
		job_id: 		the specific collection in Machine<machine_id>

	@return
		nil: 	Upload to DB successfully
		err: 	errors that fail to upload to DB
*/
func SaveManyToMongoDB(input_data []interface{}, machine_id int, job_id int) error {
	err := testConnectToMongoDB()
	if err != nil {
		return err
	}

	machine_name := fmt.Sprint("machine", machine_id)
	job_name := fmt.Sprint("job", job_id)

	machineDatabase := client.Database(machine_name)
	machineCollection := machineDatabase.Collection(job_name)

	_, err = machineCollection.InsertMany(context.TODO(), input_data)

	if err != nil {
		fmt.Println("insert fail")
		log.Panic(err)
		return err
	} else {
		fmt.Println("insert OK!")
	}

	return nil

}

/*
	This function is to search the data from MongoDB. The _id in the Machine struct
	is the key to immediately locate the data, and get the sequence with len == end after this data.
	The sequence is [start + 1, end].

	@brief
		Get data from the db

	@param
		machine_id: 	the specific database id for Machine<machine_id>
		job_id: 		the specific collection in Machine<machine_id>
		start:			the head of the query is start + 1. start is excluded.
		end:			the tail of the query. tail is inclusive.

	@return
		result, nil: 	Upload to DB successfully
		nil, err: 	errors that fail to upload to DB

*/
func getDataFromMongoDB(machine_id int, job_id int, start int, end int) ([]Machine, error) {
	err := testConnectToMongoDB()
	if err != nil {
		return nil, err
	}

	machine_name := fmt.Sprint("machine", machine_id)
	job_name := fmt.Sprint("job", job_id)

	machineDatabase := client.Database(machine_name)
	machineCollection := machineDatabase.Collection(job_name)

	// filter := bson.M{
	// 	"Job_ID":    job_id,
	// 	"Observer":  observer,
	// 	"Job_order": bson.M{"$gte": start, "$lte": end},
	// }

	filter := bson.M{
		"_id": bson.M{"$gt": start},
	}

	option := options.Find().SetLimit(int64(end))

	cursor, err := machineCollection.Find(context.TODO(), filter, option)

	if err != nil {
		return nil, err
	}

	defer cursor.Close(context.TODO())

	var results []Machine
	if err = cursor.All(context.TODO(), &results); err != nil {
		log.Panic(err)
		return nil, err
	}

	return results, err

}
