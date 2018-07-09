#include <Ice/Ice.h>
#include <Printer.h>

using namespace std;
using namespace Demo;
class PrinterI : public Printer {
public:
    virtual void printString(const string& s, const Ice::Current&);
};


void 
PrinterI::
printString(const string& s, const Ice::Current&)
{
    cout << s << endl;
}

int
main(int argc, char* argv[])
{
    int status = 0;
    Ice::CommunicatorPtr ic;
    try {
        ic = Ice::initialize(argc, argv);
        //Ice::ObjectAdapterPtr adapter = ic->createObjectAdapterWithEndpoints("SimplePrinterAdapter", "default -h localhost -p 10000");
        Ice::ObjectAdapterPtr adapter = ic->createObjectAdapter("SimplePrinterAdapter");
        printf("adapter success.\n");
        Ice::ObjectPtr object = new PrinterI;
        printf("object success.\n");
        adapter->add(object, ic->stringToIdentity("SimplePrinter"));
        printf("adapter object success.\n");
        adapter->activate();
        printf("adapter activate success.\n");
        ic->waitForShutdown();
    } catch (const Ice::Exception& e) {
        cerr << e << endl;
        status = 1;
    } catch (const char* msg) {
        cerr << msg << endl;
        status = 1;
    }
    if (ic) {
        try {
            ic->destroy();
        } catch (const Ice::Exception& e) {
            cerr << e << endl;
            status = 1;
        }
    }
    return status;
}