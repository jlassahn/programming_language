
#include <stdio.h>

void clib_print_Real64(double x)
{
	printf("Real64: %g\n", x);
}

void clib_print_Int64(long long x)
{
	printf("Int64: %lld\n", x);
}

extern void janus_main(void);

int main(void)
{
	janus_main();
	return 0;
}

