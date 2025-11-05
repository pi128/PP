import network, socket, time
from machine import ADC, Pin

SSID = "wifi"
PASS = "password"
PORT = 12345

wlan = network.WLAN(network.STA_IF)
wlan.active(True)
if not wlan.isconnected():
    wlan.connect(SSID, PASS)
    print("Connecting to Wi-Fi...")
    while not wlan.isconnected():
        time.sleep(0.2)

print("Connected! IP:", wlan.ifconfig()[0])

adc = ADC(Pin(3))
adc.atten(ADC.ATTN_11DB)

srv = socket.socket()
srv.bind(("", PORT))
srv.listen(1)
print("Listening on port", PORT)

while True:
    conn, addr = srv.accept()
    print("Client connected:", addr)
    try:
        while True:
            voltage = adc.read() * (3.3 / 4095)
            msg = "{:.3f}\n".format(voltage)
            conn.send(msg.encode())
            time.sleep(0.1)
    except OSError:
        conn.close()
        print("Client disconnected")