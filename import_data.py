import csv
import requests

with open("./data/weather.dat", newline="") as csvfile:
    reader = csv.reader(csvfile, delimiter="	")
    for row in reader:
        r = requests.post(
            "http://localhost:8080/api/data/weather",
            json={
                "date": row[0],
                "humidity": float(row[1]),
                "temperature": float(row[2]),
            },
        )
        print(f"Status Code: {r.status_code}")
