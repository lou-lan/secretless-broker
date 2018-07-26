package main

import (
	"encoding/json"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/conjurinc/secretless/internal/app/secretless/providers"
	plugin_v1 "github.com/conjurinc/secretless/pkg/secretless/plugin/v1"

	_ "github.com/joho/godotenv/autoload"
)

// TestConjur_Provider tests the ability of the ConjurProvider to provide a Conjur accessToken
// as well as secret values.
func TestConjur_Provider(t *testing.T) {
	var err error

	name := "conjur"

	options := plugin_v1.ProviderOptions{
		Name: name,
	}

	provider := providers.ProviderFactories[name](options)

	Convey("Has the expected provider name", t, func() {
		So(provider.GetName(), ShouldEqual, "conjur")
	})

	Convey("Can provide an access token", t, func() {
		value, err := provider.GetValue("accessToken")
		So(err, ShouldBeNil)

		token := make(map[string]string)
		err = json.Unmarshal(value, &token)
		So(err, ShouldBeNil)
		So(token["protected"], ShouldNotBeNil)
		So(token["payload"], ShouldNotBeNil)
	})

	Convey("Can provide a secret to a fully qualified variable", t, func() {
		value, err := provider.GetValue("dev:variable:db/password")
		So(err, ShouldBeNil)

		So(string(value), ShouldEqual, "secret")
	})

	Convey("Can provide the default Conjur account name", t, func() {
		value, err := provider.GetValue("variable:db/password")
		So(err, ShouldBeNil)

		So(string(value), ShouldEqual, "secret")
	})

	Convey("Can provide the default Conjur account name and resource type", t, func() {
		value, err := provider.GetValue("db/password")
		So(err, ShouldBeNil)

		So(string(value), ShouldEqual, "secret")
	})

	Convey("Cannot provide an unknown value", t, func() {
		_, err = provider.GetValue("foobar")
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldEqual, "Variable 'foobar' not found in account 'dev'")
	})
}
