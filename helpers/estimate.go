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

	willProvideLiquor := SafeBoolDefaultFalse(form.WillProvideLiquor)
	willProvideBeerAndWine := SafeBoolDefaultFalse(form.WillProvideBeerAndWine)
	willProvideMixers := SafeBoolDefaultFalse(form.WillProvideMixers)
	willProvideJuices := SafeBoolDefaultFalse(form.WillProvideJuices)
	willProvideCups := SafeBoolDefaultFalse(form.WillProvideCups)
	willProvideSoftDrinks := SafeBoolDefaultFalse(form.WillProvideSoftDrinks)
	willRequireBar := SafeBoolDefaultFalse(form.WillRequireBar)
	willRequireGlassware := SafeBoolDefaultFalse(form.WillRequireGlassware)
	willProvideIce := SafeBoolDefaultFalse(form.WillProvideIce)

	if willProvideLiquor {
		totalCost += floatGuests * constants.PerPersonAlcoholFee
	}
	if willProvideBeerAndWine {
		totalCost += floatGuests * constants.PerPersonBeerAndWineFee
	}
	if packageTypeID == constants.FullOpenBarPackageTypeID {
		totalCost += floatGuests * constants.FullOpenBarFee
	}
	if packageTypeID == constants.PartialOpenBarPackageTypeID {
		totalCost += floatGuests * constants.PartialOpenBarFee
	}
	if willProvideMixers {
		totalCost += floatGuests * constants.PerPersonMixersFee
	}
	if willProvideJuices {
		totalCost += floatGuests * constants.PerPersonJuicesFee
	}
	if willProvideSoftDrinks {
		totalCost += floatGuests * constants.PerPersonSoftDrinksFee
	}
	if willProvideCups {
		totalCost += floatGuests * constants.PerPersonCupsFee
	}
	if willProvideIce {
		totalCost += floatGuests * constants.PerPersonIceFee
	}
	if willRequireGlassware {
		totalCost += floatGuests * constants.PerPersonGlasswareFee
	}

	if willRequireBar {
		numBarsFloat := float64(SafeInt(form.NumBars))
		totalCost += numBarsFloat * constants.BarRentalCost

		totalCost += numBarsFloat * constants.BarSetupAndBreakdownFee
	}

	return totalCost
}
