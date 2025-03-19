package controllers

import (
	"log"
	"net/http"
	"time"

	"github.com/arjunsaxaena/Moonrider-Assignment/model"
	"github.com/arjunsaxaena/Moonrider-Assignment/repository"
	"github.com/gin-gonic/gin"
)

type IdentifyRequest struct {
	Email       *string `json:"email"`
	PhoneNumber *string `json:"phoneNumber"`
}

type IdentifyResponse struct {
	PrimaryContactID    string   `json:"primaryContactId"`
	Emails              []string `json:"emails"`
	PhoneNumbers        []string `json:"phoneNumbers"`
	SecondaryContactIDs []string `json:"secondaryContactIds"`
}

func IdentifyContact(c *gin.Context) {
	var req IdentifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("Invalid request payload:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if req.Email == nil && req.PhoneNumber == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Either email or phoneNumber must be provided"})
		return
	}

	if (req.Email != nil && *req.Email == "") && (req.PhoneNumber != nil && *req.PhoneNumber == "") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: Email or PhoneNumber required"})
		return
	}

	var email, phoneNumber string
	if req.Email != nil {
		email = *req.Email
	}
	if req.PhoneNumber != nil {
		phoneNumber = *req.PhoneNumber
	}

	contacts, err := repository.Get(repository.DB, email, phoneNumber)
	if err != nil {
		log.Println("Error fetching contacts:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	if len(contacts) == 0 {
		newContact := &model.Contact{
			Email:          req.Email,
			PhoneNumber:    req.PhoneNumber,
			LinkPrecedence: "primary",
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}
		if err := repository.Create(repository.DB, newContact); err != nil {
			log.Println("Error creating contact:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		response := IdentifyResponse{
			PrimaryContactID:    newContact.ID,
			Emails:              []string{*newContact.Email},
			PhoneNumbers:        []string{*newContact.PhoneNumber},
			SecondaryContactIDs: []string{},
		}
		c.JSON(http.StatusOK, response)
		return
	}

	primaryContact := contacts[0]
	emails := map[string]bool{}
	phoneNumbers := map[string]bool{}
	secondaryContactIDs := []string{}

	// Add primaryContact's email and phone number to the maps
	if primaryContact.Email != nil {
		emails[*primaryContact.Email] = true
	}
	if primaryContact.PhoneNumber != nil {
		phoneNumbers[*primaryContact.PhoneNumber] = true
	}

	for _, contact := range contacts {
		if contact.LinkPrecedence == "primary" {
			primaryContact = contact
		} else {
			secondaryContactIDs = append(secondaryContactIDs, contact.ID)
		}

		if contact.Email != nil {
			emails[*contact.Email] = true
		}
		if contact.PhoneNumber != nil {
			phoneNumbers[*contact.PhoneNumber] = true
		}
	}

	if req.PhoneNumber != nil && !phoneNumbers[*req.PhoneNumber] {
		phoneNumbers[*req.PhoneNumber] = true
	}

	if (req.Email != nil && !emails[*req.Email]) || (req.PhoneNumber != nil && !phoneNumbers[*req.PhoneNumber]) {
		newSecondaryContact := &model.Contact{
			Email:          req.Email,
			PhoneNumber:    req.PhoneNumber,
			LinkedID:       &primaryContact.ID,
			LinkPrecedence: "secondary",
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}
		if err := repository.Create(repository.DB, newSecondaryContact); err != nil {
			log.Println("Error creating secondary contact:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
		secondaryContactIDs = append(secondaryContactIDs, newSecondaryContact.ID)

		if newSecondaryContact.Email != nil {
			emails[*newSecondaryContact.Email] = true
		}
		if newSecondaryContact.PhoneNumber != nil {
			phoneNumbers[*newSecondaryContact.PhoneNumber] = true
		}
	}

	emailList := []string{}
	for email := range emails {
		emailList = append(emailList, email)
	}

	phoneNumberList := []string{}
	for phoneNumber := range phoneNumbers {
		phoneNumberList = append(phoneNumberList, phoneNumber)
	}

	response := IdentifyResponse{
		PrimaryContactID:    primaryContact.ID,
		Emails:              emailList,
		PhoneNumbers:        phoneNumberList,
		SecondaryContactIDs: secondaryContactIDs,
	}
	c.JSON(http.StatusOK, response)
}
