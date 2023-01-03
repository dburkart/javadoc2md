/**
 * This test ensures that @link tags work.
 */
 public class LinkTest {

    /**
     * This is a method that links to {@link JavaClass} to make sure links
     * work.
     */
    public void testLinkToOtherClass() {

    }

    /**
     * We should be able to link to {@link #testLinkToOtherClass} from within
     * this class.
     */
    public void testLinkToOtherMethodInClass() {

    }

    /**
     * Links should not eat periods: {@link #testLinkToOtherMethodInClass}. Did
     * the period disappear?
     */
     public void testLinksDontEatPeriods() {

     }

    /**
     * Link to a class with a package name: {@link com.foo.bar.JavaClass#doSomething}
     */
     public void testLinksIncludingPackageNames() {

     }

     /**
      * Link to a class with arguments: {@link JavaClass#doSomething(long)}
      * Link to a class without arguments: {@link #testLinksDontEatPeriods()}
      * Link to a class with arguments and package name: {@link com.foo.bar.JavaClass#doSomething(long)}
      * Link to a class with multiple arguments: {@link FunctionDefOverSeveralLines#thisFunctionIsLongWinded(int, int, int)}
      * Link to an overloaded method: {@link #overloadedMethod}
      * Link to the other overloaded method: {@link #overloadedMethod(int, int)}
      * Link spanning multiple lines: {@link #overloadedMethod(int,
      * int)}
      * Link (plain): {@linkplain #overloadedMethod(int,int)}
      */
      public void testLinksWithArguments() {

      }

      /**
       * This method will be overloaded below.
       *
       * @param A this is an integer
       */
      public void overloadedMethod(int A) {}

      /**
       * This method overloads the one above.
       *
       * @param A this is an integer
       * @param B this is also an integer
       */
      public void overloadedMethod(int A, int B)
 }