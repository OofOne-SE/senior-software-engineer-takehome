import websocket
import json

up_data = []


def on_open(ws):
    print("Connected to", ws.url)


def on_message(ws, message):
    data = json.loads(message)
    up_data.append(data)
    print(up_data)


def on_close(ws, code, msg):
    print("Connection closed")


def on_error(ws, error):
    print("Error:", error)


def Connect():
    wsURL = f"ws://0.0.0.0:8080/api/ws"
    ws = websocket.WebSocketApp(
        wsURL,
        on_open=on_open,
        on_message=on_message,
        on_close=on_close,
        on_error=on_error,
    )
    ws.run_forever()


Connect()
