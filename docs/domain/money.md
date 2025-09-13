# Money Value Object

## Representation
- Amount stored as signed 64-bit integer representing cents.
- Currency stored as ISO code (`USD` in v1).
- No floating point math; arithmetic uses integer operations.

## Operations
- Addition/subtraction **MUST** check currency equality.
- Formatting to user locale **SHOULD** occur at UI layer.
- Conversions between currencies are out of scope for v1.

## Invariants
- Overflow **MUST NOT** occur; guard with bounds checks.
- Zero amount represents absence of value but retains currency.

## Examples
```
Money(2500, "USD")  # $25.00
Money(-120000, "USD")  # -$1,200.00
```

