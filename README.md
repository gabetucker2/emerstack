<h1>TODO</h1>

Completed:

* Defined bh before the conditionals (ctrl+f, "bh :=")
* Created new functions to facilitate ease of updating tensors (ctrl+f, "GABE CHANGES") (can be slightly optimized with some minor changes (map instead of IndexOf)) and implemented it with friends as an example for showcase and to prove no compile/rt errors

ToDo:

* REPLACE true/false BOOL WITH COMBINED ARRAY WHERE IF >7, THEN IT'S FALSE
* Create struct for inc/dec vals for environmental/interoceptive state, current behavior, and other enviro/intero cues
* Attempt to solve redundancy problem on own
* FIX GIT

To ask:

* Ask Dr. Read about his placement of the eating every frame... do I misunderstand?
* Clarify difference between training and trials- when should dynamics run?
- Git
  * How to get emergent-factory in the leabra@v1.2.0 folder since it not being a direct descendent of examples yields an error
- Argument-passing redundancy
  * How to make an instance, static return method in golang... had so much trouble with it
  * Method within a function yields errors... workaround?
- Data saving
  *  Before implementing WorldChanges.tsv... it looks like the World tsv file is not being updated, yet it is being accessed by the script... is this intentional?  They're both always 0
  * Why is WorldChanges.tsv commented out throughout the code?  Is it mostly, but not fully, implemented?

Sidebar:
  * Organize code layout with Andy
  * Create Go syntactic extensions
   - Ternary operator
   - IndexOf function
   - Ambiguously-defined variable access..?
   - Single-line lambda array return sort function like in js

Notes:
* Angle brackets have dimension of layer/tensor, first column of layer have brackets
* Need to decide dt
* defacto trial count
* $Name column is the real value
* Have major stochastic changes (e.g., friend suddenly leaves) - create log for major events