#include <iostream>
#include <fstream>
#include <vector>
#include "func.h"

using namespace std;

const char *readFile(const string &filename) {
    ifstream reader(filename.c_str(), ios::binary | ios::ate);
    if (!reader.is_open()) {
        return "";
    }
    vector<char> buffer;
    long long origSize = reader.tellg();
    long long size = origSize;
    long long pos = 0;

    reader.seekg(pos, ios::beg);

    buffer.resize(size);
    reader.read(buffer.data(), size);
    string result(buffer.begin(), buffer.end());
    reader.close();

    const char *out = result.c_str();
    return out;
}

/*FILE *f = fopen("icon.png", "rb+");
    fseek(f, 0L, SEEK_END);
    long filesize = ftell(f); // get file size
    fseek(f, 0L ,SEEK_SET); //go back to the beginning
    char* buffer = new char[filesize]; // allocate the read buf
    fread(buffer, 1, filesize, f);
    fclose(f);
    //delete[] buffer;
    */