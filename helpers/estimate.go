package helpers

import (
	"github.com/davidalvarez305/yd_cocktails/constants"
	"github.com/davidalvarez305/yd_cocktails/types"
)

func calculateBartendersNeeded(guests int) int {
	if guests <= 50 {
		return 1
	} else if guests >= 51 && guests <= 90 {
		return 2
	} else if guests >= 91 && guests <= 120 {
		return 3
	} else if guests >= 121 {
		return 4
	}
	return 0
}

func calculateBartenderRate(numBartenders, hours float64) float64 {
	return constants.BartendingRate * numBartenders * hours
}

func CalculatePackagePrice(form types.EstimateForm) float64 {
	var totalCost float64

	guests := SafeInt(form.Guests)
	if guests <= 0 {
		return 0.00
	}

	floatGuests := float64(guests)

	hours := SafeInt(form.Hours)

	bartendersNeeded := calculateBartendersNeeded(guests)

	bartenderRate := calculateBartenderRate(float64(bartendersNeeded), float64(hours))

	totalCost += bartenderRate

	packageTypeID := SafeInt(form.PackageTypeID)

	numBarsFloat := float64(SafeInt(form.NumBars))
	totalCost += numBarsFloat * constants.BarRentalCost

	totalCost += numBarsFloat * constants.MobileBarFee * constants.TimeToSetUpAndBreakDown

	if form.WillProvideLiquor != nil && !*form.WillProvideLiquor {
		totalCost += floatGuests * constants.PerPersonAlcoholFee
	}
	if form.WillProvideBeerAndWine != nil && !*form.WillProvideBeerAndWine {
		totalCost += floatGuests * constants.PerPersonBeerAndWineFee
	}
	if packageTypeID == constants.FullOpenBarPackageTypeID {
		totalCost += floatGuests * constants.FullOpenBarFee
	}
	if packageTypeID == constants.PartialOpenBarPackageTypeID {
		totalCost += floatGuests * constants.PartialOpenBarFee
	}
	if form.WillProvideMixers != nil && !*form.WillProvideMixers {
		totalCost += floatGuests * constants.PerPersonMixersFee
	}
	if form.WillProvideJuices != nil && !*form.WillProvideJuices {
		totalCost += floatGuests * constants.PerPersonJuicesFee
	}
	if form.WillProvideSoftDrinks != nil && !*form.WillProvideSoftDrinks {
		totalCost += floatGuests * constants.PerPersonSoftDrinksFee
	}
	if form.WillProvideCups != nil && !*form.WillProvideCups {
		totalCost += floatGuests * constants.PerPersonCupsFee
	}
	if form.WillProvideIce != nil && !*form.WillProvideIce {
		totalCost += floatGuests * constants.PerPersonIceFee
	}
	if form.WillRequireGlassware != nil && *form.WillRequireGlassware {
		totalCost += floatGuests * constants.PerPersonGlasswareFee
	}

	return totalCost
}
