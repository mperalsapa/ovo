
$(document).ready(function () {

    // Add listeners to buttons. TODO: This should be done in a more elegant way
    console.log("Document ready, adding listeners...")

    $("#addLibraryOpenModal").click("click", () => {
        $(".modalContainer").toggleClass("hidden");
    })

    $(".modalContainer").click("click", (e) => {
        if (e.target === e.currentTarget) {
            $(".modalContainer").toggleClass("hidden");
        }
    })

});

function OpenAddLibraryModal() {
    Swal.fire({
        title: "Add new Library",
        input: "text",
        inputAttributes: {
            autocapitalize: "off"
        },
        showCancelButton: true,
        confirmButtonText: "Create",
        showLoaderOnConfirm: true,
        preConfirm: async (login) => {
            try {
                const addLibEndpoint = `./`;
                const method = "POST";
                const libraryName = $("#addLibraryModalInput").val();
                const body = {
                    "libraryName": libraryName
                };
                const response = await fetch(addLibEndpoint, {
                    method: method,
                    headers: {
                        "Content-Type": "application/json"
                    },
                    body: JSON.stringify(body)
                });
                if (!response.ok) {
                    return Swal.showValidationMessage(`
                        ${JSON.stringify(await response.json())}
                        `);
                }
                return response.json();
            } catch (error) {
                Swal.showValidationMessage(`
                    Request failed: ${error}
                    `);
            }
        },
        allowOutsideClick: () => !Swal.isLoading()
    }).then((result) => {
        if (result.isConfirmed) {
            Swal.fire({
                title: `${result.value.login}'s avatar`,
                imageUrl: result.value.avatar_url
            });
        }
    })
}