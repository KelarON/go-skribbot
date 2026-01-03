# go-scribbot (beta version)

Simple auto-drawing tool for skribbl.io written in Go.

## Instructions

### Setup

- On first run, a config file named config.yaml will be created with default settings.
- Make sure print_coords_mode is set to true.
- Start a game on skribbl.io (preferably a private room with two browser tabs).
- When it is your turn to draw, run the program and hover over the white color on the palette.
- After 5 seconds, the program will print coordinates (Your coordinates X: position_x, Y: position_y). Copy these values into the config file.
- After entering the correct coordinates, set print_coords_mode to false.

### Drawing

- Place any image file in the same directory as the program.
- When it is your turn to draw, run the program and open the game in your browser.
- If you need to stop drawing, quickly move your mouse cursor to the left.
