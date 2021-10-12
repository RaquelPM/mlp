#include "readData.h"
#include <fstream>
#include <iostream>

using namespace std;

int main(int argc, char ** argv){

    int dimension;
    double ** cost;
    
    readData(argc, argv, &dimension, &cost);

    ofstream file("distance_matrix");

    file << dimension << endl;
    for(int i = 1; i <= dimension; i++){
        for(int j = i+1; j <= dimension; j++){
            file << cost[i][j] << " ";
            if (cost[i][j] == 0.0) {
                cout << "yess\n";
                cout << i << " " << j << endl;
            }
        }
        if(i != dimension)
            file << endl;
        else
            file << "EOF" << endl;;

    }


    file.close();

    return 0;
}
