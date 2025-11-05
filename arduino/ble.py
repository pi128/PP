from machine import ADC, Pin
from time import sleep
import bluetooth
from ble_simple_peripheral import BLESimplePeripheral  

adc = ADC(Pin(3))
adc.atten(ADC.ATTN_11DB)
ble = bluetooth.BLE()
sp = BLESimplePeripheral(ble, name="NanoESP32_ADC")

while True:
    voltage = adc.read() * (3.3 / 4095)
    msg = "{:.3f}\n".format(voltage)
    if sp.is_connected():
        sp.send(msg)
    sleep(0.05)

    