using System;
using System.Collections.Generic;

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

    public class FooUnmarshaller {

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

    public interface IMessageHandler {
        void Handle(List<object> messageArgs);
    }

    public class BarRouter : IMessageHandler {
        private static readonly FooUnmarshaller BazFooUnmarshaller = new FooUnmarshaller();
        private static readonly FooUnmarshaller WhoopFooUnmarshaller = new FooUnmarshaller();

        private readonly IBarHandler _barHandler;

        private BarRouter() { }

        public BarRouter(IBarHandler barHandler) {
            _barHandler = barHandler;
        }

        public void Handle(List<object> data) {
            var route = data[0].ToString();
            switch (route) {
                case "baz":
                    _barHandler.HandleBazFoo(BazFooUnmarshaller.Unmarshal(data.GetRange(1, data.Count - 1)));
                    break;
                case "whoop":
                    _barHandler.HandleWhoopFoo(WhoopFooUnmarshaller.Unmarshal(data.GetRange(1, data.Count - 1)));
                    break;
                default:
                    throw new ArgumentOutOfRangeException("route", route, "Unsupported message routing arg");
            }
        }
    }

    public interface IBarHandler {
        void HandleBazFoo(Foo foo);
        void HandleWhoopFoo(Foo foo);
    }
}
