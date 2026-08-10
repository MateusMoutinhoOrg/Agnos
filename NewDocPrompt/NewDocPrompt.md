
### Objective:
reorganize the doc to make informations more easy to find


### New Organization
The documentation will be divided in 3 main directories:

- CliUsage: for the command line interface usage
- LibUsage: for the library usage
- Development: for the development process and specifications
- Templating: for the templating system

each theme will be divided into References,Tutorials and Index
- Index: the main index of the theme 
  - [example](./Specs/Index/Sample.md)

- References: explanations about that subject
- Tutorials: workflows steps of how to make something.



### Tree Glossary:
 - CurrentDir: its the actual dir, that will need to be moved 
 - modifications.md: the instructions of what needs to be done
 
 


### Doc Tree:
- Readme.md 
  - [modifications](./README.md/modifications.md)
- docs/CliUsage/Index.md 

- docs/CliUsage/Tutorials/ 
- docs/CliUsage/Tutorials/QuickStart.md
- docs/CliUsage/References/
- docs/CliUsage/References/Commands.md



- docs/LibUsage/Index.md
- docs/LibUsage/Tutorials/
- docs/LibUsage/Tutorials/QuickStart.md
- docs/LibUsage/References/
- docs/LibUsage/References/PublicApi/
  - [CurrentDir](/docs/PublicApi/)
- docs/LibUsage/References/PublicApi.md
  - [CurrentDir](/docs/PublicApi.md)

- docs/Development/
- docs/Development/Index.md
- docs/Development/Tutorials/
- docs/Development/References/
- docs/Development/References/Specs/
  - [CurrentDir](/docs/Meta/)



- docs/Templating/
- docs/Templating/Index.md
- docs/Templating/Tutorials/
- docs/Templating/References/