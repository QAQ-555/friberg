import websocket
import threading
import json
import time
import requests

WS_URL = "ws://localhost:8001/friberg/socket"
HTTP_URL = "http://localhost:8001"

# 全局状态
state = {
    "uuid": None
}

# WebSocket 回调
def on_message(ws, message):
    try:
        data = json.loads(message)
        print("Received WS ["+str(data["msgType"])+"] message:",message)
        # 假设服务端返回 {"data": "uuid值"}
        if "data" in data and data["msgType"] == 0:
            state["uuid"] = data["data"]
            print("Updated UUID:", state["uuid"])
    except json.JSONDecodeError:
        # 直接是文本
        state["uuid"] = message
        print("Updated UUID:", state["uuid"])

def on_error(ws, error):
    print("WS Error:", error)

def on_close(ws, close_status_code, close_msg):
    print("WS Closed:", close_status_code, close_msg)

def on_open(ws):
    print("WS Opened")

def run_ws():
    rand = str(int(time.time()))
    url= f"{WS_URL}?user_name={rand}"
    ws = websocket.WebSocketApp(
        url,
        on_open=on_open,
        on_message=on_message,
        on_error=on_error,
        on_close=on_close
    )
    ws.run_forever()

# HTTP 请求函数
def create_room(uuid):
    url=f"{HTTP_URL}/friberg/room/create"
    payload = {"uuid": uuid}
    try:
        resp = requests.post(url, json=payload)
        print("HTTP Response:", resp.json())
        return resp.json()
    except Exception as e:
        print("HTTP Request error:", e)
        return None
def exit_room(uuid, room_id):
    url=f"{HTTP_URL}/friberg/room/exit"
    payload = {"uuid": uuid, "room_id": room_id}
    try:
        resp = requests.post(url, json=payload)
        print("HTTP Response:", resp.json())
        return resp.json()
    except Exception as e:
        print("HTTP Request error:", e)
        return None
# 主线程处理用户输入和逻辑
def main_loop():
    while True:
        time.sleep(1)
        if state["uuid"]:
            cmd = input("Enter command (create/join/quit): ").strip()
            if cmd == "create":
                create_room(state["uuid"])
            elif cmd == "quit":
                
                break
            else:
                print("Unknown command")

if __name__ == "__main__":
    # 启动 WS 线程
    wst = threading.Thread(target=run_ws)
    wst.daemon = True
    wst.start()

    # 主循环处理输入和 HTTP
    main_loop()
