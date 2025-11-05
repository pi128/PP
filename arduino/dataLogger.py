import socket
import csv
from datetime import datetime

ESP32_IP = "10.0.0.195"  # <-- replace with your board's IP
PORT = 12345
OUTFILE = "log.csv"

with socket.create_connection((ESP32_IP, PORT)) as s, open(OUTFILE, "w", newline="") as f:
    writer = csv.writer(f)
    writer.writerow(["timestamp", "voltage_V"])
    print("Connected to ESP32. Logging to", OUTFILE)
    try:
        for line in s.makefile():
            try:
                voltage = float(line.strip())
                writer.writerow([datetime.now().isoformat(), voltage])
                f.flush()
                print(voltage)
            except ValueError:
                pass
    except KeyboardInterrupt:
        print("\nStopped logging.")