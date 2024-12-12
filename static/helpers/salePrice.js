const TAX_RATE = 0.07;
const CREDIT_CARD_FEE = 0.0595;
const PROFIT_MARGIN = 0.50;

export const calculateSalePrice = (unitCost, hasCommission, commissionRate) => {
    // Step 1: Calculate the cost after tax
    let costWithTax = unitCost * (1 + TAX_RATE);

    // Step 2: Calculate the initial sale price to achieve a 50% profit margin (before transaction fee and commission)
    let salePriceBeforeFee = costWithTax / PROFIT_MARGIN;

    // Step 3: Adjust the sale price to account for the transaction fee
    let salePrice = salePriceBeforeFee / (1 - CREDIT_CARD_FEE);

    // If there's no commission, return the sale price
    if (!hasCommission) return salePrice.toFixed(2);

    // Step 4: Calculate the credit card fee for the transaction
    let creditCardFee = salePrice * CREDIT_CARD_FEE;

    // Step 5: Define the target net profit after all costs and commission
    let adjustedSalePrice = (costWithTax + creditCardFee) / (PROFIT_MARGIN * (1 - commissionRate));

    // Return the adjusted sale price rounded to two decimal places
    return adjustedSalePrice.toFixed(2);
}