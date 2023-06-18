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
	popq %rax
	ret
.global _start
_start:
	callq main.main
	movq %rax, %rdi
	movq $60, %rax
	syscall
