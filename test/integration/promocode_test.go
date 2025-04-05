package control_test

import (
	"APIs/internal/common/entities"
	"APIs/internal/common/models"
	"APIs/internal/services/promocode/ports"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Promocode IT", func() {

	Context("POST /promocodes", func() {
		It("OK - should insert a new item into the database", func() {
			// Prepare the HTTP POST request
			reqBody := `{
				"name": "WeatherCode",
				"advantage": { "percent": 20 },
				"restrictions": [
					{
						"date": {
							"after": "2024-01-01",
							"before": "2026-06-30"
						}
					},
					{
						"or": [
						{
							"age": {
								"eq": 40
							}
						},
						{
							"and": [
							{
								"age": {
									"lt": 30,
									"gt": 15
								}
							},
							{
								"weather": {
									"is": "clear",
									"temp": {
										"gt": 15
									}
								}
							}
							]
						}
						]
					}
				]
			}`
			req := httptest.NewRequest(http.MethodPost, "/v1/promocodes", strings.NewReader(reqBody))
			rec := httptest.NewRecorder()

			// Serve the request using the router
			router.ServeHTTP(rec, req)

			// Verify the HTTP response code
			Expect(rec.Code).To(Equal(http.StatusCreated))

			var body ports.Promocode
			err := json.NewDecoder(rec.Body).Decode(&body)
			Expect(err).ToNot(HaveOccurred())
			Expect(body.Name).To(Equal("WeatherCode"))

			// Verify the item was inserted into the database
			model := &models.Promocode{}
			err = db.First(model, body.Id).Error
			Expect(err).ToNot(HaveOccurred())
			Expect(model.Name).To(Equal("WeatherCode"))
		})
	})

	Context("POST /promocodes/_valdiate", func() {
		It("OK - should validate promocode", func() {
			// Prepare the HTTP POST request
			reqBody := `{
				"promocode_name": "WeatherCode",
				"arguments": {
					"age": 30,
					"town": "Lyon"
				}
			}`
			req := httptest.NewRequest(http.MethodPost, "/v1/promocodes/_validate", strings.NewReader(reqBody))
			rec := httptest.NewRecorder()

			// Serve the request using the router
			router.ServeHTTP(rec, req)

			// Verify the HTTP response code
			Expect(rec.Code).To(Equal(http.StatusOK))

			var body ports.PromocodeValidationResponse
			err := json.NewDecoder(rec.Body).Decode(&body)
			Expect(err).ToNot(HaveOccurred())
			Expect(body.PromocodeName).To(Equal("WeatherCode"))
			Expect(body.Status).To(Equal(ports.Accepted))

			// Verify the item was inserted into the database
			model := &models.Weather{}
			err = db.Where("town = ?", "Lyon").First(model).Error
			Expect(err).ToNot(HaveOccurred())
			Expect(model.Type).To(Equal(entities.Clear))
		})
	})
})
