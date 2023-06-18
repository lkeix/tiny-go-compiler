.data
a:
  .quad 1
b:
  .quad 1

.text
main.main:
  movq a+0(%rip), %rax
  movq b+0(%rip), %rdi
  addq %rdi, %rax
  pushq %rax
  ret
.global _main
_main:
  callq main.main
  popq %rax
  movq 0x60, %rax
  syscall
