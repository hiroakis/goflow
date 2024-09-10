package goflow_test

import (
	"testing"

	"github.com/hiroakis/goflow"
)

func TestAnalyzeFunctionF1(t *testing.T) {
	err := goflow.AnalyzeFunction("testdata/src/example/example.go", "f1")
	if err != nil {
		t.Errorf("analyzeFunction failed: %v", err)
	}
	expect := `@startuml
start
:f1;
:fmt.Println("f1");
@enduml
`
	if goflow.GetUML() != expect {
		t.Errorf("got %v, want %v", goflow.GetUML(), expect)
	}
	t.Log(goflow.GetUML())
}

func TestAnalyzeFunctionF2(t *testing.T) {
	err := goflow.AnalyzeFunction("testdata/src/example/example.go", "f2")
	if err != nil {
		t.Errorf("analyzeFunction failed: %v", err)
	}
	expect := `@startuml
start
:f2;
:s := strconv.Itoa(arg);
if (s == "1") then (yes)
:err := f3(s);
if (err != nil) then (yes)
:return true, nil;
end
endif
else
:return true, fmt.Errorf("error: %d", arg);
end
endif
:return false, nil;
end
@enduml
`
	if goflow.GetUML() != expect {
		t.Errorf("got %v, want %v", goflow.GetUML(), expect)
	}
	t.Log(goflow.GetUML())
}

func TestAnalyzeFunctionF3(t *testing.T) {
	err := goflow.AnalyzeFunction("testdata/src/example/example.go", "f3")
	if err != nil {
		t.Errorf("analyzeFunction failed: %v", err)
	}
	expect := `@startuml
start
:f3;
:return nil;
end
@enduml
`
	if goflow.GetUML() != expect {
		t.Errorf("got %v, want %v", goflow.GetUML(), expect)
	}
	t.Log(goflow.GetUML())
}

func TestAnalyzeFunctionF4(t *testing.T) {
	err := goflow.AnalyzeFunction("testdata/src/example/example.go", "f4")
	if err != nil {
		t.Errorf("analyzeFunction failed: %v", err)
	}
	expect := `@startuml
start
:f4;
:return &t4{}, nil;
end
@enduml
`
	if goflow.GetUML() != expect {
		t.Errorf("got %v, want %v", goflow.GetUML(), expect)
	}
	t.Log(goflow.GetUML())
}

func TestAnalyzeFunctionF5(t *testing.T) {
	err := goflow.AnalyzeFunction("testdata/src/example/example.go", "f5")
	if err != nil {
		t.Errorf("analyzeFunction failed: %v", err)
	}
	expect := `@startuml
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
`
	if goflow.GetUML() != expect {
		t.Errorf("got %v, want %v", goflow.GetUML(), expect)
	}
	t.Log(goflow.GetUML())
}

func TestAnalyzeFunctionF6(t *testing.T) {
	err := goflow.AnalyzeFunction("testdata/src/example/example.go", "f6")
	if err != nil {
		t.Errorf("analyzeFunction failed: %v", err)
	}
	expect := `@startuml
start
:f6;
:f1();
:b, err := f2(1);
if (err != nil) then (yes)
:return err;
end
endif
if (b) then (yes)
:fmt.Println();
endif
:err := f3("str");
if (err != nil) then (yes)
:return err;
end
endif
:arg := &t4{};
:v, err := f4(arg);
if (err != nil) then (yes)
:return err;
end
endif
:fmt.Println(v);
while (range []int)
:fmt.Println(v);
endwhile
:d := []int;
while (range d)
:fmt.Println(i, v);
endwhile
:return nil;
end
@enduml
`
	if goflow.GetUML() != expect {
		t.Errorf("got %v, want %v", goflow.GetUML(), expect)
	}
	t.Log(goflow.GetUML())
}
