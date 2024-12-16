function calculateBartendersNeeded(x) {
    if (x <= 50) {
        return 1;
    } else if (x >= 41 && x <= 80) {
        return 2;
    } else if (x >= 81 && x <= 120) {
        return 3;
    } else if (x >= 121) {
        return 4;
    }
}

function calculateBartenderRate(numBartenders, hours) {
    return parseFloat(document.getElementById("bartendingRate").value) * numBartenders * hours;
}

function validateGuestsInput(value) {
    const guests = Number(value);

    if (isNaN(guests) || guests <= 0) {
        return false;
    } else {
        return true;
    }
}

function calculateTotalCost(numBarRentals, willProvideAlcohol, willProvideMixers, willProvideJuices, willProvideSoftDrinks, willProvideCups, willProvideIce, guests) {
    let totalCost = 0;

    totalCost += numBarRentals * 150; // $150 per bar rental
    totalCost += numBarRentals * 50 * 4; // $50 * 4 for breakdown & setup fee

    switch (true) {
        case willProvideAlcohol === false:
            totalCost += guests * 15; // $15 per guest for alcohol
            break;

        case willProvideMixers === false:
            totalCost += guests * 3; // $3 per guest for mixers
            break;

        case willProvideJuices === false:
            totalCost += guests * 2; // $2 per guest for juices
            break;

        case willProvideSoftDrinks === false:
            totalCost += guests * 2.50; // $2.50 per guest for soft drinks
            break;

        case willProvideCups === false:
            totalCost += guests * 2; // $2 per guest for cups
            break;

        case willProvideIce === false:
            totalCost += guests * 2; // $2 per guest for ice
            break;

        default:
            break;
    }

    return totalCost;
}