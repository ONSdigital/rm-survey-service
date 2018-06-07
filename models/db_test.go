package models

import (
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	. "github.com/smartystreets/goconvey/convey"
)

func TestDBConnectionTested(t *testing.T) {
	Convey("db connection tested before returned", t, func() {
		db, _, err := sqlmock.New()
		defer db.Close()

		err = testDBConnection(db)

		So(err, ShouldBeNil)
	})
}
