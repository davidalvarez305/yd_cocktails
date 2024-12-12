const SCROLL_THRESHOLD = 50;
const TRIGGER_THRESHOLD = 60000;
let hasScrolledPast75 = false;
let hasScrolledUpAbove50 = false;
let alreadyPoppedUp = false;

function showModal() {
	fetch("/partials/pop-up-modal")
		.then((response) => {
			if (response.ok) {
				return response.text();
			} else {
				throw new Error("Error: " + response.statusText);
			}
		})
		.then((modalHTML) => {
			const modalOverlay = document.getElementById("popUpModalOverlay");
			modalOverlay.outerHTML = modalHTML;
			modalOverlay.style.display = "flex";

			const modalContent = document.getElementById("popUpModalContent");

			// Prevent clicks inside the modal content from closing the modal
			modalContent.addEventListener("click", (e) => {
				e.stopPropagation();
			});

			// Hide the modal when clicking outside the modal content
			const modal = document.getElementById("popUpModalOverlay");
			modal.addEventListener("click", () => {
				modal.outerHTML = "";
				alreadyPoppedUp = true;
			});

			// Bring form to focus
			applyButtonlogic();
		})
		.catch(console.error);
}

// Trigger after time
window.addEventListener("load", () => {
	setTimeout(() => {
		if (!alreadyPoppedUp) showModal();
	}, TRIGGER_THRESHOLD);
});

// Trigger after scroll
window.addEventListener("scroll", function () {
	const pageHeight = document.documentElement.scrollHeight;
	const scrollPercent = window.scrollY / (pageHeight - window.innerHeight);

	if (scrollPercent >= 0.75) {
		hasScrolledPast75 = true;
	}

	// Check if scrolled back above 50% and conditions are met
	if (hasScrolledPast75 && scrollPercent < 0.50) {
		if (!hasScrolledUpAbove50 && !alreadyPoppedUp) {
			showModal();
			hasScrolledUpAbove50 = true; // Ensure the alert only fires once
		}
	}
});
