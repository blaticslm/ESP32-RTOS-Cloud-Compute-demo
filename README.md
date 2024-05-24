# ESP32-RTOS-Cloud Compute demo
The original project was to collect the 3D printer's XYZ positional data and the actual positional data measured by sensors, combine them, and send them to the cloud server. 
Later, the machine learning team could obtain the data from the cloud server to do analysis

**Don't create the exact application of what I did! There are more real-life challenges that you can address.**

## Project Topology
![图片](https://github.com/blaticslm/ESP32-RTOS-Cloud-Compute-demo/assets/47236078/326d7604-183d-40f3-ac77-a1f027582a76)

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

## 
