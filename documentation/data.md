* World will always have concrete values in the first time, each subs. line will be based on prev, only thing that matters is first line of World bc that's starting state, basis of future rows
* WorldChanges is exogenous preprogrammed changes
* wanttobeh is obsolete
* world is testing data, everything else is training

// TODO: finish later
// temporarily removed worldChanges until we figure out how to program exogenous changes

make sure trained.wts goes to data rather than directory head

"init" crashes, ignore for the time being since windows rendering error
"dynamics" doesnts
"train" doesnt
"train PIT" doesnt

it knows pvlv then instr vs instr then pvlv based on GUI encodement, TRN etable.Table allows you to update it

We can remove one of the two and just reencode based on what we want in theory by re selecting the file in Trn USING: InstrThenPvlv

r.Wt or s.Wt used to test (receive/sending) what current val is, it may be that no starting value was initialized and therefore it isnt updated properly

s shows send to, r shows receive from
"Act" should show stuff

line @@@@@@ is crashing in NN.go

a couple lines of code, RecordSyns, there's a function

not updating the gui properly, init render not working, @@@@@@

* possible vulkan newer version not playing well with code
* have to possibly use old go.mod version

in the meantime, get the world stuff working
