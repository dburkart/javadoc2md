# LinkTest

```java
public class LinkTest
```

This test ensures that @link tags work.

## `public void testLinkToOtherClass()` {#testLinkToOtherClass}

This is a method that links to [JavaClass](JavaClass) to make sure links
work.

## `public void testLinkToOtherMethodInClass()` {#testLinkToOtherMethodInClass}

We should be able to link to [testLinkToOtherClass](LinkTest#testLinkToOtherClass) from within
this class.

## `public void testLinksDontEatPeriods()` {#testLinksDontEatPeriods}

Links should not eat periods: [testLinkToOtherMethodInClass](LinkTest#testLinkToOtherMethodInClass). Did
the period disappear?

