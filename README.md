# ESP32-RTOS-Cloud Compute demo
The original project was to collect the 3D printer's XYZ positional data and the actual positional data measured by sensors, combine them, and send them to the cloud server. 
Later, the machine learning team could obtain the data from the cloud server to do analysis

### **Attention: You can modify it to make your application. Please don't create the exact application of what I did! There are more real-life challenges that you can address.**

## Project Topology
![图片](https://github.com/blaticslm/ESP32-RTOS-Cloud-Compute-demo/assets/47236078/66cdb9b9-1d4d-48e6-a910-a14d7c92b2ad)


## Why is this project useful?
This is a small-scale edge-to-cloud application. Tesla collects road data and sends it to the local data center(cloud), and the machine-learning team can obtain it from the cloud to analyze.

## Expect to learn
1. Connect ESP32 to enterprise-level WiFi (WPA2, like Eduroam)
2. Some important aspects of Real-time operating system (RTOS)
3. Address the "sending real-time data using slow transmission protocol" issue by using RTOS
5. Know the trade-off between memory space and program operational speed
6. Create a backend server and establish the connection
7. (Optional) Architect the database structure


## Project materials
1. ESP32 with external RAMs
   - https://www.amazon.com/dp/B087TNPQCV?ref=ppx_yo2ov_dt_b_product_details&th=1

2. Visual Studio Code
   -  https://code.visualstudio.com/
     
3. Platform IO
   - This is the Visual Studio Add-on for Arduino and other types of MCUs
     
## MCU program
### Topology
![micro_controller_upload_process](https://github.com/blaticslm/ESP32-RTOS-Cloud-Compute-demo/assets/47236078/dc341cb5-57ac-4cbb-a757-e0cdb5c4eaa3)

## Server
![图片](https://github.com/blaticslm/ESP32-RTOS-Cloud-Compute-demo/assets/47236078/a314aa4a-14f4-4679-aeae-0914fe764811)

## Database Architecture
![229964391-f5841ff9-3207-4ff4-a074-1adcfb7bb369](https://github.com/blaticslm/ESP32-RTOS-Cloud-Compute-demo/assets/47236078/e59a49c1-4284-4ca3-b42c-cbddfb07be1a)


