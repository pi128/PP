use macroquad::prelude::*;
use path_core::{Grid, bfs};

#[macroquad::main("Maze Pathfinding")]
async fn main() {
    let grid = Grid::new(20, 15);
    let path = bfs(&grid, (0,0), (19,14));

    loop {
        clear_background(WHITE);

        
        // draw cells here later
        // draw path here later

        next_frame().await;
    }
}