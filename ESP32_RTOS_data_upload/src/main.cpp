#include <Arduino.h>
#include <WiFi.h>
#include <esp_wpa2.h>
#include <esp_wifi.h>
#include <HTTPClient.h>
#include <esp_task_wdt.h>
#include <ArduinoJson.h>
#include <time.h>
#include <cmath>
#include "Certificate_Meduroam.h"

//Creator: Mingcheng Li
//Date: 5/23/2024

//FreeRTOS parameters of ESP32
#define TASK2_STACK_SIZE 1<<12
#define TASK1_STACK_SIZE 1<<13
#define DOCSIZE 384400 

//Machine ID, please see section of Database topology in Readme
#define MACHINE_ID 30

//Sensor calibration, change it accordingly to your project
//ADXL335 calibrations
#define XRawMin 649
#define YRawMin 670
#define XRawMax 3139
#define YRawMax 3160

//pins, see ESP32 pinout map
#define x_pin 33
#define y_pin 32
#define SWITCH_PIN 4
#define LED_PIN 2

// My application uses second serial
// I did not find GPIO25 and 26 to be a UART port but it works
// GPIO16 and 17 are for external rams. 
#define RXD2 26
#define TXD2 25

//SpiRam method, I am using exteral RAM of ESP32 Wrover-E chip
// https://arduinojson.org/v6/how-to/use-external-ram-on-esp32/
struct SpiRamAllocator {
  void* allocate(size_t size) {
    return heap_caps_malloc(size, MALLOC_CAP_SPIRAM);
  }

  void deallocate(void* pointer) {
    heap_caps_free(pointer);
  }

  void* reallocate(void* ptr, size_t new_size) {
    return heap_caps_realloc(ptr, new_size, MALLOC_CAP_SPIRAM);
  }
};
using SpiRamJsonDocument = BasicJsonDocument<SpiRamAllocator>;

//Base URL
//Change accordingly to your server
const char* PRE_ADDR = "http://3.21.128.133:8080";

//Timezone: America/Detroit
//When connecting to WAP2 enterprise WiFI, we need to setup the right time for machine 
//The code below is the EST
const char *ntpServer = "pool.ntp.org";
const char *timezoneEST = "EST5EDT,M3.2.0/2,M11.1.0";


//buffer and task switching variables, please see MCU program section
bool task1docInit = true;
bool task2docInit = false;
bool processing = false; 
bool sending = false;
bool ready = true; //indicator for sending
int buffer_index = 0;
int job_id = 1; 

//Unused Z variable
int layer = 1;

//for determining the collection and sending 
//Please see MCU program section
int counter = 0; 

//the _id field, start at one. This is overriding the database ID.
// Overriding the _ID field is easier for us to get the data (if the database is MongoDB)
unsigned long job_order = 1; 


//functional variables
TaskHandle_t Task1;
TaskHandle_t Task2; 

// SpiRamJson can't be initialized in the void setup()
// I declared two global SpiRamJson pointers so that two tasks can directly use them
SpiRamJsonDocument *holder1;
SpiRamJsonDocument *holder2;

HTTPClient http; 


void showHeapRAM();
bool isEqual(float x, float y);
void Task1code( void * pvParameters );
void Task2code( void * pvParameters );

