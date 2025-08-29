# RFC: Enum Support in osc-proto Type System

## Summary

This RFC proposes adding enum support to the osc-proto type system, enabling type-safe enumeration values in OSC message fields. Enums will be defined in the type definitions and supported by all generators (Go, C#).

## Motivation

Currently, the type system supports primitive types (`bool`, `int32`, `float32`, `string`), but lacks support for enumerated values. This limitation forces developers to use raw strings or integers for what should be type-safe enumeration values, leading to:

- Runtime errors from invalid enum values
- Lack of type safety and IDE support
- Unclear API contracts
- No validation of enum values

## Design

### Type System Changes

#### Enum Definition Syntax

Enums will be defined as a new top-level section in the YAML configuration:

```yaml
enums:
  PlayerState:
    values:
      - IDLE
      - RUNNING
      - JUMPING
      - FALLING
  Direction:
    values:
      - NORTH
      - SOUTH
      - EAST
      - WEST

types:
  PlayerUpdate:
    - state: PlayerState
    - direction: Direction
    - x: float32
    - y: float32
```

#### Type System Implementation

**New types in `pkg/types/types.go`:**

```go
type Enums map[EnumName]EnumDefinition

type EnumName string

type EnumDefinition struct {
    Values []EnumValue `yaml:"values"`
}

type EnumValue string

// Updated Definitions struct
type Definitions struct {
    Types   Types   `yaml:"types"`
    Enums   Enums   `yaml:"enums"`
    Routers Routers `yaml:"routers"`
}
```

**Validation:**
- Enum names must be valid identifiers
- Enum values must be valid identifiers
- Enum values within a definition must be unique
- Field types referencing enums must exist in the enum definitions

### Generator Implementation

#### Go Generator Changes

**Enum Generation:**
```go
// Generated Go code example
type PlayerState int32

const (
    PlayerStateIDLE PlayerState = iota
    PlayerStateRUNNING
    PlayerStateFUMPING  
    PlayerStateFALLING
)

func (e PlayerState) String() string {
    switch e {
    case PlayerStateIDLE: return "IDLE"
    case PlayerStateRUNNING: return "RUNNING"
    case PlayerStateFUMPING: return "FUMPING"
    case PlayerStateFALLING: return "FALLING"
    default: return fmt.Sprintf("PlayerState(%d)", int32(e))
    }
}

func ParsePlayerState(s string) (PlayerState, error) {
    switch s {
    case "IDLE": return PlayerStateIDLE, nil
    case "RUNNING": return PlayerStateRUNNING, nil
    case "FUMPING": return PlayerStateFUMPING, nil
    case "FALLING": return PlayerStateFALLING, nil
    default: return 0, fmt.Errorf("invalid PlayerState: %s", s)
    }
}
```

**Field Handling:**
- Enum fields in message args will be converted to their string representation
- Parsing will use the generated Parse functions with error handling

#### C# Generator Changes  

**Enum Generation:**
```csharp
public enum PlayerState {
    IDLE = 0,
    RUNNING = 1,
    FUMPING = 2,
    FALLING = 3
}

public static class PlayerStateExtensions {
    public static PlayerState ParsePlayerState(string value) {
        switch (value) {
            case "IDLE": return PlayerState.IDLE;
            case "RUNNING": return PlayerState.RUNNING; 
            case "FUMPING": return PlayerState.FUMPING;
            case "FALLING": return PlayerState.FALLING;
            default: throw new ArgumentException($"Invalid PlayerState: {value}");
        }
    }
    
    public static string ToOSCString(this PlayerState state) {
        return state.ToString();
    }
}
```

**Field Handling:**
- Enum fields will use the generated Parse methods in unmarshalling
- Enum values will be serialized as strings in OSC messages

### OSC Protocol Representation

Enums will be transmitted as strings over OSC to maintain human readability and debugging capabilities. This follows OSC best practices and ensures compatibility with existing OSC implementations.

### Migration Path

This is a purely additive change:
- Existing YAML configurations without enums continue to work unchanged
- New enum fields can be added incrementally
- No breaking changes to existing generated code

## Implementation Plan

### Phase 1: Core Type System
1. Add enum types to `pkg/types/types.go`
2. Update YAML unmarshaling to support enum definitions
3. Add validation for enum definitions and references
4. Update `generate.Definitions` struct

### Phase 2: Go Generator
1. Add enum generation to Go templates
2. Update field handling for enum types
3. Generate parse/string conversion functions
4. Add enum support to `fieldArgFuncs()`

### Phase 3: C# Generator  
1. Add enum generation to C# templates
2. Update field handling and type conversions
3. Generate parse/extension methods
4. Update `typeConversions()` map

### Phase 4: Testing & Documentation
1. Add comprehensive test cases for all enum scenarios
2. Update existing tests to include enum examples  
3. Add validation error test cases
4. Update README with enum usage examples

## Examples

### Configuration Example
```yaml
enums:
  LogLevel:
    values:
      - DEBUG
      - INFO  
      - WARN
      - ERROR

types:
  LogMessage:
    - level: LogLevel
    - message: string
    - timestamp: int32

routers:
  logger:
    log: LogMessage
```

### Generated Go Usage
```go
// Send a log message
msg := LoggerLog(LogLevelINFO, "Application started", 1234567890)
client.Send(msg)
```

### Generated C# Usage  
```csharp
// Handle a log message
public void HandleLog(LogMessage msg) {
    Console.WriteLine($"[{msg.Level()}] {msg.Message()}");
}
```

## Considerations

### Backwards Compatibility
- Fully backwards compatible - existing configurations work unchanged
- Generated code remains compatible for non-enum types

### Performance Impact
- Minimal performance overhead from string parsing
- Enum validation happens at generation time, not runtime
- OSC message size unchanged (strings vs integers trade-off favors debugging)

### Error Handling
- Parse errors for invalid enum values are handled at unmarshaling time
- Clear error messages indicating invalid enum values and expected values

### Extensibility
- Design supports future enum features (custom values, aliases)
- Could support integer-backed enums in future if needed

## Alternatives Considered

1. **Integer-based enums**: Rejected due to loss of human readability in OSC messages
2. **Inline enum definitions**: Rejected due to lack of reusability across types
3. **String constants**: Current workaround, but lacks type safety and validation

## Conclusion

Adding enum support will significantly improve the type safety and developer experience of osc-proto while maintaining full backwards compatibility. The proposed design is minimal, extensible, and follows established patterns in both Go and C# ecosystems.