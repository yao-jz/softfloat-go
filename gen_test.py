# generate test for softfloat

import struct
import os

def generate_pair():
    random_bytes_32 = os.urandom(4)
    random_float_32 = struct.unpack('f', random_bytes_32)[0]

    random_bytes_321 = os.urandom(4)
    random_float_321 = struct.unpack('f', random_bytes_321)[0]

    random_bytes_64 = os.urandom(8)
    random_float_64 = struct.unpack('d', random_bytes_64)[0]

    random_bytes_641 = os.urandom(8)
    random_float_641 = struct.unpack('d', random_bytes_641)[0]

    return random_float_32, random_float_321, random_float_64, random_float_641

def f2i(f):
    binary_data = struct.pack('f', f)
    return int.from_bytes(binary_data, 'little')

def f2h(f):
    binary_data = struct.pack('f', f)
    return hex(int.from_bytes(binary_data, 'little'))

def i2f(i):
    binary_data = i.to_bytes(4, 'little')
    return struct.unpack('f', binary_data)[0]

def d2i(d):
    binary_data = struct.pack('d', d)
    return int.from_bytes(binary_data, 'little')

def d2h(d):
    binary_data = struct.pack('d', d)
    return hex(int.from_bytes(binary_data, 'little'))

def i2d(i):
    binary_data = i.to_bytes(8, 'little')
    return struct.unpack('d', binary_data)[0]

