#include <static.h>
#include <string.h>

extern item staticnames[];
extern int64_t staticlen;

item* finditem(char *name) {
    
    int n = (int)(staticlen);

    for (int i = 0; i < n; i++) {
        if (strcmp(name, staticnames[i].name) == 0) {
            return staticnames+i;
        }

    }
    return 0;
}

int hasitem(char *name) {
    item* i = finditem(name);
    return i != NULL;
}

asm ( ".include \"static.s\"");