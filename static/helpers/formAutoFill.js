export const formAutofill = () => {
    const params = new URLSearchParams(window.location.search);

    for (const [key, value] of Object.entries(params)) {
        const field = document.querySelector(key);

        if (field && field.closest('form')) {
            field.value = value;
        }
    }
}