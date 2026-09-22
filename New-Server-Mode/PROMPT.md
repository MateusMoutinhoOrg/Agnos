

## new Routes Mechanic 
using  [sampples](New-Server-Mode/samples) refactor the route system.



### Priority Mechanic:
the prioriy system workss from top to button, if a  route has priority 0 and other has priority 5, it will run first priority 5 if it not made a write on body or not returned a err, it will run priority 0 until it makes a write on body or return a err

