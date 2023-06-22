<h1>Welcome to our TODO list</h1>

 <h2>Toward a first working implementation</h2>
 
 *In order from most important to least important.*
 
 * INTEGRATE [data.md](data.md) CONTENTS INTO THIS TODO LIST 
 * Fix line 354 @@@@ nn.go
 * Finish replacing all mentions of layers, actions, etc with EmergentStack functions
 * Run the model, ensuring it produces realistic results
 * Remove all direct references of enviro and intero from `Model.go`; these should be called via some function that iterates over an array in which they're stored.
 * Have system for mapping actions/complex actions in NN to deliberate actions in Unity (complex actions will likely raise some issues with our current system)
 * Implement method for easily switching on and off Requirements condition for updates
 * Run the model with Unity, ensuring it produces realistic results
 
 <h2>Completed</h2>

 * Have no "enviro" state updates (like in current Model implementation/diagram); define enviro state updates in Unity. For instance, upon eating food:
   * Intero state is updated in neural net model
   * Enviro-global state is updated in Unity => Enviro-local state updated in the neural net model
   * BUT have it so non-Unity implementation updates Enviro manually, whereas Unity implementation ignores Enviro manual updates
 * Model Unity and NN message enviro equation 
 * Plan a schematic, git page, etc for the project
 * Refactored Model.go to use go map when defining parameters, rather than two arrays
 * Pluralized `b` and `dx` values
 * Mathematically modeled and implemented TimeIncrements
 * Remove cost function from model.  Perhaps we just need the state updates to intrinsically act as the cost (e.g., going on a walk increases need for food and sleep).  The neural network model hidden layers will (ideally) take care of representing "cost".
 * Added units to all variable names
 * Made layer setup explicit
 * Removed "current" tensors
 * Defined EnviroTsr and InteroTsr procedurally based on `Model.go` code
 * Created MathFunctions script, mathematiclaly modeled and implemented EnviroUnityToNN
 * Migrate to leabra v1.2.3, build initial model!

 <h2>Future goals</h2>
 
 * *TODO: add*
 