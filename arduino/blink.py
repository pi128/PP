""" 

   To run this script on your Arduino:
    to find the device ls /dev/cu.*   
    source venv/bin/activate
    mpremote connect /dev/cu.usbmodem101 cp blink.py :
    mpremote connect /dev/cu.usbmodem101 exec 'import blink'
    
"""


try:
    from machine import Pin
    import time


except ImportError:


    exit(1)

led = Pin(48, Pin.OUT)

def blink(times=5, delay=0.5):
    print(f'Blinking LED {times} times...')
    for i in range(times):
        led.on()
        time.sleep(delay)
        led.off()
        time.sleep(delay)
    print('Done!')

if __name__ == '__main__':
    blink(times=10, delay=0.5)
