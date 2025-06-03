MOVW sp, #0
MOVT sp, #0
MOVW R4, #0x0000
MOVT R4, #0x3F20
ADD R2, R4, #0x8
LDR R3, [R2]
ORR R3, R3, #0x8
STR R3, [R2]
MOVW R2, #0x0000  ; pin select
MOVT R2, #0x0020  ; pin select

start:
ADD R3, R4, #0x1C ; R3 now points to GPSET
STR R2, [R3]      ; write to GPSET
MOVW R5, #0x0900  ; store 4 mil delay
MOVT R5, #0x003D  ; store 4 mil delay
STR R5, [sp]!, #4 ; store delay at the stack pointer, then increment by 4
BL delay

ADD R3, R4, #0x28 ; R3 now points to GPCLR
STR R2, [R3]      ; write to GPCLR
MOVW R5, #0x4F20  ; store 1 mil delay
MOVT R5, #0x000F  ; store 1 mil delay
STR R5, [sp]!, #4 ; store delay at the stack pointer, then increment by 4
BL delay

B start
delay:
LDR R6, [sp]!, #4 ; load into r6, the value at the stack pointer, then decrement the sp by 4
loop:
SUBS R6, R6, #0x01
BPL loop

BX lr