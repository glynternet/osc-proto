# osc-proto

A code generator that creates OSC (Open Sound Control) message handling code for Go and C# from YAML type definitions.

## Overview

osc-proto generates type-safe OSC message constructors and routers from YAML configurations. It supports:

- **Go code generation**: Creates message argument constructors and OSC message functions
- **C# code generation**: Creates immutable structs, unmarshallers, and message routers with interfaces

## Type Definitions

The generator operates on YAML files defining types and optional routers:

```yaml
types:
  message_name:
    - field_name: field_type
    - another_field: string

routers:
  router_name:
    route_name: message_name
```

### Supported Field Types

- `bool` - boolean values (converted to int32 in Go for OSC compatibility)
- `int32` - 32-bit integers  
- `float32` - 32-bit floating point numbers
- `string` - text strings

## Go Generator Behavior

The Go generator (`pkg/generate/golang/`) creates:

1. **Message argument constructors**: Functions that convert typed parameters to `[]interface{}`
   - Named as `{TypeName}MessageArgs()`
   - Handle bool-to-int32 conversion for OSC compatibility via `boolInt32()`

2. **Router message functions** (if routers defined): Functions that create `osc.Message` structs
   - Named as `{RouterName}{RouteName}()`
   - Include OSC address as `/{router_name}`
   - Prepend route name to message arguments

Example output for a type `foo` with field `enabled: bool`:
```go
func FooMessageArgs(enabled bool) []interface{} {
    return []interface{}{
        boolInt32(enabled),
    }
}
```

## C# Generator Behavior  

The C# generator (`pkg/generate/csharp/`) creates:

1. **Immutable structs**: `readonly struct` types with private fields and getter methods
   - Constructor takes all field parameters
   - Getter methods named as `{FieldName}()` with proper casing

2. **Unmarshaller classes**: Convert `List<object>` data to typed structs
   - Named as `{TypeName}Unmarshaller`
   - Include field type information as comments
   - Handle type parsing (bool from int, etc.)

3. **Router classes** (if routers defined): Implement message routing with type safety
   - Implement `IMessageHandler` interface
   - Include handler interface `I{RouterName}Handler` with typed methods
   - Route switching based on first message argument

Example output for a type `foo` with field `enabled: bool`:
```csharp
public readonly struct Foo {
    private readonly bool _enabled;
    public Foo(bool enabled) { _enabled = enabled; }
    public bool Enabled() { return _enabled; }
}
```

## Error Handling

- **Empty types**: Returns nil/null output (no code generated)
- **Unsupported field types**: Returns descriptive errors referencing field and type names  
- **Missing router references**: Returns errors when routers reference undefined types

## Package Structure

- `pkg/generate/generator.go`: Core `Generator` interface and `Definitions` struct
- `pkg/types/`: Type definitions with YAML unmarshalling and validation
- `pkg/routers/`: Router definitions linking routes to types  
- `pkg/generate/golang/`: Go code generator implementation
- `pkg/generate/csharp/`: C# code generator implementation