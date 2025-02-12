package helpers

import "github.com/davidalvarez305/yd_cocktails/models"

func CalculatePackageQuote(services []models.QuoteService) float64 {
	var total float64
	for _, service := range services {
		total += float64(service.Units) * service.PricePerUnit
	}

	return total
}
