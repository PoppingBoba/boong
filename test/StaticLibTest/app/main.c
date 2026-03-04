#include <stdio.h>

extern int add__(int a, int b);

int main()
{
    printf("Test Build: %d\n", add__(1, 2));
    return 0;   
}