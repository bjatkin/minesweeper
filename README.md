# TODO
## features
  - [ ] refactor/ clean up the code so it's easier to work with
  - [ ] make transitions from game to world map a little less jerky (maybe a wipe?, or level itesm drow into the scean?)
  - [ ] add a full screen mode (look into the layout options that ebiten provides)
  - [ ] add a save game functionality
  - [ ] add some load game functionality
  - [ ] limit the number of a specific power up that you can bring into a map
  - [x] play an animation showing when you unlock powerups
  - [x] make power up's unlockable
  - [x] hide powerups on the menu screen if they haven't been unlocked yet
  - [x] Add a continue button to the pause menu (the play pause buttons in the corner are too small)
  - [x] test non-square board layouts (start with the heart)
  - [x] create a converter from png images to board layouts
  - [x] make levels with the count down timer possible
  - [x] make it so you can bring less than 3 powerups into battle

## graphics:
  - [x] mock a save game ui/ flow
  - [x] mock a load game ui/ flow
  - [x] mock a new pause menu (with added resume/ continue button)
  - [x] update ui button colors to be easier to read/ scan
  - [x] mock a power up unlock screen

## bugs
  - [ ] best times for levels are not being saved/ loaded correctly
  - [ ] when starting map you can acidently end up freezing the 'safe' tile, preventing any files from being flipped
  - [ ] add a way to eat inputs to prevent clicking through the UI or clicking on menues right when they pop up
  - [x] powerups look fully charged even when they are not. The last row of pixels should not be drawn in until the power up is 100% charged or it feels bad
  - [x] when you use the ESC key to exit the pause menu it dosen't work if you entered the menu using the mouse
  - [x] cat / tidal wave power ups are not getting re-uped after rechargeing?? (verify this)
  - [x] shuffel seems to still be a little bit buggy?? (verify this)