package plexhooks

import (
	"encoding/json"
	"io/ioutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Plexhooks", func() {
	Context("a music play", func() {
		It("parses it properly", func() {
			body, err := ioutil.ReadFile("test-fixtures/music.json")
			Expect(err).ShouldNot(HaveOccurred())

			response, err := ParseWebhook(body)
			Expect(err).ShouldNot(HaveOccurred())

			responseString, err := json.Marshal(response)
			Expect(err).ShouldNot(HaveOccurred())

			newResponse, err := ParseWebhook(responseString)
			Expect(err).ShouldNot(HaveOccurred())

			Expect(response).To(Equal(newResponse))

			Expect(newResponse.Event).To(Equal("media.play"))
			Expect(newResponse.Account.Id).To(Equal(1))
			Expect(newResponse.Server.Title).To(Equal("Office"))
			Expect(newResponse.Player.Title).To(Equal("Plex Web (Safari)"))

			Expect(newResponse.Metadata.Type).To(Equal("track"))
			Expect(newResponse.Metadata.Title).To(Equal("Love The One You're With"))
		})
	})

	Context("a tv play", func() {
		It("parses it properly", func() {
			body, err := ioutil.ReadFile("test-fixtures/tv.json")
			Expect(err).ShouldNot(HaveOccurred())

			response, err := ParseWebhook(body)
			Expect(err).ShouldNot(HaveOccurred())

			responseString, err := json.Marshal(response)
			Expect(err).ShouldNot(HaveOccurred())

			newResponse, err := ParseWebhook(responseString)
			Expect(err).ShouldNot(HaveOccurred())

			Expect(response).To(Equal(newResponse))

			Expect(newResponse.User).To(Equal(true))
			Expect(newResponse.Account.Title).To(Equal("testyboi"))
			Expect(newResponse.Server.Title).To(Equal("nice"))
			Expect(newResponse.Player.PublicAddress).To(Equal("200.200.200.200"))

			Expect(newResponse.Metadata.Type).To(Equal("episode"))
			Expect(newResponse.Metadata.Title).To(Equal("A Clone of My Own"))

			Expect(newResponse.Metadata.Director[0].Tag).To(Equal("Rich Moore"))
			Expect(newResponse.Metadata.Writer[0].Filter).To(Equal("writer=49503"))
		})
	})

	Context("a movie play", func() {
		It("parses it properly", func() {
			body, err := ioutil.ReadFile("test-fixtures/movie.json")
			Expect(err).ShouldNot(HaveOccurred())

			response, err := ParseWebhook(body)
			Expect(err).ShouldNot(HaveOccurred())

			responseString, err := json.Marshal(response)
			Expect(err).ShouldNot(HaveOccurred())

			newResponse, err := ParseWebhook(responseString)
			Expect(err).ShouldNot(HaveOccurred())

			Expect(response).To(Equal(newResponse))

			Expect(newResponse.User).To(Equal(true))
			Expect(newResponse.Account.Title).To(Equal("testyboi"))
			Expect(newResponse.Server.Title).To(Equal("what"))
			Expect(newResponse.Player.PublicAddress).To(Equal("200.200.200.200"))

			Expect(newResponse.Metadata.Type).To(Equal("movie"))
			Expect(newResponse.Metadata.Studio).To(Equal("Hawk Films"))
			Expect(newResponse.Metadata.ContentRating).To(Equal("PG"))
			Expect(newResponse.Metadata.AudienceRating).To(Equal(float32(9.2)))

			Expect(newResponse.Metadata.Director[0].Tag).To(Equal("Stanley Kubrick"))
			Expect(newResponse.Metadata.Writer[0].Filter).To(Equal("writer=7"))
			Expect(newResponse.Metadata.Producer[0].Id).To(Equal(42))
			Expect(newResponse.Metadata.Country[0].Tag).To(Equal("United Kingdom"))

			Expect(len(newResponse.Metadata.Role)).To(Equal(43))
			Expect(newResponse.Metadata.Role[25].Tag).To(Equal("Anthony Herrick"))

			Expect(len(newResponse.Metadata.Similar)).To(Equal(20))
			Expect(newResponse.Metadata.Similar[8].Tag).To(Equal("Touch of Evil"))

			Expect(len(newResponse.Metadata.Genre)).To(Equal(4))
			Expect(newResponse.Metadata.Genre[3].Tag).To(Equal("War"))
		})
	})
})
