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
 }