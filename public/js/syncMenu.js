import { Routes } from "./routes.js";
import { PlayerIframe } from "./player-iframe.js";

export class SyncMenu {
    button;
    constructor() {
        let button = document.getElementById('nav-sync-button');
        if (!button) {
            return;
        }

        this.button = button;

        this.button.addEventListener('click', this.OpenSweetAlert.bind(this));
    }

    async OpenSweetAlert() {
        let syncPlayGroupList = await this.GetSyncplayList();
        // Generate list of groups
        let groupListElement = `<div class="flex flex-col gap-2 sync-modal-container">`;

        if (syncPlayGroupList.currentGroup) {
            groupListElement += `<ul>` + syncPlayGroupList.groups.filter((group) => group.id == syncPlayGroupList.currentGroup).map((group) => {
                return `<li>${group.name}</li>` + group.users.map((user) => `<li>${user}</li>`).join("")
            }).join("") + `</ul>`
                + `<button class="iframe-browser-button button button-primary" class="button button-primary">Go to Player</button>`
                + `<button class="leave-syncplay button button-danger">Leave Group</button>`
        } else {
            groupListElement += syncPlayGroupList.groups.map((group) => {
                return `<button class="join-syncplay button button-primary" data-group-id="${group.id}">${group.name}</button>`
            }).join("") + `<button class="create-syncplay button button-primary">Create Group</button>`
        }

        groupListElement += `</div>`;

        // Display modal
        Swal.fire({
            position: 'top-end',
            title: 'Syncplay',
            html: groupListElement,
            showConfirmButton: false,
            buttonsStyling: false,
            scrollbarPadding: false,
        })
        // Add event listeners
        $(".join-syncplay").click((e) => {
            let groupID = e.target.getAttribute("data-group-id");
            this.JoinSyncplayGroup(groupID, Swal.close);
        });
        $(".create-syncplay").click((e) => {
            this.CreateSyncplayGroup(Swal.close);

        })
        $(".leave-syncplay").click((e) => {
            this.LeaveSyncplayGroup(Swal.close);
        })
        $(".sync-modal-container").on("click", ".iframe-browser-button", () => {
            let playerIframe = new PlayerIframe();
            playerIframe.AddIframe();
            Swal.close();
        })
    }

    JoinSyncplayGroup(groupID, onOk) {
        fetch(Routes.ApiRoutes.SyncplayGroups, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                ID: groupID
            })
        }).then((res) => {
            if (res.ok) {
                onOk();
            }
        })
    }

    LeaveSyncplayGroup(onOk) {
        fetch(Routes.ApiRoutes.SyncplayGroups, {
            method: 'DELETE',
            headers: {
                'Content-Type': 'application/json'
            }
        }).then((res) => {
            if (res.ok) {
                onOk();
            }
        })
    }



    CreateSyncplayGroup(onOk) {
        fetch(Routes.ApiRoutes.SyncplayGroups, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            }
        }).then((res) => {
            if (res.ok) {
                onOk();
            }
        })

    }

    async GetSyncplayList() {
        const response = await fetch(Routes.ApiRoutes.SyncplayGroups).then((res) => res.json()).catch((err) => { console.error(err) });
        let groups = {
            currentGroup: response.currentGroup,
            groups: Object.keys(response.groups).map((key) => {
                return {
                    id: key,
                    name: response.groups[key].Name,
                    users: response.groups[key].Users
                }
            })
        }
        return groups;
    }


}