def gen_test(filename, number):

    outfile = open(filename, "w")

    outfile.write("""
    package test

    import (
        "fmt"
        "testing"
        . "github.com/yao-jz/softfloat-go/comp"
    )

    func Test"""
    )
    outfile.write(str(number))
    outfile.write("""(t *testing.T) {
                  """)


    n1, n2, d1, d2 = generate_pair()
    n1_hex = f2h(n1)
    n2_hex = f2h(n2)
    d1_hex = d2h(d1)
    d2_hex = d2h(d2)
    outfile.write("// test for 32 bits floating point number "+str(n1)+" and "+str(n2)+"\n")
    outfile.write("// test for 64 bits floating point number "+str(d1)+" and "+str(d2)+"\n\n")

    # addf32
    outfile.write("fmt.Println(\"test the addf32\")\n")
    outfile.write("fmt.Println(F32_add("+n1_hex+", "+n2_hex+"))\nfmt.Println(uint64("+str(f2i(n1+n2))+"))\n")
    # subf32
    outfile.write("fmt.Println(\"test the subf32\")\n")
    outfile.write("fmt.Println(F32_sub("+n1_hex+", "+n2_hex+"))\nfmt.Println(uint64("+str(f2i(n1-n2))+"))\n")
    # mulf32
    outfile.write("fmt.Println(\"test the mulf32\")\n")
    outfile.write("fmt.Println(F32_mul("+n1_hex+", "+n2_hex+"))\nfmt.Println(uint64("+str(f2i(n1*n2))+"))\n")
    # divf32
    outfile.write("fmt.Println(\"test the divf32\")\n")
    outfile.write("fmt.Println(F32_div("+n1_hex+", "+n2_hex+"))\nfmt.Println(uint64("+str(f2i(n1/n2))+"))\n")
    # eqf32
    outfile.write("fmt.Println(\"test the eqf32\")\n")
    outfile.write("fmt.Println(F32_eq("+n1_hex+", "+n2_hex+"))\nfmt.Println(\""+str(n1==n2)+"\")\n")
    # lef32
    outfile.write("fmt.Println(\"test the lef32\")\n")
    outfile.write("fmt.Println(F32_le("+n1_hex+", "+n2_hex+"))\nfmt.Println(\""+str(n1<=n2)+"\")\n")
    # ltf32
    outfile.write("fmt.Println(\"test the ltf32\")\n")
    outfile.write("fmt.Println(F32_lt("+n1_hex+", "+n2_hex+"))\nfmt.Println(\""+str(n1<n2)+"\")\n")
    # remf32
    outfile.write("fmt.Println(\"test the remf32\")\n")
    outfile.write("fmt.Println(F32_rem("+n1_hex+", "+n2_hex+"))\nfmt.Println(uint64("+str(f2i(n1%n2))+"))\n")
    # roundToIntf32
    outfile.write("fmt.Println(\"test the roundToIntf32\")\n")
    outfile.write("fmt.Println(F32_roundToInt("+n1_hex+",false))\nfmt.Println(int("+str(f2i(round(n1)))+"))\n")
    # sqrtf32
    outfile.write("fmt.Println(\"test the sqrtf32\")\n")
    newn1 = abs(n1)
    newn1_hex = f2h(newn1)
    outfile.write("fmt.Println(F32_sqrt("+newn1_hex+"))\nfmt.Println(uint64("+str(f2i(newn1**0.5))+"))\n")
    # mulAddf32
    outfile.write("fmt.Println(\"test the mulAddf32\")\n")
    outfile.write("fmt.Println(F32_mulAdd("+n1_hex+", "+n2_hex+", "+n1_hex+"))\nfmt.Println(uint64("+str(f2i(n1*n2+n1))+"))\n")

    # addf64
    outfile.write("fmt.Println(\"test the addf64\")\n")
    outfile.write("fmt.Println(F64_add("+d1_hex+", "+d2_hex+"))\nfmt.Println(uint64("+str(d2i(d1+d2))+"))\n")
    # subf64
    outfile.write("fmt.Println(\"test the subf64\")\n")
    outfile.write("fmt.Println(F64_sub("+d1_hex+", "+d2_hex+"))\nfmt.Println(uint64("+str(d2i(d1-d2))+"))\n")
    # mulf64
    outfile.write("fmt.Println(\"test the mulf64\")\n")
    outfile.write("fmt.Println(F64_mul("+d1_hex+", "+d2_hex+"))\nfmt.Println(uint64("+str(d2i(d1*d2))+"))\n")
    # divf64
    outfile.write("fmt.Println(\"test the divf64\")\n")
    outfile.write("fmt.Println(F64_div("+d1_hex+", "+d2_hex+"))\nfmt.Println(uint64("+str(d2i(d1/d2))+"))\n")
    # eqf64
    outfile.write("fmt.Println(\"test the eqf64\")\n")
    outfile.write("fmt.Println(F64_eq("+d1_hex+", "+d2_hex+"))\nfmt.Println(\""+str(d1==d2)+"\")\n")
    # lef64
    outfile.write("fmt.Println(\"test the lef64\")\n")
    outfile.write("fmt.Println(F64_le("+d1_hex+", "+d2_hex+"))\nfmt.Println(\""+str(d1<=d2)+"\")\n")
    # ltf64
    outfile.write("fmt.Println(\"test the ltf64\")\n")
    outfile.write("fmt.Println(F64_lt("+d1_hex+", "+d2_hex+"))\nfmt.Println(\""+str(d1<d2)+"\")\n")
    # remf64
    outfile.write("fmt.Println(\"test the remf64\")\n")
    outfile.write("fmt.Println(F64_rem("+d1_hex+", "+d2_hex+"))\nfmt.Println(uint64("+str(d2i(d1%d2))+"))\n")
    # roundToIntf64
    outfile.write("fmt.Println(\"test the roundToIntf64\")\n")
    outfile.write("fmt.Println(F64_roundToInt("+d1_hex+",false))\nfmt.Println(int("+str(d2i(round(d1)))+"))\n")
    # sqrtf64
    newd1 = abs(d1)
    newd1_hex = d2h(newd1)
    outfile.write("fmt.Println(\"test the sqrtf64\")\n")
    outfile.write("fmt.Println(F64_sqrt("+newd1_hex+"))\nfmt.Println(uint64("+str(d2i(newd1**0.5))+"))\n")
    # mulAddf64
    outfile.write("fmt.Println(\"test the mulAddf64\")\n")
    outfile.write("fmt.Println(F64_mulAdd("+d1_hex+", "+d2_hex+", "+d1_hex+"))\nfmt.Println(uint64("+str(d2i(d1*d2+d1))+"))\n")


    outfile.write("}\n")
    outfile.close()



for i in range(100):
    gen_test("test/"+str(i)+"_test.go", i)