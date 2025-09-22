import threading
import random
import time
import math

try:
    import pygame
except Exception as e:
    pygame = None


class Dog:
    def __init__(self, name, world_rect, house_center, house_radius, speed_range=(60, 120)):
        self.name = name
        self.radius = 10
        # Spawn at house center with tiny jitter to avoid overlap
        jitter = 2.0
        self.x = house_center[0] + random.uniform(-jitter, jitter)
        self.y = house_center[1] + random.uniform(-jitter, jitter)
        self.radius = 10
        self.color = (200, 180, 60)
        self.world_rect = world_rect
        self.house_center = house_center
        self.house_radius = house_radius
        self.speed = random.uniform(*speed_range)  # pixels per second
        self.direction = random.uniform(0, math.tau)
        self.alive = True
        self.exploded_at = None
        self.lock = threading.Lock()
        self.spawn_time = time.perf_counter()
        self.left_house = False

    def get_name(self):
        return self.name

    def update_direction_randomly(self):
        turn = random.uniform(-math.pi / 4, math.pi / 4)
        self.direction = (self.direction + turn) % math.tau

    def move_step(self, dt):
        if not self.alive:
            return
        # Occasionally change direction
        if random.random() < 0.1:
            self.update_direction_randomly()
        # Clamp dt to avoid huge jumps if the system hitches
        dt = min(dt, 0.05)
        # Ease in speed for the first 0.4s after spawn
        t_since_spawn = time.perf_counter() - self.spawn_time
        ease = min(1.0, t_since_spawn / 0.4)
        eff_speed = self.speed * (0.3 + 0.7 * ease)
        dx = math.cos(self.direction) * eff_speed * dt
        dy = math.sin(self.direction) * eff_speed * dt
        with self.lock:
            self.x += dx
            self.y += dy
            # Track if the dog has left the house area
            if not self.left_house:
                if ((self.x - self.house_center[0]) ** 2 + (self.y - self.house_center[1]) ** 2) ** 0.5 > (self.house_radius + self.radius):
                    self.left_house = True
            # Explode when leaving the world bounds (after grace and after leaving house), using margin = radius
            bounds_for_center = self.world_rect.inflate(-self.radius * 2, -self.radius * 2)
            if t_since_spawn > 0.35 and self.left_house and not bounds_for_center.collidepoint(self.x, self.y):
                self.alive = False
                self.exploded_at = (self.x, self.y)

    def roam(self, stop_event):
        prev = time.perf_counter()
        while not stop_event.is_set() and self.alive:
            now = time.perf_counter()
            dt = now - prev
            prev = now
            self.move_step(dt)
            time.sleep(0.01)


def run_game():
    if pygame is None:
        print("pygame is not installed. Activate your venv and run: pip install pygame")
        return

    pygame.init()
    width, height = 960, 640
    screen = pygame.display.set_mode((width, height))
    pygame.display.set_caption("Roaming Dogs - Out of Bounds Boom")
    clock = pygame.time.Clock()

    world_rect = pygame.Rect(40, 40, width - 80, height - 80)
    # House in the middle
    house_center = (world_rect.centerx, world_rect.centery)
    house_radius = 28

    font = pygame.font.SysFont(None, 18)

    try:
        dog_count = int(input("how many dogs do you want? "))
    except Exception:
        dog_count = 5

    dogs = []
    threads = []
    stop_event = threading.Event()

    for i in range(dog_count):
        name = input(f"Name {i+1}: ") or f"Dog{i+1}"
        dog = Dog(name, world_rect, house_center, house_radius)
        dogs.append(dog)
        t = threading.Thread(target=dog.roam, name=f"{dog.get_name()} thread", args=(stop_event,), daemon=True)
        threads.append(t)
        t.start()

    explosions = []  # list of dicts: {pos: (x,y), start: time, duration: s}

    running = True
    while running:
        for event in pygame.event.get():
            if event.type == pygame.QUIT:
                running = False
        # Collect explosions for newly dead dogs
        for dog in dogs:
            if not dog.alive and dog.exploded_at is not None:
                explosions.append({
                    "pos": dog.exploded_at,
                    "start": time.perf_counter(),
                    "duration": 0.6
                })
                dog.exploded_at = None  # ensure added once

        # Draw
        screen.fill((25, 28, 34))

        # World bounds
        pygame.draw.rect(screen, (60, 70, 90), world_rect, width=2)
        # House (center circle)
        pygame.draw.circle(screen, (90, 160, 200), house_center, house_radius)
        pygame.draw.circle(screen, (30, 60, 100), house_center, house_radius, width=2)

        # Dogs
        alive_count = 0
        for dog in dogs:
            if dog.alive:
                alive_count += 1
                with dog.lock:
                    pos = (int(dog.x), int(dog.y))
                pygame.draw.circle(screen, dog.color, pos, dog.radius)
                label = font.render(dog.name, True, (220, 220, 220))
                screen.blit(label, (pos[0] - label.get_width() // 2, pos[1] - dog.radius - 16))

        # Explosions
        now = time.perf_counter()
        remaining_explosions = []
        for ex in explosions:
            tnorm = (now - ex["start"]) / ex["duration"]
            if tnorm < 1.0:
                radius = int(10 + 40 * tnorm)
                alpha = int(255 * (1 - tnorm))
                draw_explosion(screen, ex["pos"], radius, alpha)
                remaining_explosions.append(ex)
        explosions = remaining_explosions

        # HUD
        hud = font.render(f"Alive: {alive_count}/{len(dogs)}  -  House in center", True, (230, 230, 230))
        screen.blit(hud, (12, 12))

        pygame.display.flip()
        clock.tick(60)

        # End when all dogs are gone
        if alive_count == 0 and len(explosions) == 0:
            running = False

    stop_event.set()
    # No need to join daemon threads on quit, but we can wait briefly
    for t in threads:
        t.join(timeout=0.2)

    pygame.quit()
    print("all dogs sleeping peacefully...")


def draw_explosion(screen, pos, radius, alpha):
    # Create a temporary surface for alpha blending
    surf = pygame.Surface((radius * 2 + 4, radius * 2 + 4), pygame.SRCALPHA)
    center = (radius + 2, radius + 2)
    # Outer ring
    pygame.draw.circle(surf, (255, 120, 40, max(0, alpha // 2)), center, radius)
    # Inner core
    pygame.draw.circle(surf, (255, 220, 60, alpha), center, max(2, radius // 3))
    # Spikes
    for i in range(8):
        angle = i * (math.tau / 8)
        x = int(center[0] + math.cos(angle) * radius)
        y = int(center[1] + math.sin(angle) * radius)
        pygame.draw.circle(surf, (255, 180, 80, alpha // 3), (x, y), max(2, radius // 6))
    screen.blit(surf, (int(pos[0] - center[0]), int(pos[1] - center[1])))


if __name__ == "__main__":
    run_game()