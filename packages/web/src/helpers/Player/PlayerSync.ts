window.PlayerjsEvents = function (event, id, info) {
    console.log("Event received:", event, info);
    switch (event) {
        case "play":
            localStorage.setItem("player-sync", JSON.stringify({ type: "play" }));
            break;
        case "pause":
            localStorage.setItem("player-sync", JSON.stringify({ type: "pause" }));
            break;
        case "time":
            localStorage.setItem("player-sync", JSON.stringify({ type: "seek", value: info }));
            break;
        case "volume":
            localStorage.setItem("player-sync", JSON.stringify({ type: "volume", value: info }));
            break;
        default:
            break;
    }
};