// Code generated by pegomock. DO NOT EDIT.
package matchers

import (
	versioned2 "github.com/jenkins-x/jx/pkg/client/clientset/versioned"
	"github.com/petergtz/pegomock"
	"reflect"
)

func AnyVersioned2Interface() versioned2.Interface {
	pegomock.RegisterMatcher(pegomock.NewAnyMatcher(reflect.TypeOf((*(versioned2.Interface))(nil)).Elem()))
	var nullValue versioned2.Interface
	return nullValue
}

func EqVersioned2Interface(value versioned2.Interface) versioned2.Interface {
	pegomock.RegisterMatcher(&pegomock.EqMatcher{Value: value})
	var nullValue versioned2.Interface
	return nullValue
}
