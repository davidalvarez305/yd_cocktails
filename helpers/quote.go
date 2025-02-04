package helpers

import (
	"github.com/davidalvarez305/yd_cocktails/constants"
	"github.com/davidalvarez305/yd_cocktails/types"
)

func calculateBartenderRate(numBartenders, hours float64) float64 {
	return constants.BartendingRate * numBartenders * hours
}

func CalculatePackageQuote(form types.LeadQuoteForm) float64 {
	var totalCost float64

	guests := SafeInt(form.Guests)
	if guests <= 0 {
		return 0.00
	}

	floatGuests := float64(guests)

	hours := SafeInt(form.Hours)

	bartenderRate := calculateBartenderRate(float64(SafeInt(form.NumberOfBartenders)), float64(hours))

	totalCost += bartenderRate

	weWillProvideAlcohol := SafeBoolDefaultFalse(form.WeWillProvideAlcohol)
	weWillProvideBeer := SafeBoolDefaultFalse(form.WeWillProvideBeer)
	weWillProvideWine := SafeBoolDefaultFalse(form.WeWillProvideWine)
	weWillProvideMixers := SafeBoolDefaultFalse(form.WeWillProvideMixers)
	weWillProvideJuices := SafeBoolDefaultFalse(form.WeWillProvideJuice)
	weWillProvideCupsStrawsNapkins := SafeBoolDefaultFalse(form.WeWillProvideCupsStrawsNapkins)
	weWillProvideSoftDrinks := SafeBoolDefaultFalse(form.WeWillProvideSoftDrinks)
	weWillProvideIce := SafeBoolDefaultFalse(form.WeWillProvideIce)

	willRequireBar := SafeBoolDefaultFalse(form.WillRequireBar)
	willRequireGlassware := SafeBoolDefaultFalse(form.WillRequireGlassware)

	if weWillProvideAlcohol {
		totalCost += floatGuests * constants.PerPersonAlcoholFee
	}

	if weWillProvideBeer {
		totalCost += floatGuests * constants.PerPersonBeerFee
	}

	if weWillProvideWine {
		totalCost += floatGuests * constants.PerPersonWineFee
	}

	if weWillProvideMixers {
		totalCost += floatGuests * constants.PerPersonMixersFee
	}
	if weWillProvideJuices {
		totalCost += floatGuests * constants.PerPersonJuicesFee
	}
	if weWillProvideSoftDrinks {
		totalCost += floatGuests * constants.PerPersonSoftDrinksFee
	}
	if weWillProvideCupsStrawsNapkins {
		totalCost += floatGuests * constants.PerPersonCupsStrawsNapkinsFee
	}
	if weWillProvideIce {
		totalCost += floatGuests * constants.PerPersonIceFee
	}

	if willRequireGlassware {
		totalCost += floatGuests * constants.PerPersonGlasswareFee
	}

	if willRequireBar {
		numBarsFloat := float64(SafeInt(form.NumBars))
		totalCost += numBarsFloat * constants.BarRentalCost
	}

	return totalCost
}
