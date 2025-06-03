# ARM Assembly Syntax Guide

## Data Movement Instructions

### MOV - Move Data
```
MOV{cond}{S} Rd, <Operand2>
MOVW{cond} Rd, #<imm16>
MOVT{cond} Rd, #<imm16>
```
where:
- **MOV** moves data from one register to another or moves an immediate value
- **MOVW** moves a 16-bit immediate value to the lower 16 bits of the destination register
- **MOVT** moves a 16-bit immediate value to the upper 16 bits of the destination register
- **{cond}** two-character condition mnemonic (see Condition Codes section)
- **{S}** if S is present, the instruction updates the condition flags
- **Rd** is the destination register (R0-R15)
- **#<imm16>** is a 16-bit immediate value (0x0000 to 0xFFFF)
- **<Operand2>** can be a register or immediate value

### LDR/STR - Load/Store Register
```
<LDR|STR>{cond}{B}{T} Rd, <Address>
```
where:
- **LDR** load from memory into a register
- **STR** store from a register into memory
- **{cond}** two-character condition mnemonic
- **{B}** if B is present then byte transfer, otherwise word transfer
- **{T}** if T is present the W bit will be set in a post-indexed instruction, forcing non-privileged mode for the transfer cycle. T is not allowed when a pre-indexed addressing mode is specified or implied
- **Rd** is the destination/source register
- **<Address\>** can be:
  - **[Rn]** - simple register indirect addressing
  - **[Rn]!, #<expression>** - post-indexed with immediate offset and writeback

### LDM/STM - Load/Store Multiple Registers
```
<LDM|STM>{cond}<EA> Rn{!}, <register_list>
```
where:
- **LDM** load multiple registers from memory
- **STM** store multiple registers to memory
- **{cond}** two-character condition mnemonic
- **<EA>** addressing mode (REQUIRED - sets P and U bits):
  - **IA** - Increment After (P=0, U=1)
  - **IB** - Increment Before (P=1, U=1) 
  - **EA** - Empty Ascending (equivalent to IA for LDM, DB for STM)
- **Rn** is the base register
- **{!}** if present, write back the final address to the base register
- **<register_list\>** is a list of registers in braces, e.g., {R0-R12} or {R1,R3,R5}

## Arithmetic Instructions

### ADD/SUB - Addition/Subtraction
```
<ADD|SUB>{cond}{S} Rd, Rn, <Operand2>
```
where:
- **ADD** performs addition
- **SUB** performs subtraction
- **{cond}** two-character condition mnemonic
- **{S}** if S is present, the instruction updates condition flags
- **Rd** is the destination register
- **Rn** is the first operand register
- **<Operand2>** can be:
  - **#<immediate>** - immediate value
  - **Rm** - register
  - **Rm, <shift>** - shifted register

### SUBS - Subtract and Set Flags
```
SUBS{cond} Rd, Rn, <Operand2>
```
where:
- **SUBS** always sets condition flags (S bit implicitly set)
- Parameters same as ADD/SUB above

## Logical Instructions

### ORR - Logical OR
```
ORR{cond}{S} Rd, Rn, <Operand2>
```
where:
- **ORR** performs bitwise logical OR operation
- **{cond}** two-character condition mnemonic
- **{S}** if S is present, updates condition flags
- **Rd** is the destination register
- **Rn** is the first operand register
- **<Operand2>** is the second operand (immediate or register)

## Branch Instructions

### B - Branch
```
B{cond} <label>
```
where:
- **B** unconditional or conditional branch
- **{cond}** two-character condition mnemonic
- **<label\>** is the target address or label (label identifiers must be preceded by >)

### BL - Branch with Link
```
BL{cond} <label>
```
where:
- **BL** branch with link (saves return address in LR/R14)
- **{cond}** two-character condition mnemonic
- **<label\>** is the target address or label (label identifiers must be preceded by >)

### BX - Branch and Exchange
```
BX{cond} Rm
```
where:
- **BX** branch and exchange instruction sets
- **{cond}** two-character condition mnemonic
- **Rm** is the register containing the target address

## Condition Codes

| Code | Flags | Meaning |
|------|-------|---------|
| EQ | Z=1 | Equal |
| PL | N=0 | Plus/Positive or Zero |
| AL | - | Always (default) |

## Stack Operations

### Stack Pointer Conventions
- **sp** (R13) is the stack pointer register
- **Full stack**: stack pointer points to the last pushed item
- **Empty stack**: stack pointer points to the next free location
- **Ascending**: stack grows toward higher memory addresses
- **Descending**: stack grows toward lower memory addresses

### Common Stack Usage Patterns
```assembly
; Save registers to stack (push)
STMEA sp!, {R0-R12}    ; Store multiple, empty ascending with writeback

; Restore registers from stack (pop)  
LDMEA sp!, {R0-R12}    ; Load multiple, empty ascending with writeback

; Push single value with post-increment
STR R5, [sp]!, #4      ; Store and increment sp by 4

; Pop single value with post-decrement
LDR R6, [sp]!, #4      ; Load and decrement sp by 4
```

## Special Notes

1. **Labels**:
   - Label declarations must be followed by a colon (e.g., my_label:)
   - Label identifiers must be preceded by `>` (e.g., >my_label)

2. **Immediate Values**: 
   - Must be prefixed with # (e.g., #0x1000)
   - 16-bit immediates for MOVW/MOVT: #0x0000 to #0xFFFF

3. **Memory Addressing**:
   - Byte addresses are incremented by 1
   - Word addresses are incremented by 4
   - All addresses should be word-aligned for word transfers

4. **LDM/STM EA Postfix**: 
   - The EA (Empty Ascending) postfix is required and sets the P and U bits correctly
   - For stack operations, EA assumes stack grows upward (ascending)
   - Register list order doesn't matter - lowest numbered registers are always stored/loaded first 