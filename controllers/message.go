package controllers

import (
	"net/http"

	"gopkg.in/mgo.v2/bson"

	"github.com/moehlone/mongodm_sample/models"
)

type MessageController struct {
	baseController
}

/**
 * @api {get} /message Get All
 *
 * @apiName Index
 * @apiGroup Message
 *
 * @apiVersion 0.0.1
 *
 * @apiDescription Fetches all messages
 *
 * @apiUse ErrorInternal
 *
 * @apiUse Pagination
 *
 * @apiSuccess {Array(Message)} messages The array of all messages
 * @apiSuccess {String} message.id User ID
 * @apiSuccess {String} message.createdAt Timestamp of first creation
 * @apiSuccess {String} message.updatedAt Timestamp of last update
 * @apiSuccess {User} message.sender User object of sender
 * @apiSuccess {Array(User)} message.receiver List of user objects (receiver)
 * @apiSuccess {String} message.text Message content
 *
 * @apiSuccessExample  {json} Success-Response
 *    HTTP/1.1 201 Created
 *    {
  "messages": [
    {
      "id": "564ca693e5feae2532000002",
      "createdAt": "2015-11-18T17:25:55.052+01:00",
      "updatedAt": "2015-11-18T17:25:55.052+01:00",
      "sender": {
        "id": "564c87e1e5feae0e3a000001",
        "createdAt": "2015-11-18T15:14:57.081+01:00",
        "updatedAt": "2015-11-18T15:14:57.081+01:00",
        "firstname": "Bob",
        "lastname": "Marley",
        "username": "",
        "email": "bob@marley.de",
        "address": {
          "street": "Liebknechtstraße 55a",
          "city": "Erfurt",
          "zip": 99085
        }
      },
      "receiver": [
        {
          "id": "564b3879e5feaed231000002",
          "createdAt": "2015-11-17T15:23:53.439+01:00",
          "updatedAt": "2015-11-17T15:23:53.439+01:00",
          "firstname": "Bob",
          "lastname": "Marley",
          "username": "",
          "email": "bob@marley.de",
          "address": {
            "street": "Liebknechtstraße 55a",
            "city": "Erfurt",
            "zip": 99085
          }
        },
        {
          "id": "564c87e1e5feae0e3a000001",
          "createdAt": "2015-11-18T15:14:57.081+01:00",
          "updatedAt": "2015-11-18T15:14:57.081+01:00",
          "firstname": "Bob",
          "lastname": "Marley",
          "username": "",
          "email": "bob@marley.de",
          "address": {
            "street": "Liebknechtstraße 55a",
            "city": "Erfurt",
            "zip": 99085
          }
        }
      ],
      "text": "Hello!"
    }
  ],
  "pagination": {
    "first": "/api/message?limit=1&skip=0",
    "last": "/api/message?limit=1&skip=19",
    "previous": "",
    "next": "/api/message?limit=1&skip=1",
    "totalRecords": 20,
    "limit": 1,
    "pages": 20,
    "currentPage": 1
  }
}
 *
*/
func (self *MessageController) GetAll() {

	Message := self.db.Model("message")
	messages := []*models.Message{}

	query := bson.M{"deleted": false}

	queryCount, queryErr := Message.Find(query).Count()

	if queryErr != nil {

		self.response.Error(http.StatusInternalServerError)
		return
	}

	err := Message.Find(query).Sort("sender").Skip(self.paging.skip).Limit(self.paging.limit).Populate("Sender").Exec(&messages)

	if err != nil {
		self.response.Error(http.StatusInternalServerError)
		return
	}

	for _, message := range messages {

		err := message.Populate("Receiver")

		if err != nil {
			self.response.Error(http.StatusInternalServerError)
			return
		}
	}

	self.response.CreatePaging(self.paging.skip, self.paging.limit, queryCount, len(messages))
	self.response.AddContent("messages", messages)
	self.response.ServeJSON()
}

/**
 * @api {post} /message Create Message
 *
 * @apiName Create
 * @apiGroup Message
 *
 * @apiVersion 0.0.1
 *
 * @apiDescription Create a new message with the specified, required fields.
 *
 * @apiUse ErrorInternal
 *
 *
 * @apiParam {Object} message The message object with all required fields
 * @apiParam {String(ObjectID)} message.sender The objectID of the sender (user)
 * @apiParam {Array(String(ObjectID))} message.receiver List containing objectIDs for all receivers (users)
 * @apiParam {String} message.text Message content
 *
 *
 * @apiSuccess {Object} message The message object which was created (unpopulated).
 *
 * @apiSuccessExample  {json} Success-Response
 *    HTTP/1.1 201 Created
 *    {
  "message": {
    "id": "571101e7e5feae48d8000001",
    "createdAt": "2016-04-15T16:59:51.97721049+02:00",
    "updatedAt": "2016-04-15T16:59:51.97721049+02:00",
    "sender": "564ca693e5feae2532000002",
    "receiver": [
      "564ca693e5feae2532000002"
    ],
    "text": "Some message.."
  }
}
 *
 * @apiError NotWrapped The message attributes were not wrapped in a message object (key)
 *
 * @apiErrorExample {json} Object not wrapped in typename
{
  "error": {
    "code": 400,
    "message": "Bad Request",
    "userInfo": [
      {
        "message": "object not wrapped in typename"
      }
    ]
  }
}
 *
 * @apiError MissingInvalidFields Some of the required fields are missing or invalid
 *
 * @apiErrorExample {json} Missing or invalid fields
 *HTTP/1.1 400 Bad Request
{
  "error": {
    "code": 400,
    "message": "Bad Request",
    "userInfo": [
      {
        "message": "Field 'sender' is required."
      },
      {
        "message": "Field 'receiver' is required."
      },
      {
        "message": "Field 'text' is required."
      }
    ]
  }
}
*/
func (self *MessageController) Create() {

	Message := self.db.Model("Message")
	message := &models.Message{}

	err, _ := Message.New(message, self.Ctx.Input.RequestBody)

	if err != nil {
		self.response.Error(http.StatusBadRequest, err)
		return
	}

	if valid, issues := message.Validate(); valid {

		err = message.Save()

		if err != nil {

			self.response.Error(http.StatusInternalServerError)
			return
		}

	} else {

		self.response.Error(http.StatusBadRequest, issues)
		return
	}

	self.response.AddContent("message", message)
	self.response.ServeJSON()
}
