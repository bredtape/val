package val

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNil(t *testing.T) {
	Convey("Given nil *Errors", t, func() {
		var e *Errors

		Convey("should be nil", func() {
			So(e, ShouldBeNil)
		})
	})
}

func TestFields(t *testing.T) {
	Convey("Given a field err", t, func() {
		e1 := New("name1", "desc1")

		Convey("should not be nil", func() {
			So(e1, ShouldNotBeNil)

			Convey(".Error() should equal", func() {
				So(e1.Error(), ShouldEqual, "validation error: 'name1' desc1")
			})
		})

		Convey("Concat nothing", func() {
			e := Concat()

			Convey("should be nil", func() {
				So(e, ShouldBeNil)
			})
		})

		Convey("Concat nil", func() {
			e := Concat(nil)

			Convey("should be nil", func() {
				So(e, ShouldBeNil)
			})
		})

		Convey("Concat with other", func() {
			e2 := New("name2", "desc2")

			e := Concat(e1, e2)

			Convey("should not be nil", func() {
				So(e, ShouldNotBeNil)

				Convey(".Error() should mention both", func() {
					So(e.Error(), ShouldEqual, "validation errors: 'name1' desc1, 'name2' desc2")
				})
			})

			Convey("should not have changed the inputs", func() {
				So(e1, ShouldResemble, New("name1", "desc1"))
				So(e2, ShouldResemble, New("name2", "desc2"))
			})
		})
	})
}

func TestReflection(t *testing.T) {
	Convey("Initialize as 'error'", t, func() {
		var e error = New("name1", "desc1")

		Convey("converting back to *Errors", func() {
			e1, ok := e.(*Errors)

			Convey("should be possible", func() {
				So(ok, ShouldBeTrue)

				Convey("should have same fields", func() {
					So(e1.GetFields(), ShouldResemble, []Field{{Name: "name1", Description: "desc1"}})
				})
			})
		})
	})
}

func TestChildren(t *testing.T) {
	Convey("constructing a hierachy", t, func() {
		child := New("child1", "y")

		Convey("parent with 1 child", func() {
			e := NewChildren("parent", child)

			Convey("should have .Error()", func() {
				So(e.Error(), ShouldEqual, "validation error: 'parent.child1' y")
			})
		})

		Convey("parent with child and no field name", func() {
			e := NewChildren("parent", New("", "x"))
			Convey("should have .Error()", func() {
				So(e.Error(), ShouldEqual, "validation error: 'parent' x")
			})
		})
	})
}
