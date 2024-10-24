#include <Wire.h>
#include <dht11.h>

#define DHT11PIN 7


dht11 DHT11;

void setup() {
  Serial.begin(9600);
  Wire.begin(0x55);  // Join I²C bus as slave with address 0x08
  Wire.onRequest(requestEvent);  // Register function to respond to requests
}

int humidity;
int temperature;

void loop() {
  delay(1000);
  Serial.println();

  int chk = DHT11.read(DHT11PIN);
  Serial.print("Humidity (%): ");
  Serial.println((float)DHT11.humidity, 2); //for checking in serial monitor
  Serial.print("Temperature  (C): ");
  Serial.print((float)DHT11.temperature, 2);
  humidity = DHT11.humidity;
  temperature = DHT11.temperature;

  //delay(2000);
  // Do nothing in the loop; the Arduino only responds when asked.
  delay(100);
}

// This function sends data to the Raspberry Pi when requested
void requestEvent() {

  int* buffer[2] = {humidity, temperature};
  Wire.write((char*) buffer, 4);
  
  //Serial.print("ping");  // Respond to master with this message
}

