# ebitengine-hexboard

a simple example for a hexboard with ebitengine.    
the example draws two floors and you can flip the tiles x and y.    

ebitengine: https://ebitengine.org/    
good information for hex: https://www.redblobgames.com/grids/hexagons/    

screenshot:    
![Pic1](screenshotsmall.jpg)

# new changes:     
- mousecoordinates and corresponding hextile
- selected hextile greenlighted
- first sprites
- recognition of mouseclicked sprite
- added astar pathfinding for hexboard ... program seems to turn into a jungle...maybe i should put the pathfinding later to other file/library?...    
    ```
    astar := NewAStar(ngrid)
    var path *Stack[*ANode]
	  path=astar.FindPath(Vector2{4,0}, Vector2{5,3})
	  fmt.Println("pathlen: ", path.Count())
    ```
    
  
