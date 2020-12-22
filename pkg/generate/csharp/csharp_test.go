package csharp_test

import (
	"testing"

	"github.com/glynternet/osc-proto/pkg/generate/csharp"
	"github.com/glynternet/osc-proto/pkg/generate/generatetest"
	"github.com/glynternet/osc-proto/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEmptyTypesShouldYieldEmptyFile(t *testing.T) {
	in := types.Types{}
	var expected map[string][]byte
	out, err := csharp.Generator{}.Generate(in)
	require.NoError(t, err)
	assert.Equal(t, expected, out)
}

func TestMultipleTypesShouldYieldUnsupportedError(t *testing.T) {
	in := types.Types{
		"foo": {},
		"bar": {},
	}

	_, err := csharp.Generator{}.Generate(in)
	require.Error(t, err)
}

func TestNonboolFieldTypeShouldYieldError(t *testing.T) {
	in := types.Types{
		"foo": {{
			FieldName: "fieldFoo",
			FieldType: "nonbool",
		}},
	}
	_, err := csharp.Generator{Namespace: "foo"}.Generate(in)
	require.EqualError(t, err, "type:foo contains non-bool field type:nonbool for field:fieldFoo")
}

func TestSingleTypeSingleFieldShouldYieldResult(t *testing.T) {
	in := types.Types{
		"foo": {{
			FieldName: "fieldFoo",
			FieldType: "bool",
		}},
	}
	out, err := csharp.Generator{Namespace: "namespaceBar"}.Generate(in)
	const expected = `using System;
using System.Collections.Generic;
using avvaunity.GOH.Unity.Message.Unmarshaller;

namespace namespaceBar {

    public readonly struct Foo {
        private readonly bool _fieldFoo;

        public Foo(bool fieldFoo) {
            _fieldFoo = fieldFoo;
        }

        public bool FieldFoo() {
            return _fieldFoo;
        }
    }

    public class FooUnmarshaller : IMessageUnmarshaller<Foo> {

        // <fieldFoo:bool>
        public Foo Unmarshal(List<object> data) {
            if (data.Count != 1) {
                throw new ArgumentException($"Expected 1 item in arg list but got {data.Count}");
            }
            var fieldFoo = ParseBool(data[0].ToString());
            return new Foo(fieldFoo);
        }

        private static bool ParseBool(string value) {
            try {
                return int.Parse(value) != 0;
            }
            catch (Exception e) {
                throw new ArgumentException($"cannot parse {value} to bool", e);
            }
        }
    }
}
`
	require.NoError(t, err)
	generatetest.AssertEqualContentLayout(t, map[string]string{
		"foo.cs": expected,
	}, out)
}

func TestSingleTypeMultipleFieldsShouldYieldResult(t *testing.T) {
	in := types.Types{
		"foo": {{
			FieldName: "fieldFoo",
			FieldType: "bool",
		}, {
			FieldName: "fieldBar",
			FieldType: "bool",
		}},
	}
	out, err := csharp.Generator{Namespace: "namespaceBar"}.Generate(in)
	const expected = `using System;
using System.Collections.Generic;
using avvaunity.GOH.Unity.Message.Unmarshaller;

namespace namespaceBar {

    public readonly struct Foo {
        private readonly bool _fieldFoo;
        private readonly bool _fieldBar;

        public Foo(bool fieldFoo, bool fieldBar) {
            _fieldFoo = fieldFoo;
            _fieldBar = fieldBar;
        }

        public bool FieldFoo() {
            return _fieldFoo;
        }

        public bool FieldBar() {
            return _fieldBar;
        }
    }

    public class FooUnmarshaller : IMessageUnmarshaller<Foo> {

        // <fieldFoo:bool>
        // <fieldBar:bool>
        public Foo Unmarshal(List<object> data) {
            if (data.Count != 2) {
                throw new ArgumentException($"Expected 2 item in arg list but got {data.Count}");
            }
            var fieldFoo = ParseBool(data[0].ToString());
            var fieldBar = ParseBool(data[1].ToString());
            return new Foo(fieldFoo, fieldBar);
        }

        private static bool ParseBool(string value) {
            try {
                return int.Parse(value) != 0;
            }
            catch (Exception e) {
                throw new ArgumentException($"cannot parse {value} to bool", e);
            }
        }
    }
}
`
	require.NoError(t, err)
	generatetest.AssertEqualContentLayout(t, map[string]string{
		"foo.cs": expected,
	}, out)
}
