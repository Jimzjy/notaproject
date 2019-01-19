#include <stdlib.h>

#define FACE_DETECT 0
#define BODY_DETECT 1

typedef struct Rect
{
    int x0;
    int y0;
    int x1;
    int y1;
} Rect;

typedef struct Rects
{
    unsigned int size;
    Rect* rects;
} Rects;

#ifdef __cplusplus
extern "C" {
#endif

#ifdef __cplusplus
typedef ncnn::Net* Ncnnnet;
#else
typedef void* Ncnnnet;
#endif

Ncnnnet newNcnnnet();
void ncnnnetLoad(char* param, char* model, Ncnnnet net);
Rects detectFromByte(unsigned char* data, int cols, int rows, Ncnnnet net, int mode);
Rects detectFromPath(Ncnnnet net, int mode, char* camPath);

#ifdef __cplusplus    
}
#endif