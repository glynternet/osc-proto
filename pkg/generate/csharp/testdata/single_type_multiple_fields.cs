using System;
using System.Collections.Generic;

namespace namespaceBar {

    // Code generated by osc-proto (version 🧨) DO NOT EDIT.

    public readonly struct Foo {
        private readonly bool _fieldFoo;
        private readonly int _fieldBar;
        private readonly string _fieldString;
        private readonly float _fieldFloat32;

        public Foo(bool fieldFoo, int fieldBar, string fieldString, float fieldFloat32) {
            _fieldFoo = fieldFoo;
            _fieldBar = fieldBar;
            _fieldString = fieldString;
            _fieldFloat32 = fieldFloat32;
        }

        public bool FieldFoo() {
            return _fieldFoo;
        }

        public int FieldBar() {
            return _fieldBar;
        }

        public string FieldString() {
            return _fieldString;
        }

        public float FieldFloat32() {
            return _fieldFloat32;
        }
    }

    public class FooUnmarshaller {

        // <fieldFoo:bool>
        // <fieldBar:int32>
        // <fieldString:string>
        // <fieldFloat32:float32>
        public Foo Unmarshal(List<object> data) {
            if (data.Count != 4) {
                throw new ArgumentException($"Expected 4 item in arg list but got {data.Count}");
            }
            var fieldFoo = ParseBool(data[0].ToString());
            var fieldBar = int.Parse(data[1].ToString());
            var fieldString = (data[2].ToString());
            var fieldFloat32 = float.Parse(data[3].ToString());
            return new Foo(fieldFoo, fieldBar, fieldString, fieldFloat32);
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
