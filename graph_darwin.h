//go:build darwin

#ifndef GRAPH_DARWIN_H
#define GRAPH_DARWIN_H

#include <stdint.h>

void toggleTrafficGraphWindow(void);
void updateTrafficGraph(uint64_t upSpeed, uint64_t downSpeed);
int isTrafficGraphVisible(void);
void resetTrafficGraphPeaks(void);

#endif // GRAPH_DARWIN_H
