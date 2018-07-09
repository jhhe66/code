/********************************************
//文件名:os.h
//功能:系统包含文件
//作者:钟何明
//创建时间:2010.07.06
//修改记录:

*********************************************/

#ifndef _JW_OS_H_
#define _JW_OS_H_

#ifndef JW_PROTOCOL_API
    #ifdef _WIN32
        #ifdef JWPROTOCOL_EXPORTS
            #define JW_PROTOCOL_API extern "C" __declspec(dllexport)
        #else		
            #define JW_PROTOCOL_API extern "C" __declspec(dllimport)
        #endif

    #else
        #define JW_PROTOCOL_API extern "C"        
    #endif // _WIN32
#endif


#endif

