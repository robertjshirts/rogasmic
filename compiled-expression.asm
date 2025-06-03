movw sp, #0
movt sp, #0
movw r4, #0x0000
movt r4, #0x3f20
add r2, r4, #0x8
ldr r3, [r2]
orr r3, r3, #0x8
str r3, [r2]
movw r2, #0x0000 ; bits for our gpio pin (select the specific bit for the specific pin)
movt r2, #0x0020 ; bits for our gpio pin (select the specific bit for the specific pin)

reset:
movw r7, #0x4  ; blink 4 times before long delay
movt r7, #0x0  ; blink 4 times before long delay

blink_loop:
; Turn pin on
add r3, r4, #0x1c   ; r3 now points at GPSET
str r2, [r3]        ; write r2 to GPSET
stmea sp!, {r0-r12} ; store everything
movw r5, #0x4f20    ; store 1 mil
movt r5, #0x000f    ; store 1 mil
str r5, [sp]!, #4   ; pass value to delay
bl >delay            ; delay subroutine
ldmea sp!, {r0-r12} ; restore all registers

; Turn pin off
add r3, r4, #0x28   ; r3 now points at GPCLR
str r2, [r3]        ; write r2 to GPCLR
stmea sp!, {r0-r12} ; store everything
movw r5, #0x4f20    ; store 1 mil
movt r5, #0x000f    ; store 1 mil
str r5, [sp]!, #4   ; pass value to delay
bl >delay            ; delay subroutine
ldmea sp!, {r0-r12} ; restore all registers

; repeat
subs r7, r7, #1     ; decrement counter
bpl >blink_loop      ; blink again if counter is positive

; long delay
stmea sp!, {r0-r12} ; store everything
movw r5, #0x0900     ; store 4 mil
movt r5, #0x003d     ; store 4 mil
str r5, [sp]!, #4   ; pass value to delay
bl >delay            ; delay subroutine
ldmea sp!, {r0-r12} ; restore all registers

b >reset

delay:
ldr R6, [sp]!, #4
loop:
subs r6, r6, #0x01
bpl >loop

bx lr