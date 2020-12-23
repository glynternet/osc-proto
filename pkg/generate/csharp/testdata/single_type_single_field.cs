using System;
using System.Collections.Generic;
using avvaunity.GOH.Unity.Message.Unmarshaller;

namespace namespaceBar {

    // Code generated by osc-proto (version unknown) DO NOT EDIT.

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