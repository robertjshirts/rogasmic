MOVW R4, #0x0000
MOVT R4, #0x3F20
ADD R2, R4, #0x8
LDR R3, [R2]
ORR R3, R3, #0x8
STR R3, [R2]

ADD R3, R4, #0x1C       ;       -12
MOVW R2, #0x0000        ;       -11
MOVT R2, #0x0020        ;       -10
STR R2, [R3]           ;       -9
BL #0x000005           ; -2    -8

ADD R3, R4, #0x28       ; -1    -7
MOVW R2, #0x0000        ;  0    -6
MOVT R2, #0x0020        ;  1    -5
STR R2, [R3]           ;  2    -4
BL #0x000000           ;  3 -2 -3

B #0xFFFFF4             ;  4 -1 -2

MOVW R5, #0xFFFF        ;  5  0
MOVT R5, #0x000F

SUBS R5, R5, #0x01     ;  -3
BPL #0xFFFFFD          ;  -2

BX R14