void setup() {
  Serial.begin(115200);
  Serial2.begin(115200, SERIAL_8N1, RXD2, TXD2); //RXD2 and TXD2 must be other than 16 and 17 since PSRAM is using these pins

  Serial.println("ESP32 MAC address: ");
  Serial.println(WiFi.macAddress());

  //connection indicator
  pinMode(LED_PIN, OUTPUT);
  pinMode(x_pin, INPUT);
  pinMode(y_pin, INPUT);

  WiFi.disconnect(true); //disconnect form wifi to set new wifi connection
  WiFi.mode(WIFI_STA);   //init wifi mode

  //School wifi
  //You need to donload your school's ca_certificate
  Serial.printf("Connecting to WiFi: %s ", SSID1);
  esp_wifi_sta_wpa2_ent_set_ca_cert((uint8_t *)incommon_ca, strlen(incommon_ca)+1);
  esp_wifi_sta_wpa2_ent_set_identity((uint8_t *)IDENTITY, strlen(IDENTITY));
  esp_wifi_sta_wpa2_ent_set_username((uint8_t *)IDENTITY, strlen(IDENTITY));
  esp_wifi_sta_wpa2_ent_set_password((uint8_t *)PASSWORD, strlen(PASSWORD));
  esp_wifi_sta_wpa2_ent_enable();

  WiFi.begin(SSID1);
  
  while (WiFi.status() != WL_CONNECTED)
  {
    delay(500);
    Serial.print(F("."));
    counter++;
    if (counter >= 100)
    { //after 60 seconds timeout - reset board
      Serial.println("connect fail");
      //counter = 0;
      ESP.restart();
    }
  }
  Serial.println(F(" connected!"));
  Serial.print(F("IP address set: "));
  Serial.println(WiFi.localIP()); //print LAN IP

  //set Eastern standard time
  configTime(0, 0, ntpServer);
  setenv("TZ", timezoneEST, 1);

  //REQUIRE 3 steps in INITIALIZATION!:
  //Step 1, test the connection of the servers and initialize the database connection
  String testURL = PRE_ADDR + (String)"/test";
  http.begin(testURL);

  if(http.GET() <= 0){
    Serial.println("Address fail! Please change the address.");
    return; 
  }

  //Step 2, get new empty job collection index from given machine id
  String getCollectionID_URL = PRE_ADDR + (String)"/getNewJobId/" + String(MACHINE_ID);

  Serial.println(getCollectionID_URL);

  http.begin(getCollectionID_URL);

  if(http.GET() <= 0){
    Serial.println("Get collection ID fail! Please change the address.");
    return;
  } 
  job_id = http.getString().toInt();

  http.end();


  //Step 3, get ready for the actual transmission
  String groupUpload_URL = PRE_ADDR + (String)"/groupUpload/" + String(MACHINE_ID) + '/' + job_id;
  Serial.println(groupUpload_URL);
  
  http.begin(groupUpload_URL);  
  http.addHeader("Content-Type", "application/json");  
  http.setReuse(true);

  //Indicator: the next work is ready to transmit!
  digitalWrite(LED_PIN, HIGH);


  //RTOS programs
  //Sending process
  xTaskCreatePinnedToCore(
          Task2code,  /* Task function. */
          "Task2",    /* name of task. */
          TASK2_STACK_SIZE,      /* Stack size of task */
          NULL,       /* parameter of the task */
          2,          /* priority of the task */
          &Task2,     /* Task handle to keep track of created task */
          1);         /* pin task to core 1 */    

  //Collecting process
  //Protocol Core! Try to run light weight process
  xTaskCreatePinnedToCore(
          Task1code,  /* Task function. */
          "Task1",    /* name of task. */
          TASK1_STACK_SIZE,      /* Stack size of task */
          NULL,       /* parameter of the task */
          2,          /* priority of the task */
          &Task1,     /* Task handle to keep track of created task */
          0);         /* pin task to core 0 */    
  
}

void loop() {
  // put your main code here, to run repeatedly:
  vTaskDelete(NULL);
}

void showHeapRAM() {
  Serial.println((String)"Total heap: "+ ESP.getHeapSize());
  Serial.println((String)"Free heap: "+ ESP.getFreeHeap());
  Serial.println();
}

bool isEqual(float x, float y) {
  float tolerance = 0.01;
  return fabs(x - y) < tolerance; 
}

