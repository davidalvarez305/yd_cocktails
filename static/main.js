function handleCloseAlertModal() {
	const closeModalButton = document.querySelectorAll(".closeModal");
	const modal = document.getElementById("alertModal");

	closeModalButton.forEach((button) => {
		button.addEventListener("click", () => {
			modal.style.display = "none";
		});
	});
}

function setUser() {
    const user = {
        landingPage: typeof window !== "undefined" && window.location ? window.location.href : null,
        referrer: typeof document !== "undefined" && document.referrer ? document.referrer : null,
    };

    localStorage.setItem("user", JSON.stringify(user));
}

function preserveQuerystring() {
	if (["crm"].some(link => window.location.pathname.includes(link))) return;

	const links = document.querySelectorAll('a');

	links.forEach(link => {
		link.addEventListener('click', function () {
			const qs = window.location.search;
			if (qs.length > 0) {
				const url = new URL(this.href, window.location.origin);
				url.search = url.search ? `${url.search}&${qs.substring(1)}` : qs;
				this.href = url.toString();
			}
		});
	});
}

document.addEventListener("DOMContentLoaded", () => {
	// Preserve querystring
	preserveQuerystring();

	// Get user variables
	localStorage.getItem("user") ?? setUser();
});