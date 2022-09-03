    digit               = "0" … "9" .
    ascii_letter        = "A" … "Z" | "a" … "z" .
    letter              = ascii_letter | "_" .
    word                = ( letter ) { letter | digit } .
    time                = ( digit ) ":" ( digit ) ( "AM" | "PM" ) .
    
    
    Program           = { ProgramStatement | BlockStatement } .
    Block             = "{" { BlockStatement } "}"  | BlockStatement .
    ProgramStatement  = SceneStatement .
    BlockStatement    = SetStatement | GetStatement | VarStatement | AtStatement | WhenStatement 
                            StartStatement | StopStatement .
    SetStatement      = "set" PathMatch Value .
    VarStatement      = "var" word "=" GetStatement .
    GetStatement      = "get" Path .
    SceneStatement    = "scene" word Block .
    AtStatement       = "at" Time Action word .
    Time              = word | { digit } ":" { digit } ( "AM" | "PM" )
    Action            = ( "start" | "stop" )
    WhenStatement     = "when" PathMatch "is" Value "wait" duration Block  | "when" PathMatch "is" Value Block .
    Path              = "$" | { word "/" } word .
    PathMatch         = "$" | { ( word | "+" ) "/" } ( word | "+"  | "#" ) .
    StartStatement    = "start" word .
    StopStatement     = "stop" word .

  
---- 

    word                = ( letter ) { letter | digit } .
    ascii_letter        = "A" … "Z" | "a" … "z" .
    letter              = ascii_letter | "_" . 
    
    Program             = { RuleStatement | BlockStatement } .
    BlockStatement      = ActivateStatement | DectivateStatement | DoStatement |OrderedActionStatement | RandomActionStatement |  .
    RuleStatement       = "rule" "{" WhenStatement ThenStatement "}" as word .
    
    WhenStatement       = "when" (word | word "." word) ("<" | ">" | "==" | "!=")  (word | word "." word)|digit.
    ThenStatement       = "then" BlockStatement .
    
    
    
    ActivateStatement   = "activate" word  .
    DectivateStatement  = "deactivate" word  .
              
    RandomActionStatement   = "RandomAction" "["{word ","}"]" "as" word
    OrderedActionStatement  = "OrderedAction" "["{word ","}"]" "as" word
    DoStatement             = "Do" word "as" word.


---- 
    Arkham Grammer

    digit               = "0" … "9" .
    word                = ( letter ) { letter | digit } .
    ascii_letter        = "A" … "Z" | "a" … "z" .
    letter              = ascii_letter | "_" .
    string              = "'" { digit | letter } "'"

    ccodeStatement      = "ccode" word .
    Program             =   { BlockStatement } .
    BlockStatement      = damageStatement | testStatement | printStatement | ccodeStatement | emitStatement |   .
    onStatement         = "on" word "{" Program "}" .
    emitStatement       = "emit" word ":" {word} .
    testStatement       = "test" word "against" digit "{" successStatement failureStatement "}".
    successStatement    = success { Program } .
    failureStatement    = failure { Program } .
    printStatement      = "print" string .
    damageStatement     = "damage" word ":" number word .
    interceptStatement  = "intercept" word "{" Program "}" .
----
Example

    RandomAction [attack, defense, jump, move] as namedespreset
    OrderedAction [attack, defense, jump, move] as namedespreset 
    Do attack as namedespreset 

    activate namedespreset
    deactivate namedespreset
   
    rule {
        when thing is|isnot|less|more value
            then Do
    } as namedrule
    
    rule {
        when A is 12 and b less 3 and c not "yellow" 
            then Do dodge
    } as namedrule

    activate namedurle
    deactivate namedurle

# Actions without cost
- Wisper (log.debug, only player will see this)
- Say

#  Standard Actions that end the round
- attack
- defend
- move (+1m)
- backup (+1m)
- rush (move 2m)
- flee (move -2m)
- heal
- parry

# Skills - Actions that also end the round but not everyone has
- skill_jumpattack ( +2m and attack)
- skill_


#Things we e expose as data

## me
- currenthealth
- lastAction
- classification (Tank, DD, Weak, Fast )
- distance
- statueeffect

## enemy
- currenthealth
- lastAction
- classification (Tank, DD, Weak, Fast )
- distance
- statueeffect

## other
- round
   

