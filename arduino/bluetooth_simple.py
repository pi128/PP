import bluetooth
import time
from machine import Pin

ble = bluetooth.BLE()
ble.active(True)

led = Pin(48, Pin.OUT)

def on_rx(event, data):
    if event == 3:
        print("Received data")

print("Bluetooth enabled")
print("Device name: ESP32")

while True:
    time.sleep(1)

