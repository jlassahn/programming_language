
# Data Types

FIXME to be written...


## Numeric Types

Signed integer types:
* Int64
* Int32
* Int16
* Int8
If an integer is specified without an explicit type it is assumed to be Int64.


Unsigned integer types:
* UInt64
* UInt32
* UInt16
* UInt8

Floating point (real number) types:
* Real32
* Real64
If a decimal number is specified without explicit type it is assumend to
be Real64.

Numeric constants can be written with tags at the end that control the type:
```
def x Int8 = 93s8;
def y Int64 = 93s64;
```

Some cases where an integer constant might need to be outside the bounds
of the type they're declared:
```
def x = -128s8;           # 128s8 isn't valid, but -128s8 is
def x = -100s8 + 200s8;   # 200s8 isn't vaid, but -100 + 200 is
```
A constant can be out of range for it's type, but if the result of an
operation or assignment is out of range it will be truncated and generate
a warning.

