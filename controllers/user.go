package controllers

import (
	"net/http"
	"strings"

	"gopkg.in/mgo.v2/bson"

	"github.com/moehlone/mongodm_sample/models"
)

type UserController struct {
	baseController
}

/**
 * @api {post} /user Create User
 *
 * @apiName Create
 * @apiGroup User
 *
 * @apiVersion 0.0.1
 *
 * @apiDescription Create a new user with the specified, required fields.
 *
 * @apiUse ErrorInternal
 *
 *
 * @apiParam {Object} user The user object with all required fields
 * @apiParam {String} user.email Users the email (used for login) <code>validation: email</code>
 * @apiParam {String} user.firstname Users firstname <code>minLen: 2 maxLen: 30</code>
 * @apiParam {String} user.lastname Users lastname <code>minLen: 2 maxLen: 30</code>
 * @apiParam {String} [user.username] Users nickname <code>minLen: 2 maxLen: 15</code>
 * @apiParam {Address} [address] User address object
 * @apiParam {String} [address.street] Street and Number
 * @apiParam {String} [address.city] City
 * @apiParam {Number} [address.zip] Zipcode
 * @apiParam {String} password The user password as plain text <code>minLen: 8 maxLen: 50</code>
 *
 *
 * @apiSuccess {Object} user The user object which was created.
 *
 * @apiSuccessExample  {json} Success-Response
 *    HTTP/1.1 201 Created
 *    {
 *	   "user": {
 *	     "id": "5710ed60e5feae376f000001",
 *	     "createdAt": "2016-04-15T15:32:16.670823736+02:00",
 *	     "updatedAt": "2016-04-15T15:32:16.670823736+02:00",
 *	     "firstname": "Max",
 *	     "lastname": "Mustermann",
 *	     "username": "",
 *	     "email": "test@test.de",
 *	     "address": null
 *	   }
 *	 }
 *
 * @apiError MissingInvalidFields Some of the required fields are missing or invalid
 *
 * @apiErrorExample {json} Missing or invalid fields
 *HTTP/1.1 400 Bad Request
 *{
 *  "error": {
 *    "code": 400,
 *    "message": "Bad Request",
 *    "userInfo": [
 *      {
 *        "message": "Field 'firstname' is required."
 *      },
 *      {
 *        "message": "Field 'lastname' is required."
 *      },
 *      {
 *        "message": "Field 'email' is required."
 *      },
 *      {
 *        "message": "Field 'password' is required."
 *      }
 *    ]
 *  }
 *}
 */
func (self *UserController) Create() {

	// See https://godoc.org/github.com/zebresel-com/mongodm#Connection.Model
	User := self.db.Model("User")
	user := &models.User{}

	// See https://godoc.org/github.com/zebresel-com/mongodm#Model.New
	err, requestMap := User.New(user, self.Ctx.Input.RequestBody)

	if err != nil {
		self.response.Error(http.StatusBadRequest, err)
		return
	}

	/**
	 * Use a custom validation for the field password.
	 * The default validation will be called automatically within this step.
	 */
	if valid, issues := user.Validate(requestMap["password"]); valid {

		// See https://godoc.org/github.com/zebresel-com/mongodm#DocumentBase.Save
		err = user.Save()

		if err != nil {

			self.response.Error(http.StatusInternalServerError)
			return
		}

	} else {

		self.response.Error(http.StatusBadRequest, issues)
		return
	}

	self.response.AddContent("user", user)
	self.response.SetStatus(http.StatusCreated)
	self.response.ServeJSON()
}

/**
 * @api {get} /user Get All
 *
 * @apiName Index
 * @apiGroup User
 *
 * @apiVersion 0.0.1
 *
 * @apiDescription Fetches all users
 *
 * @apiUse ErrorInternal
 *
 * @apiUse Pagination
 *
 * @apiSuccess {Array(User)} users The array of all users
 * @apiSuccess {String} users.id User ID
 * @apiSuccess {String} users.createdAt Timestamp of first creation
 * @apiSuccess {String} users.updatedAt Timestamp of last update
 * @apiSuccess {String} users.firstname Firstname
 * @apiSuccess {String} users.lastname Lastname
 * @apiSuccess {String} users.username Username
 * @apiSuccess {String} users.email Email
 * @apiSuccess {Address} address User address object
 * @apiSuccess {String} address.street Street and Number
 * @apiSuccess {String} address.city City
 * @apiSuccess {Number} address.zip Zipcode
 *
 * @apiSuccessExample  {json} Success-Response
 *    HTTP/1.1 201 Created
 *    {
 *          "pagination": {
 *              "first": "/api/v1/users?skip=0&limit=10",
 *              "last": "/api/v1/users?skip=0&limit=2",
 *              "previous": "",
 *              "next": "",
 *              "totalRecords": 2,
 *              "limit": 10,
 *              "pages": 1,
 *              "currentPage": 1
 *          },
 *          "users": [
 *              {
 *                  "id": "55e961ab113c61d6d3000001",
 *                  "createdAt": "2015-09-04T11:17:31.219+02:00",
 *                  "updatedAt": "2015-09-08T16:35:44.484+02:00",
 *                  "firstname": "Max",
 *                  "lastname": "Mustermann",
 *                  "username": "zebresel_admin",
 *                  "email": "admin@zebresel.com",
 *                  "address": {
 *				        "street": "Musterstraße 12",
 *				        "city": "Erfurt",
 *				        "zip": 99085
 *				    }
 *              },
 *              {
 *                  "id": "55e9621d113c61d6d3000002",
 *                  "createdAt": "2015-09-04T11:19:25.83+02:00",
 *                  "updatedAt": "2015-09-04T11:19:25.83+02:00",
 *                  "firstname": "Max",
 *                  "lastname": "Musterfrau",
 *                  "username": "moehlone",
 *                  "email": "hello@zebresel.com",
 *                  "address": {
 *				        "street": "Musterstraße 12",
 *				        "city": "Erfurt",
 *				        "zip": 99085
 *				    }
 *              }
 *          ]
 *      }
 *
 */
func (self *UserController) GetAll() {

	User := self.db.Model("user")
	users := []*models.User{}

	searchTag := self.GetString("search")
	fieldName := self.GetString("field")

	query := bson.M{"deleted": false}

	if searchTag != "" {

		if fieldName != "" {

			query = bson.M{"deleted": false, fieldName: searchTag}

		} else {

			splitted := strings.Split(searchTag, " ")
			regexArray := []interface{}{}

			for _, value := range splitted {

				if len(value) > 0 {
					regexArray = append(regexArray, &bson.RegEx{Pattern: value, Options: "i"})
				}
			}

			query = bson.M{
				"deleted": false,
				"$or": []interface{}{
					bson.M{"firstname": bson.M{"$in": regexArray}},
					bson.M{"lastname": bson.M{"$in": regexArray}},
					bson.M{"email": bson.M{"$in": regexArray}},
					bson.M{"username": bson.M{"$in": regexArray}},
				}}
		}
	}

	queryCount, queryErr := User.Find(query).Count()

	if queryErr != nil {

		self.response.Error(http.StatusInternalServerError)
		return
	}

	err := User.Find(query).Sort("-createdAt").Skip(self.paging.skip).Limit(self.paging.limit).Exec(&users)

	if len(users) > 0 && err != nil {
		self.response.Error(http.StatusInternalServerError)
		return
	}

	self.response.CreatePaging(self.paging.skip, self.paging.limit, queryCount, len(users))
	self.response.AddContent("users", users)

	self.response.ServeJSON()
}
