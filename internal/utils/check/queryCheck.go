package check

import (
	"fmt"
	"regexp"
	"strconv"
	"sync"

	"github.com/v.kirpichov/db_tp/internal/models"
	"github.com/v.kirpichov/db_tp/internal/utils/constants"
)

const slugReg = "^(\\d|\\w|-|_)*(\\w|-|_)(\\d|\\w|-|_)*$"

type queryCheck struct {
	slugRegExCompiled *regexp.Regexp
}

var instanceLock = &sync.Mutex{}
var instance *queryCheck

func GetInstance() (*queryCheck, error) {
	if instance == nil {
		instanceLock.Lock()
		defer instanceLock.Unlock()
		if instance == nil {
			var err error
			instance, err = create()
			if err != nil {
				return nil, err
			}
		}
	}
	return instance, nil
}

func create() (checker *queryCheck, err error) {
	checker = &queryCheck{}
	checker.slugRegExCompiled, err = regexp.Compile(slugReg)
	if err != nil {
		return nil, err
	}
	return
}

func (checker *queryCheck) CheckSlug(slug string) bool {
	return checker.slugRegExCompiled.MatchString(slug)
}

func (checker *queryCheck) CheckForumQuery(query *models.ForumQueryParams) {
	if query.Limit == 0 {
		query.Limit = 100
	}
}

func (checker *queryCheck) CheckForumUserQuery(query *models.ForumUserQueryParams) {
	if query.Limit == 0 {
		query.Limit = 100
	}
}

func (checker *queryCheck) GetSlugOrIdOrErr(slugOrId string) (slug string, id int, err error) {
	if slugOrId == "" {
		err = fmt.Errorf("пустой slug or id")
		return
	}

	id, err = strconv.Atoi(slugOrId)
	if err == nil {
		return
	}

	if checker.CheckSlug(slugOrId) == false {
		err = fmt.Errorf("неверный slug or id")
		return
	}

	err = nil
	slug = slugOrId
	return
}

func (checker *queryCheck) CheckPostsQuery(query *models.PostsQueryParams) bool {
	if query.Limit == 0 {
		query.Limit = 100
	}

	if query.SortType == "" {
		query.SortType = constants.SortFlat
	}

	if query.SortType != constants.SortParentTree &&
		query.SortType != constants.SortFlat &&
		query.SortType != constants.SortTree {
		return false
	}
	return true
}

func (checker *queryCheck) CheckVote(vote *models.Vote) bool {
	if vote.Voice != 1 && vote.Voice != -1 {
		return false
	}
	return true
}
