export function handleQuoteService(field, guests, hours, assumedBaseHours, handleFindAdHocUnitsByServiceId) {
    if (!field || !field.dataset) {
        console.error("Error: 'field' or 'field.dataset' is missing.");
        return;
    }

    let serviceId = parseInt(field.dataset.serviceId, 10);
    let suggestedPrice = parseFloat(field.dataset.suggestedPrice);
    let unitTypeId = parseInt(field.dataset.unitTypeId, 10);
    let serviceTypeId = parseInt(field.dataset.serviceTypeId, 10);
    let guestRatio = parseFloat(field.dataset.guestRatio);

    if (isNaN(serviceId) || isNaN(suggestedPrice) || isNaN(unitTypeId) || isNaN(serviceTypeId) || isNaN(guestRatio)) {
        console.error("Error: One or more required dataset values are missing or invalid.", {
            serviceId, suggestedPrice, unitTypeId, serviceTypeId, guestRatio
        });
        return;
    }

    if (!guests || isNaN(guests)) {
        console.error("Error: 'guests' parameter is missing or invalid.", guests);
        return;
    }

    if (!hours || isNaN(hours)) {
        console.error("Error: 'hours' parameter is missing or invalid.", hours);
        return;
    }

    if (!assumedBaseHours || isNaN(assumedBaseHours) || assumedBaseHours <= 0) {
        console.error("Error: 'assumedBaseHours' is missing, invalid, or zero.", assumedBaseHours);
        return;
    }

    let units = 0;
    const PER_PERSON_UNIT_TYPE = 1;
    const HOURLY_UNIT_TYPE = 2;
    const RATIO = 4;
    const AD_HOC = 5;
    const FIXED = 6;
    const HOURLY_SERVICE_TYPE_ID = 5;

    switch (unitTypeId) {
        case PER_PERSON_UNIT_TYPE:
            units = guests;
            suggestedPrice *= (hours / assumedBaseHours);
            break;
        case HOURLY_UNIT_TYPE:
            units = hours;
            break;
        case RATIO:
            units = guestRatio > 0 ? Math.ceil(guests / guestRatio) : 0;
            if (serviceTypeId === HOURLY_SERVICE_TYPE_ID && hours > 0) units *= hours;
            break;
        case AD_HOC:
            if (typeof handleFindAdHocUnitsByServiceId === "function") {
                let selector = field.dataset.selector || "";
                units = handleFindAdHocUnitsByServiceId(serviceId, selector);
            } else {
                console.error("Error: 'handleFindAdHocUnitsByServiceId' is not a function.");
                return;
            }
            break;
        case FIXED:
            units = 1;
            break;
        default:
            console.error("Error: Unknown unitTypeId:", unitTypeId);
            return;
    }

    return { suggestedPrice, units };
}