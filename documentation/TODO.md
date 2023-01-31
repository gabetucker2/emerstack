<h1>Welcome to our TODO list</h1>

 <h2>Toward a first working implementation</h2>
 
 *In order from most important to least important.*
 
 * Find out how to define EnviroTsr and InteroTsr procedurally based on `Model.go` code
 * Have no "enviro" state updates (like in current Model implementation/diagram); define enviro state updates in Unity. For instance, upon eating food:
   * Intero state is updated in neural net model
   * Enviro-global state is updated in Unity => Enviro-local state updated in the neural net model
 * Remove all direct references of enviro and intero from `Model.go`; these should be called via some function that iterates over an array in which they're stored.
 * Finish replacing all mentions of layers, actions, etc with EmergentStack functions
 * Ensure the model is bug-free
 * Run the model, ensuring it produces realistic results
 
 <h2>Completed</h2>
 
 * Plan a schematic, git page, etc for the project
 * Refactored Model.go to use go map when defining parameters, rather than two arrays
 * Pluralized `b` and `dx` values
 * Mathematically modeled and implemented TimeIncrements
 * Remove cost function from model.  Perhaps we just need the state updates to intrinsically act as the cost (e.g., going on a walk increases need for food and sleep).  The neural network model hidden layers will (ideally) take care of representing "cost".
 * Added units to all variable names
 * Made layer setup explicit
 
 <h2>Future goals</h2>
 
 * *TODO: add*
 