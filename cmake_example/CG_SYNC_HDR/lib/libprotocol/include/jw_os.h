/********************************************
//�ļ���:os.h
//����:ϵͳ�����ļ�
//����:�Ӻ���
//����ʱ��:2010.07.06
//�޸ļ�¼:

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

