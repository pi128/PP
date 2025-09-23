import java.util.concurrent.Semaphore;

class SharedResource {
    private final Semaphore semaphoreLock;

    public SharedResource(int permits) {
        this.semaphoreLock = new Semaphore(permits);
    }

    public void accessResource(String threadName) {
        try {
            semaphoreLock.acquire();
            System.out.println(threadName + " is accessing the resource");
            Thread.sleep(500); // 1/2 sec
            System.out.println(threadName + " is done");
        } catch (InterruptedException e) {
            e.printStackTrace();
        } finally {
            semaphoreLock.release();
        }
    }
}

public class semaphore {
    public static void main(String[] args) {
        SharedResource resource = new SharedResource(1); // declare variable
        Runnable task = () -> {
            String threadN = Thread.currentThread().getName();
            resource.accessResource(threadN);
        };

        Thread t1 = new Thread(task, "Thread 1");
        Thread t2 = new Thread(task, "Thread 2");
        Thread t3 = new Thread(task, "Thread 3");

        t1.start();
        t2.start();
        t3.start();
    }
}
