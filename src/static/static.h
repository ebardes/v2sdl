#include <stdint.h>
#include <unistd.h>
#include <stdlib.h>

typedef struct item {
    char *name;
    void *data;
    int64_t len;
    time_t time;
} item;

item* finditem(char *name);
int hasitem(char *name);