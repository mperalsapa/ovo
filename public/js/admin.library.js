
$(document).ready(function () {

    // Add listeners to buttons. TODO: This should be done in a more elegant way
    console.log("Document ready, adding listeners...")

    $("#addLibraryOpenModal").click("click", () => {
        EditLibrary(0)
    })

    $(".modalContainer").click("click", (e) => {
        if (e.target === e.currentTarget) {
            $(".modalContainer").toggleClass("hidden");
        }
    })

    $("#libraryContainer").find("button").on("click", (e) => {
        EditLibrary(e.currentTarget.getAttribute("data-ID"))
    })

    $("#addPathToForm").click("click", () => {
        let path = $("#newPath").val()
        AddPathToLibrary(path)
    })

    $("#submit").click("click", () => {
        SaveLibrary()
    })

});


async function EditLibrary(id) {
    console.log("Edit Library: " + id)
    let libraryData = await fetch(`/ovo/api/library/${id}`)
    libraryData = await libraryData.json()
    OpenDialog(libraryData)
}

function OpenDialog(data) {
    console.log(data)
    $("#library_id").val(data.ID)
    $("#name").val(data.name)
    SetPathList(data.paths)

    $(".modalContainer").toggleClass("hidden");
}

function SaveLibrary() {
    let id = $("#library_id").val()
    let name = $("#name").val()
    let type = $("#libraryType").val()
    let paths = $("#pathList").find(".pathElement").find("span").map((i, e) => e.innerHTML).get()

    let data = {
        "ID": parseInt(id),
        "name": name,
        "type": type,
        "paths": paths
    }

    console.log("Saving library: " + id)
    console.log(data)

    fetch(`/ovo/api/library/${id}`, {
        method: "POST",
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(data)
    }).then((response) => {
        console.log(response)
        if (response.status == 200) {
            location.reload()
        }
    })
}

function SetPathList(paths) {
    $("#pathList").empty()
    if (paths == null || paths == undefined) {
        $("#pathList").hasClass("hidden") ? null : $("#pathList").toggleClass("hidden")
        return
    }

    paths.forEach(path => {
        AddPathToLibrary(path)
    })

    // Add erase button listener
    $(".pathElement").find("button").on("click", (e) => {
        e.currentTarget.parentElement.remove()
    })
}

function AddPathToLibrary(path) {
    if (path == "") {
        return
    }

    if ($("#pathList").hasClass("hidden")) {
        $("#pathList").toggleClass("hidden")
    }

    let pathElement = document.createElement("div")
    let pathText = document.createElement("span")
    let pathDelete = document.createElement("button")

    pathText.innerHTML = path
    pathDelete.classList.add("button", "button-danger")
    pathDelete.innerHTML = "-"
    pathElement.classList.add("pathElement")
    pathElement.appendChild(pathText)
    pathElement.appendChild(pathDelete)

    $("#pathList").append(pathElement)

    $("#newPath").val("")
}