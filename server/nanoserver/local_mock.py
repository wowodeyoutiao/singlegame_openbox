# -*- coding: utf-8 -*-
"""Local go-admin mock for single-server dev.

Listens on 127.0.0.1:18080.
  POST /goadminapi/api/v1/servername
    - empty body (server admin fetchServerList) -> ServerListResponse
    - body with pageIndex (client OnGetSeverInfo) -> HTTP_SeverResponse
  POST /goadminapi/api/v1/servername/version -> HTTP_VersionRequest

Run:  python local_mock.py
"""
import json
from http.server import BaseHTTPRequestHandler, ThreadingHTTPServer

HOST = "127.0.0.1"
PORT = 18080
SERVER_ID = 9
HTTP_BASE = "http://127.0.0.1:12307"
WS_ADDR = "127.0.0.1:7000"

ADMIN_SERVER = {
    "id": SERVER_ID,
    "name": "LocalServer",
    "ipaddr": WS_ADDR,
    "server_id": SERVER_ID,
    "merge_sid": 0,
    "open_by": 0,
    "status": 1,
    "is_show": 1,
    "createdAt": "",
    "updatedAt": "",
    "createBy": 0,
    "updateBy": 0,
    "zone_id": SERVER_ID,
}

CLIENT_SEVER_ITEM = {
    "id": SERVER_ID,
    "name": "LocalServer",
    "ipaddr": HTTP_BASE,
    "status": 1,
    "server_id": SERVER_ID,
    "open_by": 0,
    "websocket": WS_ADDR,
}


class Handler(BaseHTTPRequestHandler):
    def _send(self, payload):
        data = json.dumps(payload, ensure_ascii=False).encode("utf-8")
        self.send_response(200)
        self.send_header("Content-Type", "application/json; charset=utf-8")
        self.send_header("Content-Length", str(len(data)))
        self.end_headers()
        self.wfile.write(data)

    def do_POST(self):
        length = int(self.headers.get("Content-Length", 0) or 0)
        body = self.rfile.read(length).decode("utf-8", "replace") if length else ""
        path = self.path.split("?")[0]

        if path.endswith("/version"):
            payload = {"code": 200, "msg": "",
                       "data": {"version": "0", "created_at": 0}}
            tag = "version"
        elif "pageIndex" in body:
            payload = {"code": 200, "msg": "",
                       "data": {"count": 1, "pageIndex": 1, "pageSize": 20,
                                "list": [CLIENT_SEVER_ITEM],
                                "mylist": [CLIENT_SEVER_ITEM],
                                "version_status": 0}}
            tag = "client-serverlist"
        else:
            payload = {"code": 0, "error": "",
                       "data": {"count": 1, "list": [ADMIN_SERVER]}}
            tag = "admin-serverlist"

        print("[mock] %s %s -> %s" % (self.command, path, tag), flush=True)
        self._send(payload)

    def do_GET(self):
        self.do_POST()

    def log_message(self, fmt, *args):
        pass


if __name__ == "__main__":
    print("[mock] listening on %s:%d (server_id=%d)" % (HOST, PORT, SERVER_ID), flush=True)
    ThreadingHTTPServer((HOST, PORT), Handler).serve_forever()
