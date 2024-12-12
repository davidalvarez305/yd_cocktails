const facebookLeadEventName = document.getElementById("facebookLeadEventName").textContent;
const linkedInEventID = document.getElementById("linkedInEventID").textContent;
const quoteEventName = document.getElementById("quoteEventName").textContent;

const phoneNumber = document.getElementById("phoneNumberCTA");

if (phoneNumber) phoneNumber.addEventListener("click", () => handlePhoneNumberClick())

function handlePhoneNumberClick() {
    // Report FB Conversion
    if (fbq) fbq("track", facebookLeadEventName);

    // Report LinkedIn
    if (window.lintrk) window.lintrk('track', { conversion_id: linkedInEventID });

    // Report Microsoft Ads Conversion
    window.uetq = window.uetq || [];

    window.uetq.push('event', quoteEventName, {});
}