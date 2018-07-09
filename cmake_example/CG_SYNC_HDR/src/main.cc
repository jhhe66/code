#include <iostream>
#include <boost/scoped_ptr.hpp>

#include "HandlerBase.h"
#include "HandlerFactory.h"

#include <sstream>


using namespace std;

unsigned int
StrToHex(const string& str)
{
    istringstream iss(str);

    unsigned int i = 0;
    
    iss >> hex >> i;

    return i;
}


int
main(int argc, char** argv)
{
    if ( argc < 2 || argc > 2)
    {   
        cout << "usage:\n" 
            << "cg_sync_har handlerID\n" << "example:\n" 
                << "cg_sync_har 0x0A50\n";
        
        return -1;
    }


    unsigned short handlerID = StrToHex(*(argv + 1));// atoi(*(argv + 1));
    
    Handler::CHandlerFactory factory;

    boost::scoped_ptr< Handler::CHandlerBase > handler(factory.CreateHandler(handlerID));

    cout << "Return: " << handler->Sync() << endl;
}
