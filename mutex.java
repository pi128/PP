import java.util.concurrent.locks.Lock;
import java.util.concurrent.locks.ReentrantLock;

class Counter {
    private int count = 0;
    private Lock lock = new ReentrantLock();

    public void increment(){
        lock.lock();
        try{
            count++;
            System.out.println(Thread.currentThread().getName() + " incremented count to: " + count);
        } finally {
            lock.unlock();
        }
    }
}

public class mutex {
    public static void main(String[] args) {
        Counter counter = new Counter();

        Thread t1 = new Thread(() -> {
            for (int i = 0; i < 2; i++) counter.increment();

            }, "Thread 1");
        

        Thread t2= new Thread(() -> {
            for (int i = 0; i < 2; i++) counter.increment();

            }, "Thread 2");
        

        t1.start();
        t2.start();

        try{
            t1.join();
            t2.join();
        } catch(InterruptedException e){
            e.printStackTrace();
        }
    }
}