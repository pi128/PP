import threading
import random
import time

class Dog:
    def __init__(self, name, bound):
        self.name = name
        self.bound = bound

    def get_name(self):
        return self.name
  

    def get_bound(self):
        return self._bound
    
    def set_bound(self, bound):
        self._bound = int(bound)

    def roam(self):
        while True:

            roam = random.randint(0, int(self.bound * 1.5))

            if roam > self.bound:
                print(f"[{threading.current_thread().name}] {self.name} exploded at {roam}!")
                #print("[{}] {} exploded at {}!".format(threading.current_thread().name, self.name, roam))

                break
            else:
                #print("[{}] {} chilling at {}!".format(threading.current_thread().name, self.name, roam))
                print(f"[{threading.current_thread().name}] {self.name} chilling at {roam}!")

            
            time.sleep(random.uniform(0.5, 2.0))



if __name__ == "__main__":

    dogs = []
    threads = []

    dog_count = int(input("how many dogs do you want? "))

    for i in range(dog_count):   
        name = input(f"Name: {i+1}: ")
        bound = int(input(f"Bound of {name}: "))

        dog = Dog(name, bound)
        dogs.append(dog)

        t = threading.Thread(target=dog.roam, name=f"{dog.get_name()} thread")
        threads.append(t)
    
    for t in threads:
        t.start()

    for t in threads:
        t.join()    

    print("all dogs sleeping peacfully...")