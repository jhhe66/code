#include <Ice/Ice.h>
#include "Printer.h"
using namespace std;
using namespace Demo;
int
main(int argc, char* argv[])
{
    int status = 0;
    Ice::CommunicatorPtr ic;
    try {
        ic = Ice::initialize(argc, argv);
        printf("init success.\n");
        //Ice::ObjectPrx base = ic->stringToProxy("SimplePrinter:default -h localhost -p 10000");
        Ice::ObjectPrx base = ic->propertyToProxy("client.Proxy");
        printf("base success.\n");
        PrinterPrx printer = PrinterPrx::checkedCast(base);
        printf("printer success.\n");
        if (!printer)
            throw "Invalid proxy";
        printer->printString("Hello World!");
        printf("printString success.\n");
    } catch (const Ice::Exception& ex) {
        cerr << ex << endl;
        status = 1;
    } catch (const char* msg) {
        cerr << msg << endl;
        status = 1;
    }
    if (ic)
        ic->destroy();
    return status;
}