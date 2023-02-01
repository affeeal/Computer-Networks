#include <iostream>
using namespace std;

int main(int argc, char *argv[]) 
{
	if (argc != 3) {
		cout << "Not enough arguments!" << endl;
	} else {
		int a = std::abs(std::stoi(argv[1]));
		int b = std::abs(std::stoi(argv[2]));
		cout << std::max(a, b) << endl;
	}
    return 0;
}

