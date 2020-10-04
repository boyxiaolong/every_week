#include <iostream>
using namespace std;

class A
{
        public:
                void test()
                {
                        std::cout << "class A test\n";
                }
                virtual int test1()
                {
                        std::cout << "class A test1\n";
                        return 0;
                }

};

class B : public A
{
        public:
                void test()
                {
                        std::cout << "class B test\n";
                }
                virtual int test1()
                {
                        std::cout << "class B test1\n";
                        return 0;
                }
};
int main()
{
        A* pa = new B();
        pa->test();
        pa->test1();
}
