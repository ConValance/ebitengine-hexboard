# ebitengine-hexboard

a simple example for a hexboard with ebitengine (simple cause for example only few biomes for floor0 and only some parts for floor1).    
the example draws two floors and you can flip the tiles x and y.  the hex-pathfinding works (see green path over the bridge).    
the pathfinding is now in astarhexlib local library. for working import the name of maindirectory should be "hexboard"      
or change the name of the directory in main.og in import from "hexboard/astarhexlib" to "yourdirectory/astarhexlib"        

ebitengine: https://ebitengine.org/    
good information for hex: https://www.redblobgames.com/grids/hexagons/    

screenshot:    
![Pic1](screenshotsmall.jpg)

# new changes:     
- added astar pathfinding for hexboard ... program seems to turn into a jungle...maybe i should put the pathfinding later to other file/library?...
  from the test in init:        
    ```
    astar := NewAStar(ngrid)
    var path *Stack[*ANode]
     path=astar.FindPath(Vector2{4,0}, Vector2{5,3})
     fmt.Println("pathlen: ", path.Count())
     for i := 0; i < path.Count(); i++ {
		fmt.Println("path Nr:", i, " x:", path.items[i].Position.X, " y:", path.items[i].Position.Y)
     }
    ```
- changed layout to oddq. now hexcoordinates are correct and hexpathfinding works (see green path over the bridge)!!
- mouseclick left for new startpos, mouseclick right for new target
- check if path found, new part for floor1 a tree, floor as parameter
- changed the pathfinding to a local library in file astarhexlib.go    
  
  
    
  
