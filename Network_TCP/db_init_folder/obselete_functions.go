package main

/*
//mongoDB.go:

//upload to MongoDB:

func saveOneToMongoDB(input_data interface{}) error {

	fmt.Println(input_data)

	clientOptions := options.Client().ApplyURI(MongoDB_URL)

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Panic(err) //I dont want to have os.exit(1)
		CloseClientDB(client)
		return err
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Panic(err)
		CloseClientDB(client)
		return err
	}
	//fmt.Println("Connect to D")

	object, ok := input_data.(Machine)

	if !ok {
		CloseClientDB(client)
		return fmt.Errorf("input data is not machine type")
	}

	machine_name := fmt.Sprint("machine", object.Machine_id)

	machineDatabase := client.Database("machine_print")
	machineCollection := machineDatabase.Collection(machine_name)

	_, err = machineCollection.InsertOne(context.TODO(), input_data)

	if err != nil {
		//fmt.Println("insert fail")
		log.Panic(err)
		CloseClientDB(client)
		return err
	} else {
		fmt.Println("insert OK!")
	}

	CloseClientDB(client)
	return nil //if everything is ok, then there is no error

}

func saveManyToMongoDB(input_data []interface{}, machine_id int) error {

	clientOptions := options.Client().ApplyURI(MongoDB_URL)

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Panic(err) //I dont want to have os.exit(1)
		CloseClientDB(client)
		return err
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Panic(err)
		CloseClientDB(client)
		return err
	}

	//fmt.Println("Connect to D")

	machine_name := fmt.Sprint("machine", machine_id)

	machineDatabase := client.Database("machine_print")
	machineCollection := machineDatabase.Collection(machine_name)

	_, err = machineCollection.InsertMany(context.TODO(), input_data)

	if err != nil {
		//fmt.Println("insert fail")
		log.Panic(err)
		CloseClientDB(client)
		return err
	} else {
		fmt.Println("insert OK!")
	}

	//temp solution: each request will be not disconnecting after
	CloseClientDB(client)

	return nil //if everything is ok, then there is no error

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

//Obselete!
//uploadhandler: for receiving the json
func uploadHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Receive a single upload request at ", time.Now().Format(time.ANSIC))

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	// w.Header().Set("Content-Type", "application/json")

	if r.Method == "OPTIONS" {
		fmt.Println("upload OPTIONS request detected!")
		return
	}

	//fmt.Println("new insert request is coming.")
	decoder := json.NewDecoder(r.Body)
	var machine Machine
	err := decoder.Decode(&machine)

	fmt.Println(machine._id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = SaveOneToMongoDB(Machine(machine))

	if err != nil {
		log.Panic(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responseString := fmt.Sprintf("Upload single successfully! Current time: %v.", time.Now().Format(time.ANSIC))
	w.Write([]byte(responseString))
}

//obseleted///////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// new upload a single to mongoDB
func SaveOneToMongoDB(input_data interface{}) error {
	err := testConnectToMongoDB()
	if err != nil {
		return err
	}

	object, ok := input_data.(Machine)

	if !ok {
		return fmt.Errorf("input data is not machine type")
	}

	machine_name := fmt.Sprint("machine", object.Machine_id)

	machineDatabase := Client.Database("machine_print")
	machineCollection := machineDatabase.Collection(machine_name)

	_, err = machineCollection.InsertOne(context.TODO(), input_data)

	if err != nil {
		//fmt.Println("insert fail")
		log.Panic(err)
		return err
	} else {
		fmt.Println("insert OK!")
	}
	return nil //if everything is ok, then there is no error

}


*/
