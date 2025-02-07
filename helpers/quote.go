package helpers

import (
	"github.com/davidalvarez305/yd_cocktails/constants"
	"github.com/davidalvarez305/yd_cocktails/models"
	"github.com/davidalvarez305/yd_cocktails/types"
)

func calculateBartenderRate(numBartenders, hours float64) float64 {
	return constants.BartendingRate * numBartenders * hours
}

func CalculatePackageQuote(form types.LeadQuoteForm, barTypes []models.BarType, alcoholFeeAdjustment float64) float64 {
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
	weWillProvideGarnish := SafeBoolDefaultFalse(form.WeWillProvideGarnish)

	willRequireBar := SafeBoolDefaultFalse(form.WillRequireBar)
	barTypeId := SafeInt(form.BarTypeID)

	willRequireGlassware := SafeBoolDefaultFalse(form.WillRequireGlassware)

	// Alcohol
	if weWillProvideAlcohol {
		totalCost += floatGuests * constants.PerPersonAlcoholFee * alcoholFeeAdjustment
	}
	if weWillProvideBeer {
		totalCost += floatGuests * constants.PerPersonBeerFee
	}

	if weWillProvideWine {
		totalCost += floatGuests * constants.PerPersonWineFee
	}
	// Alcohol

	// Ingredients
	if weWillProvideIce {
		totalCost += floatGuests * constants.PerPersonIceFee
	}
	if weWillProvideSoftDrinks {
		totalCost += floatGuests * constants.PerPersonSoftDrinksFee
	}
	if weWillProvideJuices {
		totalCost += floatGuests * constants.PerPersonJuicesFee
	}
	if weWillProvideMixers {
		totalCost += floatGuests * constants.PerPersonMixersFee
	}
	if weWillProvideGarnish {
		totalCost += floatGuests * constants.PerPersonGarnishFee
	}
	// Ingredients

	// Supplies
	if weWillProvideCupsStrawsNapkins {
		totalCost += floatGuests * constants.PerPersonCupsStrawsNapkinsFee
	}

	if willRequireGlassware {
		totalCost += floatGuests * constants.PerPersonGlasswareFee
	}
	// Supplies

	if willRequireBar {
		var barRentalFee float64
		for _, barType := range barTypes {
			if barType.BarTypeID == barTypeId {
				barRentalFee = barType.Price
			}
		}

		numBarsFloat := float64(SafeInt(form.NumBars))
		totalCost += numBarsFloat * barRentalFee
	}

	return totalCost
}
