const querystring = new URLSearchParams(window.location.search);

function handleBindPagination() {
    document.querySelectorAll('.pagination-link').forEach(item => {
        item.addEventListener('click', e => {
            const currentPage = querystring.get('page_num') ?? '1';
            const page = parseInt(currentPage);
            const maxPages = parseInt(document.getElementById('maxPages').textContent);
            const chevronValue = item.getAttribute('name');

            if (chevronValue === "left" && page > 1) {
                querystring.set('page_num', page - 1);
            }

            if (chevronValue === "right" && page < maxPages && maxPages > 1) {
                querystring.set('page_num', page + 1);
            }

            if (!querystring.toString()) return;
            updateURL();
        });
    });
}

function updateURL() {
    const { origin, pathname } = window.location;
    const url = new URL(origin + pathname);
    url.search = querystring.toString();

    window.location.replace(url.href);
}

document.addEventListener('DOMContentLoaded', () => handleBindPagination());