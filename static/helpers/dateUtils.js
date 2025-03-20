export function convertTimestampToEst(timestamp) {
    const date = new Date(timestamp * 1000);
    
    const estDate = date.toLocaleString('en-US', {
        timeZone: 'America/New_York',
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit',
        hour12: false,
    });

    const [month, day, year, hour, minute] = estDate.replace(',', '').split(/[\s/:]/);

    return `${year}-${month}-${day}T${hour}:${minute}`;
}