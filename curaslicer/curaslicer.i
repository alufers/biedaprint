/* File : example.i */
%module(directors="1")  curaslicer

%{
#include "CuraInterface/CuraInterface.h"
%}


%ignore cura;

%include "CuraInterface/CuraInterface.h"



