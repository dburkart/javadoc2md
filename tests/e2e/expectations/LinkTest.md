# LinkTest

## Definition

```java
public class LinkTest
```

## Overview

This test ensures that @link tags work.

### `public void testLinkToOtherClass()` {#testLinkToOtherClass()}

This is a method that links to [JavaClass](JavaClass) to make sure links
work.

### `public void testLinkToOtherMethodInClass()` {#testLinkToOtherMethodInClass()}

We should be able to link to [testLinkToOtherClass](LinkTest#testLinkToOtherClass()) from within
this class.

### `public void testLinksDontEatPeriods()` {#testLinksDontEatPeriods()}

Links should not eat periods: [testLinkToOtherMethodInClass](LinkTest#testLinkToOtherMethodInClass()). Did
the period disappear?

### `public void testLinksIncludingPackageNames()` {#testLinksIncludingPackageNames()}

Link to a class with a package name: [doSomething](JavaClass#doSomething(long))

### `public void testLinksWithArguments()` {#testLinksWithArguments()}

Link to a class with arguments: [doSomething](JavaClass#doSomething(long))
Link to a class without arguments: [testLinksDontEatPeriods](LinkTest#testLinksDontEatPeriods())
Link to a class with arguments and package name: [doSomething](JavaClass#doSomething(long))
Link to a class with multiple arguments: [thisFunctionIsLongWinded](FunctionDefOverSeveralLines#thisFunctionIsLongWinded(int,int,int))
Link to an overloaded method: [overloadedMethod](LinkTest#overloadedMethod(int))
Link to the other overloaded method: [overloadedMethod](LinkTest#overloadedMethod(int,int))
Link spanning multiple lines: [overloadedMethod](LinkTest#overloadedMethod(int,int))

### `public void overloadedMethod(int A)` {#overloadedMethod(int)}

This method will be overloaded below.

**Parameters:**

* `A` - this is an integer

### `public void overloadedMethod(int A, int B)` {#overloadedMethod(int,int)}

This method overloads the one above.

**Parameters:**

* `A` - this is an integer
* `B` - this is also an integer

