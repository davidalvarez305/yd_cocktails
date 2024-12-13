const facebookLeadEventName = document.getElementById("facebookLeadEventName").textContent;

const phoneNumber = document.getElementById("phoneNumberCTA");

if (phoneNumber) phoneNumber.addEventListener("click", () => handlePhoneNumberClick())

function handlePhoneNumberClick() {
    if (fbq) fbq("track", facebookLeadEventName);
}