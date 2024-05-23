import { Routes } from "./routes.js";

export class PlayerIframe {
    frame;

    constructor() { }
    // Add the iframe to the page
    // If we click on "play" from the item info, we can pass "newItem" as true to force a new item to the sync player.
    // Otherwise we open the player without sending new item (requesting the existing one as an empty path)
    AddIframe(newItemID = 0) {
        if (this.frame) {
            return;
        }
        let ifrm = document.createElement("iframe");

        // if (ItemID) {
        if (newItemID != 0) {
            ifrm.setAttribute("src", Routes.Routes.Player + "?item=" + newItemID);
        } else {
            ifrm.setAttribute("src", Routes.Routes.Player + "");
        }

        ifrm.style.width = "100vw";
        ifrm.style.height = "100vh";
        ifrm.style.position = "fixed";
        ifrm.style.top = "0";
        ifrm.style.left = "0";
        ifrm.style.zIndex = "+2";

        ifrm.classList.add("player-iframe");

        // Add the iframe to the page
        document.body.appendChild(ifrm);
        this.frame = ifrm;
        this.frame.addEventListener("load", this.OnIframeLoad.bind(this));
    }

    OnIframeLoad() {
        this.frame.contentWindow.document.querySelector(".close-video-iframe").addEventListener("click", this.CloseIframe.bind(this));
    }

    CloseIframe() {
        if (this.frame) {
            this.frame.remove();
        }
    }

}