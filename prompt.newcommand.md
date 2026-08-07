- the cli commands should be like:
### Command changes
the project needs to be called agnos-cli , not just agnos 


### Category managment

agnos-cli add category food 
agnos-cli remove category drinks 
agnos-cli list categories



### Transactions 

agnos spend food 10 
agnos recived salary 1000  

### Transactions with date 

agnos speend food 10 --date 06/07/2026 

### Listage 
agnos list -> list all 
agnos list --start 06/07/2026 -> list all after
agnos list --start 06/07/2026 --end 06/07/2028 
agnos list --start 06/07/2026 --end 06/07/2028 --category food --category drinks ## will list all elements betwen data, and that is food or drink