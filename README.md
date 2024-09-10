# goflow

`goflow`, a flowchart generator for Go.

# Usage

```
goflow <filename> <function name>
```

## Example

```
goflow testdata/src/example/example.go f5
```

```uml
@startuml
start
:f5;
if (arg == 0) then (yes)
:fmt.Printf("if %d\n", arg);
else
if (arg == 1) then (yes)
:fmt.Printf("else if %d\n", arg);
else
:fmt.Printf("else %d\n", arg);
endif
endif
switch (arg % 2)
case (1)
:odd = true;
case (0)
:even = true;
:return ;
end
case (default)
:err = fmt.Errorf("error: %d", arg);
:return ;
end
endswitch
switch ()
case (args == 10)
:err = fmt.Errorf("error: %d", arg);
:return ;
end
case (args == 20)
:err = fmt.Errorf("error: %d", arg);
:return ;
end
case ()
endswitch
while (for i := 0; i < 10; i++)
if (i == 5) then (yes)
:continue;
stop
endif
:return switchCase(i);
end
endwhile
:c int;
while (for)
:c++;
if (c > 10) then (yes)
:break;
break
endif
endwhile
:return ;
end
@enduml
```

![f5](/images/f5.png)


# TODO

- go routine
- channel
- select statement
- ... and more

# License

MIT
