const phoneNumbers = document.querySelectorAll(".phoneNumberCTA");

phoneNumbers.forEach(phoneNumber => {
    phoneNumber.addEventListener("click", () => handlePhoneNumberClick())
})

function handlePhoneNumberClick() {
    if (fbq) fbq("track", "Lead");
}