const phoneNumbers = document.querySelectorAll(".phoneNumberCTA");

phoneNumbers.forEach(phoneNumber => {
    phoneNumber.addEventListener("click", () => handlePhoneNumberClick())
})

function handlePhoneNumberClick() {
    if (fbq) fbq("track", "Lead");
    if (gtag) gtag("event", "generate_lead", { currency: "USD", value: 150.00 });
}