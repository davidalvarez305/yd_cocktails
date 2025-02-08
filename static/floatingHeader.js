const floatingHeader = document.getElementById("floatingHeader");

if (!floatingHeader) return;

floatingHeader.style.display = "none";

// Show floating header After 25% scroll
function handleScroll() {
    const scrollTop = window.scrollY || document.documentElement.scrollTop;
    const pageHeight = document.documentElement.scrollHeight - document.documentElement.clientHeight;

    const scrollPercentage = (scrollTop / pageHeight) * 100;

    if (scrollPercentage > 5 && floatingHeader) {
        floatingHeader.style.display = "";
    } else {
        floatingHeader.style.display = "none";
    }
}

window.addEventListener("scroll", handleScroll);