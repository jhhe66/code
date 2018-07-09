#include <stdio.h>
#include <iostream>
#include <fstream>
#include <locale>
#include <stdlib.h>
#include <time.h>
#include <string.h>

using namespace std;
 
const int testsize = 50000;
 
int main(){
    char cstr[] = "这是实验字符串这是实验字符串这是实验字符串这是实验字符串这是实验字符串这是实验字符串这是实验字符串这是实验字符串";
    wchar_t wstr[] = L"这是实验字符串这是实验字符串这是实验字符串这是实验字符串这是实验字符串这是实验字符串这是实验字符串这是实验字符串";
    char buffer[200];
    int cstrlen = strlen(cstr);
    locale defloc("");
 
    clock_t start;
    FILE *cfile;
    ofstream fout;
    wofstream wfout;
 
    // pure text ...........................
    start = clock();
    cfile = fopen("text_out_C_fputs.txt", "w");
    for(int i = 0; i < testsize; ++i){
        fputs(cstr, cfile);
        fputc('\n', cfile);
    }
    fclose(cfile);
    printf("Text out C fputs: %d\n", clock()-start);
 
    start = clock();
    cfile = fopen("test_out_C_fprintf.txt", "w");
    for(int i = 0; i < testsize; ++i){
        fprintf(cfile, "%s\n", cstr);
    }
    fclose(cfile);
    printf("Text out C fprintf: %d\n", clock()-start);
 
    fout.clear();
    start = clock();
    fout.open("text_out_Cpp_ofstream.txt");
    for(int i = 0; i < testsize; ++i){
        fout << cstr << '\n';
    }
    fout.close();
    printf("Text out Cpp ofstream: %d\n", clock()-start);
 
    fout.clear();
    start = clock();
    fout.open("text_out_Cpp_rdbuf.txt");
    for(int i = 0; i < testsize; ++i){
        fout.rdbuf()->sputn(cstr, cstrlen);
        fout.rdbuf()->sputc('\n');
    }
    fout.close();
    printf("Text out Cpp rdbuf: %d\n", clock()-start);
 
    // wchar_t text ...............................
    wfout.clear();
    start = clock();
    wfout.open("wtext_out_Cpp_wofs_defloc.txt");
    wfout.imbue(defloc);
    for(int i = 0; i < testsize; ++i){
        wfout << wstr << L'\n';
    }
    wfout.close();
    printf("wText out Cpp wofs with default locale: %d\n", clock()-start);
 
    fout.clear();
    start = clock();
    char *next;
    const wchar_t *wnext;
    mbstate_t st;
    fout.open("wtext_out_Cpp_codecvt_facet.txt");
    for(int i = 0; i < testsize; ++i){
        use_facet<codecvt<wchar_t, char, mbstate_t> >(defloc).out(
            st, wstr, wstr+sizeof(wstr)/2-1, wnext,
            buffer, buffer+sizeof(buffer)-1, next);
        fout.rdbuf()->sputn(buffer, next-buffer);
        fout.rdbuf()->sputc('\n');
    }
    fout.close();
    printf("wText out Cpp ofs with codecvt facet: %d\n", clock()-start);
 
    fout.clear();
    start = clock();
    fout.open("wtext_out_Cpp_ofs_WinAPI.txt");
    for(int i = 0; i < testsize; ++i){
        WideCharToMultiByte(CP_ACP, 0, wstr, -1, buffer, 200, NULL, NULL);
        fout << buffer << '\n';
    }
    fout.close();
    printf("wText out Cpp ofs with WideCharToMultiByte API: %d\n",
        clock()-start);
 
    // Format out ...........................................
    srand((unsigned)time(NULL));
 
    char datastr[] = "TestDataString实验格式化字符串";
    start = clock();
    cfile = fopen("format_data_out_C_fprintf.txt", "w");
    for(int i = 0; i < testsize; ++i){
        fprintf(cfile, "%d %lf %s\n",
            rand(), double(rand())/RAND_MAX, datastr);
    }
    fclose(cfile);
    printf("Format data out C fprintf: %d\n", clock()-start);
 
    fout.clear();
    start = clock();
    fout.open("format_data_out_Cpp_ofstream.txt");
    for(int i = 0; i < testsize; ++i){
        fout << rand() << ' ' << double(rand())/RAND_MAX << ' '
            << datastr << '\n';
    }
    fout.close();
    printf("Format data out Cpp ofstream: %d\n", clock()-start);
 
    fout.clear();
    start = clock();
    fout.open("format_data_out_Cpp_ofs_facet.txt");
    for(int i = 0; i < testsize; ++i){
        use_facet<num_put<char> >(locale::classic()).put(
            ostreambuf_iterator<char>(fout), fout, ' ', (long)rand());
        fout.rdbuf()->sputc(' ');
        use_facet<num_put<char> >(locale::classic()).put(
            ostreambuf_iterator<char>(fout), fout, ' ',
            double(rand())/RAND_MAX);
        fout.rdbuf()->sputc(' ');
        fout.rdbuf()->sputn(datastr, sizeof(datastr) - 1);
        fout.rdbuf()->sputc('\n');
    }
    fout.close();
    printf("Format data out Cpp ofs with facet: %d\n", clock()-start);
}
