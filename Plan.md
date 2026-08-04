

###  Rename  (done)
- execute prompts/prompts/1.Rename.md 
~~~bash
claude -p "$(cat prompts/prompts/1.Rename.md)" --dangerously-skip-permissions
~~~    


### Aply the Factory mode 
- Execute the prompts/prompts/2.FactoryMode.md script
~~~bash
claude -p "$(cat prompts/prompts/2.FactoryMode.md)" --dangerously-skip-permissions
~~~


### adapt verbe pt.1
- clone Verbe (done)
- prepare the prompts/prompts/4.VerbAdapt.md prompt 
~~~bash
claude -p "$(cat prompts/prompts/4.VerbAdapt.md)" --dangerously-skip-permissions
~~~

### adapt keep pt.1
- clone keep (done)
- prepare the prompts/prompts/3.KeepAdapt.md prompt
~~~bash
claude -p "$(cat prompts/prompts/3.KeepAdapt.md)" --dangerously-skip-permissions
~~~

### create the main  mode
- prepare the required itens on RefTree
- prepare the prompts/prompts/5.mainMode.md 
- execute: prompts/prompts/5.mainMode.md
~~~bash
claude -p "$(cat prompts/prompts/5.mainMode.md)" --dangerously-skip-permissions
~~~
