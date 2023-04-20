
FSM gen

### Interface sketch:


schema.sql
```
create table mydata (
  id bigint not null autoincrement,

  state int64 not null,
  created_at datetime not null,
  updated_at datetime not null,

  value_one int,
  value_two varchar(255),

  primary key (id)
);
```

fsm_def.go
```
import (
  ".../fsmgen/fsm"
)

const (
  StateOne fsm.State = 1
  StateTwo fsm.State = 2
  StateThree fsm.State = 3
)


type Insert_StateOne struct {
  fsm.Vertex

  ValueOne int64
}

type StateOne_to_StateTwo struct {
  fsm.Vertex

  ID int64
  ValueTwo string
}

type StateOne_to_StateTwo_SettingValOne struct {
  fsm.Vertex

  ID int64
  ValueOne string
}

type StateTwo_to_StateThree struct {
  fsm.Vertex

  ID int64
}
```

Generated code:

fsm_gen.go
```

type ErrInvalidStateTransition ... // Some error type which captures invalid transitions - should include to+from state and the fsm name in the erro message

func TransitionFSM(v fsm.Vertex) error {

  // Check this vertext type does belong to this FSM somehow...
  // Update the db with the state change and the data change...
}

//... probably generate some receivers on each of the `fsm.Vertex`s as well
```

fsm_graph.md
```
// Generate a graph visulisation (UML probably) of the fsm structure
```


Calling code:
```
import "../db/mydata"

func DoWork() {


  id, err := mydata.TransitionFSM(
    mydata.Insert_StateOne{
      ValueOne: 123,
    },
  )
  // err handle...

  err := mydata.TransitionFSM(
    mydata.StateOne_to_StateTwo{
      ID: id,
      ValueTwo: "hello world",
    },
  )
  // err handle...

// *** Alternate interfafce:

  id, err := mydata.Insert_StateOne{
    ValueOne: 123,
  }.Exec()
  if err != nil {...}

  err := mydata.StateOne_to_StateTwo{
    ID: id,
    ValueTwo: "hello world",
  }.Exec()
  if err != nil {...}

}

```
