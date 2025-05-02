package dup

import (
	"fmt"
	"github.com/jinzhu/copier"
	"math/rand"
	"testing"
)

type User struct {
	Name string
	Age  int
}

type ReportUser struct {
	Name string
	Age  int
}

func genUsers() []*User {
	var users []*User
	for i := 0; i < 10; i++ {
		users = append(users, &User{
			Name: fmt.Sprintf("U%d", i+1),
			Age:  rand.Intn(100),
		})
	}
	return users
}

func copyUserToReport(src interface{}, dest interface{}) {
	// Copy(src, dest)
}

func copyUsersToReport(src []*User) interface{} {
	var users = make([]*ReportUser, len(src))
	for i := 0; i < len(src); i++ {
		var rUser ReportUser
		copyUserToReport(src[i], &rUser)

		users[i] = &rUser
	}

	return users

}

func TestCopyOneUser(t *testing.T) {
	users := genUsers()
	// user := users[0]

	//var rUser ReportUser

	var rUsers []*ReportUser

	copyUserToReport(users, &rUsers)

	copyUsersToReport(users)

	t.Log(rUsers)
	t.Log(users[0])
}

type MenuSrc struct {
	Name string
	Code *string
	ID   uint64
	Age  uint64
}

type MenuDst struct {
	Name *string
	Code *string
	ID   uint64
	Age  *uint64
}

func copyD(src, dst interface{}) error {
	return nil
}

func TestCopy(t *testing.T) {
	var name = "google"
	var dst = &MenuDst{}

	var src = &MenuSrc{
		ID:   252271854679490561,
		Name: name,
	}

	err := copyD(src, dst)
	if err != nil {
		t.Error(err)
	}
}

func TestCopier(t *testing.T) {
	var name = "google"
	var dst = &MenuDst{}

	var src = &MenuSrc{
		ID:   252271854679490561,
		Name: name,
	}

	Copy(src, dst)
}

func TestCopierWithNil(t *testing.T) {
	src := &MenuSrc{
		ID: 0,
	}
	dst := &MenuDst{
		ID: 11,
	}

	CopyWithOption(src, dst, copier.Option{IgnoreEmpty: false})
	fmt.Printf("dst: %+v", dst)
}
