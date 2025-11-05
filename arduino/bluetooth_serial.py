import bluetooth
import struct
import time
from micropython import const

_IRQ_CENTRAL_CONNECT = const(1)
_IRQ_CENTRAL_DISCONNECT = const(2)
_IRQ_GATTS_WRITE = const(3)

_UART_UUID = bluetooth.UUID("6E400001-B5A3-F393-E0A9-E50E24DCCA9E")
_UART_TX = (bluetooth.UUID("6E400003-B5A3-F393-E0A9-E50E24DCCA9E"), bluetooth.FLAG_NOTIFY,)
_UART_RX = (bluetooth.UUID("6E400002-B5A3-F393-E0A9-E50E24DCCA9E"), bluetooth.FLAG_WRITE,)
_UART_SERVICE = (_UART_UUID, (_UART_TX, _UART_RX,),)

class BLESimplePeripheral:
    def __init__(self, ble, name="ESP32"):
        self._ble = ble
        self._ble.active(True)
        self._ble.irq(self._irq)
        ((self._handle_tx, self._handle_rx,),) = self._ble.gatts_register_services((_UART_SERVICE,))
        self._connections = set()
        self._payload = bytearray()
        self._advertise(name)

    def _irq(self, event, data):
        if event == _IRQ_CENTRAL_CONNECT:
            conn_handle, _, _ = data
            print("Connected")
            self._connections.add(conn_handle)
        elif event == _IRQ_CENTRAL_DISCONNECT:
            conn_handle, _, _ = data
            print("Disconnected")
            self._connections.remove(conn_handle)
            self._advertise()
        elif event == _IRQ_GATTS_WRITE:
            conn_handle, value_handle = data
            value = self._ble.gatts_read(value_handle)
            if value_handle == self._handle_rx and self.on_write:
                self.on_write(value)

    def send(self, data):
        for conn_handle in self._connections:
            self._ble.gatts_notify(conn_handle, self._handle_tx, data)

    def _advertise(self, name="ESP32"):
        self._ble.gap_advertise(100000, adv_data=self._advertising_payload(name=name, services=[_UART_UUID]))

    def on_write(self, value):
        pass

    @staticmethod
    def _advertising_payload(limited_disc=False, br_edr=False, name=None, services=None):
        payload = bytearray()
        def _append(adv_type, value):
            nonlocal payload
            payload += struct.pack("BB", len(value) + 1, adv_type) + value
        _append(0x01, struct.pack("B", (0x01 if limited_disc else 0x02) + (0x18 if br_edr else 0x04)))
        if name:
            _append(0x09, name.encode())
        if services:
            for uuid in services:
                b = bytes(uuid)
                _append(0x03 if len(b) == 2 else 0x07 if len(b) == 16 else 0x06, b)
        return payload


def main():
    ble = bluetooth.BLE()
    p = BLESimplePeripheral(ble, name="Arduino")
    
    def on_rx(data):
        print("Received:", data.decode().strip())
        message = data.decode().strip()
        
        if message == "hello":
            response = "Hello from Arduino!\n"
            p.send(response.encode())
        elif message == "blink":
            from machine import Pin
            led = Pin(48, Pin.OUT)
            for _ in range(3):
                led.on()
                time.sleep(0.2)
                led.off()
                time.sleep(0.2)
            response = "LED blinked!\n"
            p.send(response.encode())
        else:
            response = f"You sent: {message}\n"
            p.send(response.encode())
    
    p.on_write = on_rx
    print("Bluetooth UART active. Connect with a BLE terminal app.")
    print("Device name: Arduino")
    
    while True:
        time.sleep(1)

if __name__ == "__main__":
    main()