void Task1code( void * pvParameters ) {

  Serial.print("Collection task running on core ");
  Serial.println(xPortGetCoreID());


  //local JsonDocument, initialzing here can uses PSRAM
  SpiRamJsonDocument group_data1(DOCSIZE); 
  SpiRamJsonDocument group_data2(DOCSIZE); 
  JsonArray object;

  //using global varible to store the reference for these two JsonDocuments and process them in Task2
  holder1 = &group_data1;
  holder2 = &group_data2;

  float x_pos_prev = 0;
  float y_pos_prev = 0;
  bool printing = false;       // Indication of the data is doing print job or not
  bool update = false;         // To determine update the time difference or not
  unsigned long prev_time = 0; // for calculating the time difference between two samples 

  showHeapRAM();

  for(;;){

    if(digitalRead(SWITCH_PIN) != HIGH)  {
      vTaskDelay(1); //feeding the dog
      continue;
    }

    //buffer switching 
    if(!processing) {
      if(task1docInit && !task2docInit){
        Serial.println("Creating group data1");
        object = group_data1.createNestedArray("group_data"); 

      } else if(!task1docInit && task2docInit) {
        Serial.println("Creating group data2");
        object = group_data2.createNestedArray("group_data"); 
        
      }
      processing = true;
    }

    //waiting for the serial window has new information
    if(!update) {
      update = true;
      prev_time = millis();
    }
    

    if(!Serial2.available()) {
      vTaskDelay(1); //feeding the dog
      continue; 
    }

    //Serialized the string from the serial window to json string
    StaticJsonDocument<64> doc;
    String data = Serial2.readStringUntil('\n');
    Serial.println(data + ',' + (String)counter);
    if(deserializeJson(doc, data)) {
      Serial.println("Failed to deserialize the input string");
      doc.garbageCollect();
      update = false;
      continue;
    }

    //calculating the time difference, and get the sample acceleration
    unsigned long time_diff = millis() - prev_time;
    update = false;

    float x_pos = roundf((float)doc["X"] * 100) / 100;
    float y_pos = roundf((float)doc["Y"] * 100) / 100;

    float accelX = map(analogRead(x_pin), XRawMin, XRawMax, -3000, 3000)/1000.0;
    float accelY = map(analogRead(y_pin), YRawMin, YRawMax, -3000, 3000)/1000.0;

    if(isEqual(x_pos, x_pos_prev) && isEqual(y_pos, y_pos_prev)) {
      printing = false; 
    } else {
      printing = true;
    }

    x_pos_prev = x_pos;
    y_pos_prev = y_pos;

    //to avoid buffer overflow, once it reach the limit, we need to stop to wait
    if(counter <= 2000){ 
      JsonObject object_in = object.createNestedObject();

      object_in["_id"] = job_order;
      object_in["Job_order"] = job_order;
      object_in["Machine_ID"] = MACHINE_ID;
      object_in["Job_ID"] = job_id;
      object_in["Layer"] = layer;
      object_in["X_acc"] = accelX;
      object_in["X_input_pos"] = x_pos;
      object_in["X_act_pos"] = 0;
      object_in["Y_acc"] = accelY;
      object_in["Y_input_pos"] = y_pos;
      object_in["Y_act_pos"] = 0;
      object_in["TimeDiff"] = time_diff;
      object_in["IsPrint"] = printing;
      object_in["Observer"] = false; 

      if(job_order % 10 == 0){
        layer++;
      } 

      job_order++;

      counter++;
    }

    //Once the counter reach the lower limit and sending process is idling: sneding will happen
    if(counter >= 500 && ready) {
      String current_in_Task1 = (String)"Collecting done! counter: " + counter;
      Serial.println(current_in_Task1);

      processing = false; 
      counter = 0;
      ready = false;


      //buffer sending 
      if(task1docInit && !task2docInit){
        sending = true;
        buffer_index = 1;

        task1docInit = false;
        task2docInit = true;

      } else if(!task1docInit && task2docInit){
        sending = true;
        buffer_index = 2;

        task1docInit = true;
        task2docInit = false;
      } 
      //rest of them give to task2 to deal with
    }

  }


}

void Task2code( void * pvParameters ) {
  Serial.print("Sending task running on core ");
  Serial.println(xPortGetCoreID());

  unsigned long duration; //for testing internet sending duration
  unsigned long cur_max = 0; //record largest response time
  uint8_t* toSend = (uint8_t*)heap_caps_malloc(DOCSIZE * sizeof(uint8_t), MALLOC_CAP_SPIRAM);

  for(;;){
    if(digitalRead(SWITCH_PIN) != HIGH){
      vTaskDelay(1);
      continue;
    }

    if(!processing && sending) {
    //if processing is false, and sending is true, then we need to send buffer to cloud
      duration = millis();

      switch(buffer_index){
        case 1:

          Serial.println("Sending 1...");
          serializeJson((*holder1), toSend, DOCSIZE * sizeof(char));
          break;

        case 2:

          Serial.println("Sending 2...");
          serializeJson((*holder2), toSend, DOCSIZE * sizeof(char));
          break;

      }

      if(http.POST(toSend, DOCSIZE * sizeof(char))>0) {
        Serial.println(http.getString());
      }


      Serial.println("Cleaning");
      switch(buffer_index){
        case 1:
          (*holder1).clear();
          break;

        case 2:
          (*holder2).clear();
          break;  
      }

      //calculate the time difference
      duration = millis() - duration; 
      
      if(cur_max < duration){
        cur_max = duration;
      }
      Serial.print((String)"Duration: " + duration);
      Serial.print("\t");
      Serial.println((String)"Cur_max: " + cur_max);


      buffer_index = 0;
      sending = false;
      ready = true;

    } 
  } 
}