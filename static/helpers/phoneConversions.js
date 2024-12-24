const facebookLeadEventName = document.getElementById("facebookLeadEventName").textContent;

const phoneNumbers = document.querySelectorAll(".phoneNumberCTA");

phoneNumbers.forEach(phoneNumber => {
    phoneNumber.addEventListener("click", () => handlePhoneNumberClick())
})

function handlePhoneNumberClick() {
    if (fbq) fbq("track", facebookLeadEventName);
}