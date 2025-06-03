MOVW sp, #0
MOVT sp, #0
MOVW R4, #0x0000
MOVT R4, #0x3F20
ADD R2, R4, #0x8
LDR R3, [R2]
ORR R3, R3, #0x8
STR R3, [R2]
MOVW R2, #0x0000 ; bits for our gpio pin (select the specific bit for the specific pin)
MOVT R2, #0x0020 ; bits for our gpio pin (select the specific bit for the specific pin)

start:
ADD R3, R4, #0x1C ; change R3 so it points at GPSET
STR R2, [R3] ; write the correct bit (R2) into the GPSET address (R3)

STMEA sp!, {R0-R12} ; Store all registers
MOVW R5, #0x0900 ; store 4 mil
MOVT R5, #0x003D ; store 4 mil
STR R5, [sp!], #4 ; copy delay onto stack at the current sp, then incr by 4
BL delay ; branch to subroutine
LDMEA sp!, {R0-R12} ; Restore all registers

ADD R3, R4, #0x28 ; change R3 so it points at GPCLR
STR R2, [R3] ; write the correct bit (R2) into the GPCLR address (R3)

STMEA sp!, {R0-R12} ; store all registers
MOVW R5, #0x4F20 ; store 1 mil
MOVT R5, #0x000F ; store 1 mil
STR R5, [sp]!, #4 ; copy delay onto the stack at the current sp, then incr by 4
BL delay ; branch to subroutine
LDMEA sp! {R0-R12} ; restore all registers

B start ; restart

delay:
ldr R6, [sp]!, #4
loop:
SUBS R6, R6, #0x01
BPL loop

BX lr