import EventEmitter from "events";

export default class Connection extends EventEmitter {
  socket = null;
  isOpen = false;
  connect() {
    var loc = window.location,
      socketUrl;
    if (loc.protocol === "https:") {
      socketUrl = "wss:";
    } else {
      socketUrl = "ws:";
    }
    socketUrl += "//" + loc.host;
    socketUrl += "/ws";
    if (socketUrl.indexOf("8080") !== -1) {
      //local development
      socketUrl = "ws://localhost:4444/ws";
    }
    this.socket = new WebSocket(socketUrl);
    this.socket.addEventListener("open", ev => {
      this.emit("statusChanged", "connected");
      this.emit("open", ev);
      this.isOpen = true;
    });
    this.socket.addEventListener("error", e => {
      // eslint-disable-next-line
      console.error(e);
      this.emit("statusChanged", "error");
    });
    this.socket.addEventListener("close", () => {
      this.emit("statusChanged", "closed");
    });
    this.socket.addEventListener("message", ev => {
      let data = JSON.parse(ev.data);
      this.emit("message", data.type, data.data);
      this.emit("message." + data.type, data.data);
    });
  }
  sendMessage(type, data) {
    if (!this.isOpen) {
      this.connect("open", () => {
        this.sendMessage(type, data);
      });
      return;
    }
    this.socket.send(
      JSON.stringify({
        type,
        data
      })
    );
  }
}
