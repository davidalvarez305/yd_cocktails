import { PaginationHandler } from "./pagination";

export class NotificationHelper {
    constructor(element) {
        this.notificationModal = element || document.getElementById("alertModal");
    }

    handleNotification(notification) {
        let modalMessage = notification;

        if (!err.includes("alertModal")) {
            fetch(`/partials/error-modal?err=${err}`)
                .then(response => {
                    if (response.ok) {
                        return response.text();
                    } else {
                        return response.text().then((err) => {
                            throw new Error(err);
                        });
                    }
                })
                .then(err => {
                    modalMessage = err;
                })
                .catch(console.error);
        }

        if (this.notificationModal) {
            this.notificationModal.outerHTML = modalMessage

            this.handleCloseAlertModal();
        };
    }

    handleCloseAlertModal() {
        const closeModalButton = document.querySelector(".closeModal");

        if (!closeModalButton || !this.notificationModal) return;

        closeModalButton.addEventListener("click", () => {
            this.notificationModal.style.display = "none";
        });
    }
}

export class FormHandler {
    constructor(form) {
        this.data = new FormData(form);
        this.notification = new NotificationHelper();
        this.pagination = new PaginationHandler();
    }

    handleSubmitForm(url, headers, handleSuccess, handleError) {
        fetch(url, { body: this.data, ...headers })
            .then((response) => {
                const token = response.headers.get('X-Csrf-Token');
                if (token) {
                    const csrfInputs = document.querySelectorAll('[name="csrf_token"]');

                    csrfInputs.forEach(input => input.value = token);
                }
                if (response.ok) {
                    return response.text();
                } else {
                    return response.text().then((err) => {
                        throw new Error(err);
                    });
                }
            })
            .then(html => {
                if (!handleSuccess) {
                    this.notification.handleNotification(html);
                } else {
                    handleSuccess(html);
                }
            })
            .catch(err => {
                if (!handleError) {
                    this.notification.handleNotification(html);
                    return;
                } else {
                    handleError(err)
                }
            });
    }

    swapTable(html, elementId) {
        const table = document.getElementById(elementId);

        if (!table) return;

        table.outerHTML = html;

        this.pagination.bindPagination();
    }
}