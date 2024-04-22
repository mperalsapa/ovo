$(document).ready(function () {
    // add listener to any input text field inside pathlist
    $("#pathList").find("input[type='text']").on("input", PathListener);
});

function PathListener() {
    if (GetEmptyPathCount() === 0) {
        console.log("Adding new input field");
        let pathChildren = $("#pathList").find("input[type='text']");
        let lastPath = pathChildren[pathChildren.length - 1];
        let newPath = $(lastPath).clone();
        newPath.val("");
        pathList.append(newPath);
        newPath.on("input", PathListener);
    }

    if ($(this).val() === "" && GetEmptyPathCount() > 1) {
        console.log("Removing input field");
        let pathChildren = $("#pathList").find("input[type='text']");
        if (pathChildren.length > 1) {
            $(this).remove();
        }
    }
}

function GetEmptyPathCount() {
    let pathList = $("#pathList");
    let pathChildren = pathList.find("input[type='text']");
    let emptyCount = 0;
    for (let i = 0; i < pathChildren.length; i++) {
        if (pathChildren[i].value === "") {
            emptyCount++;
        }
    }
    return emptyCount;
}