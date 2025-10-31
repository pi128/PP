

#include <stdio.h>
#include <string.h>

int main(void) {

    FILE *fptr;
    fptr = fopen("test.txt", "r");

    if (fptr == NULL) {
    printf("Error: Could not open file.\n");
    return 1;
        }

    identifier = [a-zA-Z_]\w*\b
    char *words[]
    
    
    
    while ((ch = getc(fptr)) != EOF) {
        if (ch == ' '){
            continue;
        }
        
        else {
            printf("%c", ch);
        }
        
    }
    return 0;
  
    fclose(fptr);
}