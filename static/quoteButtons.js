let quoteButtons = document.querySelectorAll(".quoteButton");

quoteButtons.forEach((button) => {
    button.addEventListener("click", function () {
        const formModal = document.getElementById("formModalContainer");
        if (formModal) formModal.style.display = "";

        const buttonClicked = document.getElementById("button_clicked");
        if (buttonClicked) buttonClicked.value = button.getAttribute("name");

        const popUp = document.getElementById("popUpModalOverlay");
        if (popUp) popUp.style.display = "none";
    });
});