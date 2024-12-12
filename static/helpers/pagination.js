export class PaginationHandler {
    constructor() {
        this.querystring = new URLSearchParams(window.location.search);
        this.maxPages = parseInt(document.getElementById('maxPages').textContent);
        this.bindPagination();
    }

    bindPagination() {
        document.querySelectorAll('.pagination-link').forEach(item => {
            item.addEventListener('click', e => this.handlePageChange(e, item));
        });
    }

    handlePageChange(event, item) {
        const currentPage = this.querystring.get('page_num') ?? '1';
        const page = parseInt(currentPage);
        const chevronValue = item.getAttribute('name');

        if (chevronValue === "left" && page > 1) {
            this.querystring.set('page_num', page - 1);
        } else if (chevronValue === "right" && page < this.maxPages && this.maxPages > 1) {
            this.querystring.set('page_num', page + 1);
        }

        if (this.querystring.toString()) {
            this.updateURL();
        }
    }

    updateURL() {
        const { origin, pathname } = window.location;
        const url = new URL(origin + pathname);
        url.search = this.querystring.toString();

        window.location.replace(url.href);
    }
}