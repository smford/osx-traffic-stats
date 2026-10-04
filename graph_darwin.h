//go:build darwin

#ifndef GRAPH_DARWIN_H
#define GRAPH_DARWIN_H

#include <stdint.h>

void toggleTrafficGraphWindow(void);
void updateTrafficGraph(uint64_t upSpeed, uint64_t downSpeed, const char *topAppsSummary);
int isTrafficGraphVisible(void);
void resetTrafficGraphPeaks(void);
void setTrafficGraphOpacity(double opacity);
void setTrafficGraphAlwaysOnTop(int alwaysOnTop);
void snapTrafficGraphWindow(int corner);

#endif // GRAPH_DARWIN_H
