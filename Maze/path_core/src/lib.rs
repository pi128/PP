pub struct Grid {
    pub width: usize,
    pub height: usize,
    pub cells: Vec<u8>, // 0 = open, 1 = wall
}

impl Grid {
    pub fn new(width: usize, height: usize) -> Self {
        Self {
            width,
            height,
            cells: vec![0; width * height],
        }
    }

    pub fn index(&self, x: usize, y: usize) -> usize {
        y * self.width + x
    }
}

// placeholder for BFS/Dijkstra/A*
pub fn bfs(grid: &Grid, start: (usize, usize), goal: (usize, usize)) -> Vec<(usize, usize)> {
    Vec::new() // implement later
}