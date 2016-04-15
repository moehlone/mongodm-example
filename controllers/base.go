package controllers

import (
	"github.com/astaxie/beego"
	"github.com/moehlone/mongodm_sample/models"
	response "github.com/zebresel-com/beego-response"
	"github.com/zebresel-com/mongodm"
)

const DEFAULT_LIMIT int = 10
const DEFAULT_SKIP int = 0

type paging struct {
	skip  int
	limit int
}

type baseController struct {
	beego.Controller

	db       *mongodm.Connection
	user     *models.User
	response *response.Response
	paging   *paging
}

var (
	Database *mongodm.Connection
)

func (self *baseController) Prepare() {

	self.db = Database
	self.response = response.New(self.Ctx)

	self.paging = &paging{}

	limit, limitErr := self.GetInt("limit")
	skip, skipErr := self.GetInt("skip")

	if limitErr != nil || limit < 0 {
		limit = DEFAULT_LIMIT
	}

	if skipErr != nil || skip < 0 {
		skip = DEFAULT_SKIP
	}

	self.paging.limit = limit
	self.paging.skip = skip
}

/**
 * @apiDefine ErrorInternal
 *
 * @apiError InternalServerError Some server problems occured.
 *
 * @apiErrorExample {json} Internal Server Error
 *     HTTP/1.1 500 Internal Server Error
 *     {
 *          "error": {
 *              "code": 500,
 *              "message": "Internal Server Error"
 *          }
 *      }
 *
 */

/**
 * @apiDefine Pagination
 *
 * @apiParam {Number} [skip=0] The offset to start with searching
 * @apiParam {Number} [limit=10] Maximum records for one page
 *
 *
 * @apiSuccess {Object} pagination The pagination object
 * @apiSuccess {String} pagination.first Path to first page including other params
 * @apiSuccess {String} pagination.last Path to last page including other params
 * @apiSuccess {String} pagination.previous Path to previous page
 * @apiSuccess {String} pagination.next Path to next page
 * @apiSuccess {Number} pagination.totalRecords All records on all pages
 * @apiSuccess {Number} pagination.limit Records on the current page
 * @apiSuccess {Number} pagination.pages Amount of all pages
 * @apiSuccess {Number} pagination.currentPage Number of the current page
 *
 */
