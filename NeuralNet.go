// Copyright (c) 2019, The Emergent Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Personality Model
// Currently set up to do separate training of Pavlovian and Instrumental Training.
// Interaction with Environment is still a work in progress

package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/emer/emergent/emer"
	"github.com/emer/emergent/env"
	"github.com/emer/emergent/netview"
	"github.com/emer/emergent/params"
	"github.com/gabetucker2/gogenerics"
	. "github.com/gabetucker2/gostack"

	// "github.com/emer/emergent/patgen"
	"github.com/emer/emergent/prjn"
	"github.com/emer/emergent/relpos"
	"github.com/emer/etable/agg"
	"github.com/emer/etable/eplot"
	"github.com/emer/etable/etable"
	"github.com/emer/etable/etensor"
	_ "github.com/emer/etable/etview" // include to get gui views
	"github.com/emer/etable/split"
	"github.com/emer/leabra/leabra"
	"github.com/goki/gi/gi"
	"github.com/goki/gi/gimain"
	"github.com/goki/gi/giv"
	"github.com/goki/ki/ki"
	"github.com/goki/ki/kit"
	"github.com/goki/mat32"
)

func main() {
	TheSim.New()
	TheSim.Config()
	if len(os.Args) > 1 {
		TheSim.CmdArgs() // simple assumption is that any args = no gui -- could add explicit arg if you want
		} else {
			gimain.Main( func() { 
				
				guirun()
				
			})
		}
	}
	
	func guirun() {
		TheSim.Init()
		win := TheSim.ConfigGui()
		win.StartEventLoop()
	}
	
	// LogPrec is precision for saving float values in logs
	const LogPrec = 4
	
	// ParamSets is the default set of parameters -- Base is always applied, and others can be optionally
	// selected to apply on top of that
	var ParamSets = params.Sets{
		{Name: "Base", Desc: "these are the best params", Sheets: params.Sheets{
			"Network": &params.Sheet{
				{Sel: "Prjn", Desc: "norm and momentum on works better, but wt bal is not better for smaller nets",
				Params: params.Params{
					"Prjn.Learn.Norm.On":     "true",
					"Prjn.Learn.Momentum.On": "true",
					"Prjn.Learn.WtBal.On":    "false",
					"Prjn.Learn.Learn": "true",
					"Prjn.Learn.Lrate":  "0.04",
					"Prjn.WtInit.Dist": "Uniform",
					"Prjn.WtInit.Mean": "0.5",
					"Prjn.WtInit.Var":  "0.25",
					"Prjn.WtScale.Abs": "1",
					
				}},
				{Sel: "#MotiveToApproach", Desc: "Set weight to fixed value of 1",
				Params: params.Params{
					"Prjn.Learn.Learn": "false",
					"Prjn.Learn.Lrate":  "0",
					"Prjn.WtInit.Dist": "Uniform",
					"Prjn.WtInit.Mean": "1",
					"Prjn.WtInit.Var":  "0",
					"Prjn.WtScale.Abs": "1",		
				}},
				{Sel: "#MotiveToAvoid", Desc: "Set weight to fixed value of 1",
				Params: params.Params{
					"Prjn.Learn.Learn": "false",
					"Prjn.Learn.Lrate":  "0",
					"Prjn.WtInit.Dist": "Uniform",
					"Prjn.WtInit.Mean": "1",
					"Prjn.WtInit.Var":  "0",
					"Prjn.WtScale.Abs": "1",		
				}},
				{Sel: "#VTA_DAToApproach", Desc: "Set weight to fixed value of 1, and weight scale of 2",
				Params: params.Params{
					"Prjn.Learn.Learn": "false",
					"Prjn.Learn.Lrate":  "0",
					"Prjn.WtInit.Dist": "Uniform",
					"Prjn.WtInit.Mean": "1",
					"Prjn.WtInit.Var":  "0",
					"Prjn.WtScale.Abs": "2",		
				}},
				{Sel: "#VTA_DAToAvoid", Desc: "Set weight to fixed value of 1, and weight scale of 2, inhibitory connection",
				Params: params.Params{
					"Prjn.Learn.Learn": "false",
					"Prjn.Learn.Lrate":  "0",
					"Prjn.WtInit.Dist": "Uniform",
					"Prjn.WtInit.Mean": "1",
					"Prjn.WtInit.Var":  "0",
					"Prjn.WtScale.Abs": "2",	
				}},
				{Sel: "Layer", Desc: "using 2.3 inhib for all of network -- can explore",
				Params: params.Params{
					"Layer.Inhib.Layer.Gi": "2.3",
					"Layer.Act.Gbar.L":     "0.1", // set explictly, new default, a bit better vs 0.2
					"Layer.Act.XX1.Gain" : "100",
				}},
				{Sel: "#Behavior", Desc: "Make Behavior layer selective for 1 behavior",
				Params: params.Params{
					"Layer.Inhib.Layer.Gi": "2.5",
					"Layer.Act.XX1.Gain": "200",
					}},	
					{Sel: "#Approach", Desc: "",
					Params: params.Params{
						"Layer.Inhib.Layer.Gi": "2.3",
						"Layer.Act.XX1.Gain": "400",
						}},	
						{Sel: "#Avoid", Desc: "",
						Params: params.Params{
							"Layer.Inhib.Layer.Gi": "2.3",
							"Layer.Inhib.Layer.FB": "1.2",
							"Layer.Act.XX1.Thr": ".49",
							"Layer.Act.XX1.Gain": "400",
							"Layer.Act.Gbar.I": "1.35",
							}},	
							{Sel: "#Hidden", Desc: "Make Hidden representation a bit sparser",
							Params: params.Params{
								"Layer.Inhib.Layer.Gi": "2.5",
								}},	
								{Sel: ".Back", Desc: "top-down back-projections MUST have lower relative weight scale, otherwise network hallucinates",
								Params: params.Params{
									"Prjn.WtScale.Rel": "0.2",
								}},
								
							},
							"Sim": &params.Sheet{ // sim params apply to sim object
								{Sel: "Sim", Desc: "best params always finish in this time",
								Params: params.Params{
									"Sim.MaxEpcs": "100",
								}},
							},
						}},
						{Name: "DefaultInhib", Desc: "output uses default inhib instead of lower", Sheets: params.Sheets{
							"Network": &params.Sheet{
								{Sel: "#Behavior", Desc: "go back to default",
								Params: params.Params{
									"Layer.Inhib.Layer.Gi": "1.8",
								}},
							},
							"Sim": &params.Sheet{ // sim params apply to sim object
								{Sel: "Sim", Desc: "takes longer -- generally doesn't finish..",
								Params: params.Params{
									"Sim.MaxEpcs": "100",
								}},
							},
						}},
						{Name: "NoMomentum", Desc: "no momentum or normalization", Sheets: params.Sheets{
							"Network": &params.Sheet{
								{Sel: "Prjn", Desc: "no norm or momentum",
								Params: params.Params{
									"Prjn.Learn.Norm.On":     "false",
									"Prjn.Learn.Momentum.On": "false",
								}},
							},
						}},
						{Name: "WtBalOn", Desc: "try with weight bal on", Sheets: params.Sheets{
							"Network": &params.Sheet{
								{Sel: "Prjn", Desc: "weight bal on",
								Params: params.Params{
									"Prjn.Learn.WtBal.On": "true",
								}},
							},
						}},
					}
					
					// Sim encapsulates the entire simulation model, and we define all the
					// functionality as methods on this struct.  This structure keeps all relevant
					// state information organized and available without having to pass everything around
					// as arguments to methods, and provides the core GUI interface (note the view tags
					// for the fields which provide hints to how things should be displayed).
					type Sim struct {
						Net          *leabra.Network   `view:"no-inline" desc:"the network -- click to view / edit parameters for layers, prjns, etc"`
						Instr        *etable.Table     `view:"no-inline" desc:"Training pattern for Instrumental Learning"`
						Pvlv         *etable.Table     `view:"no-inline" desc:"Training pattern for Pavlovian Learning"`
						Trn    		 *etable.Table     `view:"no-inline" desc:"Table that controls type of training and number of Epochs of training"`
						Training	 string			   `view:"no-inline" desc:"Type of training: Pavlovian or Instrumental"`
						World 	 	 *etable.Table     `view:"no-inline" desc:"Table that represents starting state and then each new state of the Internal and External world"`
						// WorldChanges *etable.Table     `view:"no-inline" desc:"Table that represents change in External world at each time step"`
						TrnEpcLog    *etable.Table     `view:"no-inline" desc:"training epoch-level log data"`
						TstEpcLog    *etable.Table     `view:"no-inline" desc:"testing epoch-level log data"`
						TstTrlLog    *etable.Table     `view:"no-inline" desc:"testing trial-level log data"`
						TstErrLog    *etable.Table     `view:"no-inline" desc:"log of all test trials where errors were made"`
						TstErrStats  *etable.Table     `view:"no-inline" desc:"stats on test trials where errors were made"`
						TstCycLog    *etable.Table     `view:"no-inline" desc:"testing cycle-level log data"`
						RunLog       *etable.Table     `view:"no-inline" desc:"summary log of each run"`
						RunStats     *etable.Table     `view:"no-inline" desc:"aggregate stats on all runs"`
						Params       params.Sets       `view:"no-inline" desc:"full collection of param sets"`
						ParamSet     string            `desc:"which set of *additional* parameters to use -- always applies Base and optionaly this next if set"`
						Tag          string            `desc:"extra tag string to add to any file names output from sim (e.g., weights files, log files, params for run)"`
						MaxRuns      int               `desc:"maximum number of model runs to perform"`
						MaxEpcs      int               `desc:"maximum number of epochs to run per model run"`
						NZeroStop    int               `desc:"if a positive number, training will stop after this many epochs with zero SSE"`
						TrainEnv     env.FixedTable    `desc:"Training environment -- contains everything about iterating over input / output patterns over training"`
						TestEnv      env.FixedTable    `desc:"Testing environment -- manages iterating over testing"`
						Time         leabra.Time       `desc:"leabra timing parameters and state"`
						ViewOn       bool              `desc:"whether to update the network view while running"`
						TrainUpdt    leabra.TimeScales `desc:"at what time scale to update the display during training?  Anything longer than Epoch updates at Epoch in this model"`
						TestUpdt     leabra.TimeScales `desc:"at what time scale to update the display during testing?  Anything longer than Epoch updates at Epoch in this model"`
						TestInterval int               `desc:"how often to run through all the test patterns, in terms of training epochs -- can use 0 or -1 for no testing"`
						LayStatNms   []string          `desc:"names of layers to collect more detailed stats on (avg act, etc)"`
						
						// statistics: note use float64 as that is best for etable.Table
						TrlErr        float64 `inactive:"+" desc:"1 if trial was error, 0 if correct -- based on SSE = 0 (subject to .5 unit-wise tolerance)"`
						TrlSSE        float64 `inactive:"+" desc:"current trial's sum squared error"`
						TrlAvgSSE     float64 `inactive:"+" desc:"current trial's average sum squared error"`
						TrlCosDiff    float64 `inactive:"+" desc:"current trial's cosine difference"`
						EpcSSE        float64 `inactive:"+" desc:"last epoch's total sum squared error"`
						EpcAvgSSE     float64 `inactive:"+" desc:"last epoch's average sum squared error (average over trials, and over units within layer)"`
						EpcPctErr     float64 `inactive:"+" desc:"last epoch's average TrlErr"`
						EpcPctCor     float64 `inactive:"+" desc:"1 - last epoch's average TrlErr"`
						EpcCosDiff    float64 `inactive:"+" desc:"last epoch's average cosine difference for output layer (a normalized error measure, maximum of 1 when the minus phase exactly matches the plus)"`
						EpcPerTrlMSec float64 `inactive:"+" desc:"how long did the epoch take per trial in wall-clock milliseconds"`
						FirstZero     int     `inactive:"+" desc:"epoch at when SSE first went to zero"`
						NZero         int     `inactive:"+" desc:"number of epochs in a row with zero SSE"`
						
						// internal state - view:"-"
						SumErr       float64                     `view:"-" inactive:"+" desc:"sum to increment as we go through epoch"`
						SumSSE       float64                     `view:"-" inactive:"+" desc:"sum to increment as we go through epoch"`
						SumAvgSSE    float64                     `view:"-" inactive:"+" desc:"sum to increment as we go through epoch"`
						SumCosDiff   float64                     `view:"-" inactive:"+" desc:"sum to increment as we go through epoch"`
						Win          *gi.Window                  `view:"-" desc:"main GUI window"`
						NetView      *netview.NetView            `view:"-" desc:"the network viewer"`
						ToolBar      *gi.ToolBar                 `view:"-" desc:"the master toolbar"`
						TrnEpcPlot   *eplot.Plot2D               `view:"-" desc:"the training epoch plot"`
						TstEpcPlot   *eplot.Plot2D               `view:"-" desc:"the testing epoch plot"`
						TstTrlPlot   *eplot.Plot2D               `view:"-" desc:"the test-trial plot"`
						TstCycPlot   *eplot.Plot2D               `view:"-" desc:"the test-cycle plot"`
						RunPlot      *eplot.Plot2D               `view:"-" desc:"the run plot"`
						TrnEpcFile   *os.File                    `view:"-" desc:"log file"`
						RunFile      *os.File                    `view:"-" desc:"log file"`
						ValsTsrs     map[string]*etensor.Float32 `view:"-" desc:"map tensor for holding layer values"`
						tsrsStack    *Stack                      `view:"-" desc:"gostack map structure containing the keys and values for each stack"`
						SaveWts      bool                        `view:"-" desc:"for command-line run only, auto-save final weights after each run"`
						NoGui        bool                        `view:"-" desc:"if true, runing in no GUI mode"`
						LogSetParams bool                        `view:"-" desc:"if true, print message for all params that are set"`
						IsRunning    bool                        `view:"-" desc:"true if sim is running"`
						StopNow      bool                        `view:"-" desc:"flag to stop running"`
						NeedsNewRun  bool                        `view:"-" desc:"flag to initialize NewRun if last one finished"`
						RndSeed      int64                       `view:"-" desc:"the current random seed"`
						LastEpcTime  time.Time                   `view:"-" desc:"timer for last epoch"`
					}
						
						// this registers this Sim Type and gives it properties that e.g.,
						// prompt for filename for save methods.
						var KiT_Sim = kit.Types.AddType(&Sim{}, SimProps)
						
						// TheSim is the overall state for this simulation
						var TheSim Sim
						
						// New creates new blank elements and initializes defaults
						func (ss *Sim) New() {
							ss.Net = &leabra.Network{}
							ss.Instr = &etable.Table{}
							ss.Pvlv = &etable.Table{}
							ss.Trn = &etable.Table{}
							ss.World = &etable.Table{}
							// ss.WorldChanges = &etable.Table{}
							ss.TrnEpcLog = &etable.Table{}
							ss.TstEpcLog = &etable.Table{}
							ss.TstTrlLog = &etable.Table{}
							ss.TstCycLog = &etable.Table{}
							ss.RunLog = &etable.Table{}
							ss.RunStats = &etable.Table{}
							ss.Params = ParamSets
							ss.RndSeed = 1
							ss.ViewOn = true
							ss.TrainUpdt = leabra.AlphaCycle
							ss.TestUpdt = leabra.Cycle
							ss.TestInterval = -1
							ss.LayStatNms = []string{"Approach", "Avoid", "Hidden", "Behavior"}
						}
						
						////////////////////////////////////////////////////////////////////////////////////////////
						// 		Configs
						
						// Config configures all the elements using the standard functions
						func (ss *Sim) Config() {
							//ss.ConfigPats()
							ss.OpenPats()
							ss.ConfigEnv()
							ss.ConfigNet(ss.Net)
							ss.ConfigTrnEpcLog(ss.TrnEpcLog)
							ss.ConfigTstEpcLog(ss.TstEpcLog)
							ss.ConfigTstTrlLog(ss.TstTrlLog)
							ss.ConfigTstCycLog(ss.TstCycLog)
							ss.ConfigRunLog(ss.RunLog)
						}
						// are all of these tensors used/needed??

							func (ss *Sim) ConfigEnv() {
								if ss.MaxRuns == 0 { // allow user override
									ss.MaxRuns = 10
								}
								if ss.MaxEpcs == 0 { // allow user override
									ss.MaxEpcs = 50
									ss.NZeroStop = 5
								}
								
								ss.TrainEnv.Nm = "TrainEnv"
								ss.TrainEnv.Dsc = "training params and state"
								ss.TrainEnv.Table = etable.NewIdxView(ss.Instr)
								ss.TrainEnv.Validate()
								ss.TrainEnv.Run.Max = ss.MaxRuns // note: we are not setting epoch max -- do that manually
								
								ss.TestEnv.Nm = "TestEnv"
								ss.TestEnv.Dsc = "testing params and state"
								ss.TestEnv.Table = etable.NewIdxView(ss.Instr)
								ss.TestEnv.Sequential = true
								ss.TestEnv.Validate()
								
								// note: to create a train / test split of pats, do this:
								// all := etable.NewIdxView(ss.Pats)
								// splits, _ := split.Permuted(all, []float64{.8, .2}, []string{"Train", "Test"})
								// ss.TrainEnv.Table = splits.Splits[0]
								// ss.TestEnv.Table = splits.Splits[1]
								
								ss.TrainEnv.Init(0)
								ss.TestEnv.Init(0)
							}
							
							func (ss *Sim) ConfigNet(net *leabra.Network) {
								net.InitName(net, "Dynamic Personality Model")
								enviro := net.AddLayer2D("Environment", 1, Parameters.Size, emer.Input)
								intero := net.AddLayer2D("InteroState", 1, Parameters.Size, emer.Input)
								
								app := net.AddLayer2D("Approach", 1, 5, emer.Target)
								av := net.AddLayer2D("Avoid", 1, 2, emer.Target)
								hid := net.AddLayer2D("Hidden", 3, 7, emer.Hidden)
								beh := net.AddLayer2D("Behavior", 1, 12, emer.Target)
								motb := net.AddLayer2D("MotiveBias", 1, 7, emer.Input)
								vta := net.AddLayer2D("VTA_DA", 1, 1, emer.Input)
								
								// use this to position layers relative to each other
								// default is Above, YAlign = Front, XAlign = Center
								av.SetRelPos(relpos.Rel{Rel: relpos.RightOf, Other: "Avoid", YAlign: relpos.Front, Space: 2})
								
								// note: see emergent/prjn module for all the options on how to connect
								// NewFull returns a new prjn.Full connectivity pattern
								full := prjn.NewFull()
								motb2app := prjn.NewOneToOne()
								
								motb2av := prjn.NewOneToOne()
								motb2av.SendStart = 5
								
								net.ConnectLayers(enviro, app, full, emer.Forward)
								net.ConnectLayers(enviro, av, full, emer.Forward)
								net.ConnectLayers(intero, app, full, emer.Forward)
								net.ConnectLayers(intero, av, full, emer.Forward)
								
								net.ConnectLayers(vta, app, full, emer.Forward)
								net.ConnectLayers(vta, av, full, emer.Inhib) // Inhibitory connection
								net.ConnectLayers(motb, app, motb2app,  emer.Forward)
								net.ConnectLayers(motb, av, motb2av,  emer.Forward) 
								
								net.BidirConnectLayers(hid, app, full)
								net.BidirConnectLayers(hid, av, full)
								net.BidirConnectLayers(hid, beh, full)
								
								// Commands that are used to position the layers in the Netview
								
								app.SetRelPos(relpos.Rel{Rel: relpos.Above, Other: "Environment", YAlign: relpos.Front, XAlign: relpos.Right, XOffset: 1})
								hid.SetRelPos(relpos.Rel{Rel: relpos.Above, Other: "Approach", YAlign: relpos.Front, XAlign: relpos.Left, YOffset: 0})
								beh.SetRelPos(relpos.Rel{Rel: relpos.Above, Other: "Hidden", YAlign: relpos.Front, XAlign: relpos.Right, XOffset: 1})
								
								intero.SetRelPos(relpos.Rel{Rel: relpos.RightOf, Other: "Environment", YAlign: relpos.Front, XAlign: relpos.Right, XOffset: 1})
								av.SetRelPos(relpos.Rel{Rel: relpos.RightOf, Other: "Approach", YAlign: relpos.Front, XAlign: relpos.Right, XOffset: 1})
								motb.SetRelPos(relpos.Rel{Rel: relpos.RightOf, Other: "Avoid", YAlign: relpos.Front, XAlign: relpos.Right, XOffset: 1})
								vta.SetRelPos(relpos.Rel{Rel: relpos.RightOf, Other: "InteroState", YAlign: relpos.Front, XAlign: relpos.Right, XOffset: 1})
								
								// note: can set these to do parallel threaded computation across multiple cpus
								// not worth it for this small of a model, but definitely helps for larger ones
								// if Thread {
									// 	hid2.SetThread(1)
									// 	out.SetThread(1)
									// }
									
									// note: if you wanted to change a layer type from e.g., Target to Compare, do this:
									// out.SetType(emer.Compare)
									// that would mean that the output layer doesn't reflect target values in plus phase
									// and thus removes error-driven learning -- but stats are still computed.
									
									net.Defaults()
									ss.SetParams("Network", ss.LogSetParams) // only set Network params
									err := net.Build()
									if err != nil {
										log.Println(err)
										return
									}
									net.InitWts()
								}
								
								////////////////////////////////////////////////////////////////////////////////
								// 	    Init, utils
								
								// Init restarts the run, and initializes everything, including network weights
								// and resets the epoch log table
								func (ss *Sim) Init() {
									rand.Seed(ss.RndSeed)
									ss.ConfigEnv() // re-config env just in case a different set of patterns was
									// selected or patterns have been modified etc
									ss.StopNow = false
									ss.SetParams("", ss.LogSetParams) // all sheets
									ss.NewRun()
									ss.UpdateView(true)
								}
								
								// NewRndSeed gets a new random seed based on current time -- otherwise uses
								// the same random seed for every run
								func (ss *Sim) NewRndSeed() {
									ss.RndSeed = time.Now().UnixNano()
								}
								
								// Counters returns a string of the current counter state
								// use tabs to achieve a reasonable formatting overall
								// and add a few tabs at the end to allow for expansion..
								func (ss *Sim) Counters(train bool) string {
									if train {
										return fmt.Sprintf("Run:\t%d\tEpoch:\t%d\tTrial:\t%d\tCycle:\t%d\tName:\t%v\t\t\t", ss.TrainEnv.Run.Cur, ss.TrainEnv.Epoch.Cur, ss.TrainEnv.Trial.Cur, ss.Time.Cycle, ss.TrainEnv.TrialName)
										} else {
											return fmt.Sprintf("Run:\t%d\tEpoch:\t%d\tTrial:\t%d\tCycle:\t%d\tName:\t%v\t\t\t", ss.TrainEnv.Run.Cur, ss.TrainEnv.Epoch.Cur, ss.TestEnv.Trial.Cur, ss.Time.Cycle, ss.TestEnv.TrialName)
										}
									}
									
									func (ss *Sim) UpdateView(train bool) {
										if ss.NetView != nil && ss.NetView.IsVisible() {
											ss.NetView.Record(ss.Counters(train), 0) // TODO: 0 was added arbitrarily to compile, update this later
											// note: essential to use Go version of update when called from another goroutine
											ss.NetView.GoUpdate() // note: using counters is significantly slower..
										}
									}
									
									////////////////////////////////////////////////////////////////////////////////
									// 	    Running the Network, starting bottom-up..
									
									// AlphaCyc runs one alpha-cycle (100 msec, 4 quarters)	of processing.
									// External inputs must have already been applied prior to calling,
									// using ApplyExt method on relevant layers (see TrainTrial, TestTrial).
									// If train is true, then learning DWt or WtFmDWt calls are made.
									// Handles netview updating within scope of AlphaCycle
									func (ss *Sim) AlphaCyc(train bool) {
										// ss.Win.PollEvents() // this can be used instead of running in a separate goroutine
										viewUpdt := ss.TrainUpdt
										if !train {
											viewUpdt = ss.TestUpdt
										}
										
										// update prior weight changes at start, so any DWt values remain visible at end
										// you might want to do this less frequently to achieve a mini-batch update
										// in which case, move it out to the TrainTrial method where the relevant
										// counters are being dealt with.
										if train {
											ss.Net.WtFmDWt()
										}
										
										ss.Net.AlphaCycInit(train)
										ss.Time.AlphaCycStart()
										for qtr := 0; qtr < 4; qtr++ {
											for cyc := 0; cyc < ss.Time.CycPerQtr; cyc++ {
												ss.Net.Cycle(&ss.Time)
												if !train {
													ss.LogTstCyc(ss.TstCycLog, ss.Time.Cycle)
												}
												ss.Time.CycleInc()
												if ss.ViewOn {
													switch viewUpdt {
													case leabra.Cycle:
														if cyc != ss.Time.CycPerQtr-1 { // will be updated by quarter
															ss.UpdateView(train)
														}
													case leabra.FastSpike:
														if (cyc+1)%10 == 0 {
															ss.UpdateView(train)
														}
													}
												}
											}
											ss.Net.QuarterFinal(&ss.Time)
											ss.Time.QuarterInc()
											if ss.ViewOn {
												switch {
												case viewUpdt <= leabra.Quarter:
													ss.UpdateView(train)
												case viewUpdt == leabra.Phase:
													if qtr >= 2 {
														ss.UpdateView(train)
													}
												}
											}
										}
										
										if train {
											ss.Net.DWt()
										}
										if ss.ViewOn && viewUpdt == leabra.AlphaCycle {
											ss.UpdateView(train)
										}
										if !train {
											ss.TstCycPlot.GoUpdate() // make sure up-to-date at end
										}
									}
									
									// ApplyInputs applies input patterns from given environment.
									// It is good practice to have this be a separate method with appropriate
									// args so that it can be used for various different contexts
									// (training, testing, etc).
									func (ss *Sim) ApplyInputs(en env.Env) {
										ss.Net.InitExt() // clear any existing inputs -- not strictly necessary if always
										// going to the same layers, but good practice and cheap anyway
										
										lays := []string{"Environment", "InteroState", "Approach", "Avoid", "Behavior", "VTA_DA", "MotiveBias"}
										for _, lnm := range lays {
											ly := ss.Net.LayerByName(lnm).(leabra.LeabraLayer).AsLeabra()
											pats := en.State(ly.Nm)
											if pats != nil {
												ly.ApplyExt(pats)
											}
										}
									}
									
									// TrainTrial runs one trial of training using TrainEnv
									func (ss *Sim) TrainTrial() {
										if ss.NeedsNewRun {
											ss.NewRun()
										}
										
										ss.TrainEnv.Step() // the Env encapsulates and manages all counter state
										
										// Key to query counters FIRST because current state is in NEXT epoch
										// if epoch counter has changed
										epc, _, chg := ss.TrainEnv.Counter(env.Epoch)
										if chg {
											ss.LogTrnEpc(ss.TrnEpcLog)
											if ss.ViewOn && ss.TrainUpdt > leabra.AlphaCycle {
												ss.UpdateView(true)
											}
											if ss.TestInterval > 0 && epc%ss.TestInterval == 0 { // note: epc is *next* so won't trigger first time
											ss.TestAll()
										}
										if epc >= ss.MaxEpcs || (ss.NZeroStop > 0 && ss.NZero >= ss.NZeroStop) {
											// done with training..
											ss.RunEnd()
											if ss.TrainEnv.Run.Incr() { // we are done!
												ss.StopNow = true
												return
												} else {
													ss.NeedsNewRun = true
													return
												}
											}
										}
										
										ss.ApplyInputs(&ss.TrainEnv)
										ss.AlphaCyc(true)   // train
										ss.TrialStats(true) // accumulate
									}
									
									// RunEnd is called at the end of a run -- save weights, record final log, etc here
									func (ss *Sim) RunEnd() {
										ss.LogRun(ss.RunLog)
										if ss.SaveWts {
											fnm := ss.WeightsFileName()
											fmt.Printf("Saving Weights to: %v\n", fnm)
											ss.Net.SaveWtsJSON(gi.FileName(fnm))
										}
									}
									
									// NewRun intializes a new run of the model, using the TrainEnv.Run counter
									// for the new run value
									func (ss *Sim) NewRun() {
										run := ss.TrainEnv.Run.Cur
										ss.TrainEnv.Init(run)
										ss.TestEnv.Init(run)
										ss.Time.Reset()
										ss.Net.InitWts()
										ss.Net.SaveWtsJSON("trained.wts")
										ss.InitStats()
										ss.TrnEpcLog.SetNumRows(0)
										ss.TstEpcLog.SetNumRows(0)
										ss.NeedsNewRun = false
									}
									
									// InitStats initializes all the statistics, especially important for the
									// cumulative epoch stats -- called at start of new run
									func (ss *Sim) InitStats() {
										// accumulators
										ss.SumErr = 0
										ss.SumSSE = 0
										ss.SumAvgSSE = 0
										ss.SumCosDiff = 0
										ss.FirstZero = -1
										ss.NZero = 0
										// clear rest just to make Sim look initialized
										ss.TrlErr = 0
										ss.TrlSSE = 0
										ss.TrlAvgSSE = 0
										ss.EpcSSE = 0
										ss.EpcAvgSSE = 0
										ss.EpcPctErr = 0
										ss.EpcCosDiff = 0
									}
									
									// TrialStats computes the trial-level statistics and adds them to the epoch accumulators if
									// accum is true.  Note that we're accumulating stats here on the Sim side so the
									// core algorithm side remains as simple as possible, and doesn't need to worry about
									// different time-scales over which stats could be accumulated etc.
									// You can also aggregate directly from log data, as is done for testing stats
									func (ss *Sim) TrialStats(accum bool) {
										out := ss.Net.LayerByName("Behavior").(leabra.LeabraLayer).AsLeabra()
										ss.TrlCosDiff = float64(out.CosDiff.Cos)
										ss.TrlSSE, ss.TrlAvgSSE = out.MSE(0.5) // 0.5 = per-unit tolerance -- right side of .5
										if ss.TrlSSE > 0 {
											ss.TrlErr = 1
											} else {
												ss.TrlErr = 0
											}
											if accum {
												ss.SumErr += ss.TrlErr
												ss.SumSSE += ss.TrlSSE
												ss.SumAvgSSE += ss.TrlAvgSSE
												ss.SumCosDiff += ss.TrlCosDiff
											}
										}
										
										// TrainEpoch runs training trials for remainder of this epoch
										func (ss *Sim) TrainEpoch() {
											ss.StopNow = false
											curEpc := ss.TrainEnv.Epoch.Cur
											for {
												ss.TrainTrial()
												if ss.StopNow || ss.TrainEnv.Epoch.Cur != curEpc {
													break
												}
											}
											ss.Stopped()
										}
										
										// TrainRun runs training trials for remainder of run
										func (ss *Sim) TrainRun() {
											ss.StopNow = false
											curRun := ss.TrainEnv.Run.Cur
											for {
												ss.TrainTrial()
												if ss.StopNow || ss.TrainEnv.Run.Cur != curRun {
													break
												}
											}
											ss.Stopped()
										}
										
										// Train runs the full training from this point onward
										func (ss *Sim) Train() {
											ss.StopNow = false
											for {
												ss.TrainTrial()
												if ss.StopNow {
													break
												}
											}
											ss.Stopped()
										}
										
										// Stop tells the sim to stop running
										func (ss *Sim) Stop() {
											ss.StopNow = true
										}
										
										// Stopped is called when a run method stops running -- updates the IsRunning flag and toolbar
										func (ss *Sim) Stopped() {
											ss.IsRunning = false
											if ss.Win != nil {
												vp := ss.Win.WinViewport2D()
												if ss.ToolBar != nil {
													ss.ToolBar.UpdateActions()
												}
												vp.SetNeedsFullRender()
											}
										}
										
										// SaveWeights saves the network weights -- when called with giv.CallMethod
										// it will auto-prompt for filename
										func (ss *Sim) SaveWeights(filename gi.FileName) {
											ss.Net.SaveWtsJSON(filename)
										}
										
										// Following code is for training the network with both Pavlovian training and Instrumental training
										// controls order of training and number of Epochs of training by reading in an etable that has that information.
										func (ss *Sim) TrainPIT() {
											
											for i := 0; i < 2; i++ {
												
												ss.Training = ss.Trn.CellString("Training", i)
												ss.MaxEpcs = int(ss.Trn.CellFloat("MaxEpoch", i))
												switch ss.Training {
													
												case "INSTRUMENTAL":
													
													ss.TrainEnv.Epoch.Cur = 0  //set current epoch to 0 so that training starts from 0 epochs
													
													// Need to set up training patterns.  	
													ss.TrainEnv.Table = etable.NewIdxView(ss.Instr)
													ss.TestEnv.Table = etable.NewIdxView(ss.Instr)
													
													// Unlesion Hidden and Behavior layer to make sure all layers are unlesioned
													ss.Net.LayerByName("Hidden").SetOff(false)
													ss.Net.LayerByName("Behavior").SetOff(false)
													
													// Load saved weights
													// OpenWtsJSON opens trained weights
													ss.Net.OpenWtsJSON("trained.wts")
													// Define Approach and Avoid as Input layers
													// Define Behavior as a Target layer
													ss.Net.LayerByName("Environment").SetType(emer.Input)
													ss.Net.LayerByName("InteroState").SetType(emer.Input)
													ss.Net.LayerByName("Approach").SetType(emer.Input)
													ss.Net.LayerByName("Avoid").SetType(emer.Input)
													ss.Net.LayerByName("Behavior").SetType(emer.Target)
													// Lesion Environment and InteroState layers
													ss.Net.LayerByName("Environment").SetOff(true)
													ss.Net.LayerByName("InteroState").SetOff(true)
													ss.Train()
													// Unlesion Environment and InteroState layers
													ss.Net.LayerByName("Environment").SetOff(false)
													ss.Net.LayerByName("InteroState").SetOff(false)
													// Save weights
													ss.Net.SaveWtsJSON("trained.wts")
													// Define Environment and InteroState as Input layers
													// Define Approach and Avoid as Target layers
													// Define Behavior as a Target layer
													// Makes sure that layers are set to default.
													ss.Net.LayerByName("Environment").SetType(emer.Input)
													ss.Net.LayerByName("InteroState").SetType(emer.Input)
													ss.Net.LayerByName("Approach").SetType(emer.Target)
													ss.Net.LayerByName("Avoid").SetType(emer.Target)	
													ss.Net.LayerByName("Behavior").SetType(emer.Target)
													
												case "PAVLOV":
													
													ss.TrainEnv.Epoch.Cur = 0   //set current epoch to 0 so that training starts from 0 epochs
													
													// Pavlovian Training
													
													// Need to set up training patterns.  
													ss.TrainEnv.Table = etable.NewIdxView(ss.Pvlv)
													ss.TestEnv.Table = etable.NewIdxView(ss.Pvlv)
													
													// Define Environment and InteroState as Input layers
													// Define Approach and Avoid as Target layers
													// Define Behavior as a Compare layer
													ss.Net.LayerByName("Environment").SetType(emer.Input)
													ss.Net.LayerByName("InteroState").SetType(emer.Input)
													ss.Net.LayerByName("Approach").SetType(emer.Target)
													ss.Net.LayerByName("Avoid").SetType(emer.Target)
													ss.Net.LayerByName("Behavior").SetType(emer.Target)
													
													// Load saved weights
													ss.Net.OpenWtsJSON("trained.wts")
													// Lesion Hidden layer and Behavior layer
													ss.Net.LayerByName("Hidden").SetOff(true)
													ss.Net.LayerByName("Behavior").SetOff(true)
													// Train until number of Epochs of training reached
													ss.Train()
													
													//Unlesion Hidden Layer and Behavior layer
													ss.Net.LayerByName("Hidden").SetOff(false)
													ss.Net.LayerByName("Behavior").SetOff(false)
													// Save weights
													ss.Net.SaveWtsJSON("trained.wts")
													// Define Environment and InteroState as Input layers
													// Define Approach and Avoid as Target layers
													// Define Behavior as a Target layer
													// Makes sure that layers are set to default.
													
													ss.Net.LayerByName("Environment").SetType(emer.Input)
													ss.Net.LayerByName("InteroState").SetType(emer.Input)
													ss.Net.LayerByName("Approach").SetType(emer.Target)
													ss.Net.LayerByName("Avoid").SetType(emer.Target)	
													ss.Net.LayerByName("Behavior").SetType(emer.Target)
												}
											}
										}
										////////////////////////////////////////////////////////////////////////////////////////////
										// Testing
										
										// TestTrial runs one trial of testing -- always sequentially presented inputs
										func (ss *Sim) TestTrial(returnOnChg bool) {
											ss.TestEnv.Step()
											
											// Query counters FIRST
											_, _, chg := ss.TestEnv.Counter(env.Epoch)
											if chg {
												if ss.ViewOn && ss.TestUpdt > leabra.AlphaCycle {
													ss.UpdateView(false)
												}
												ss.LogTstEpc(ss.TstEpcLog)
												if returnOnChg {
													return
												}
											}
											
											ss.ApplyInputs(&ss.TestEnv)
											ss.AlphaCyc(false)   // !train
											ss.TrialStats(false) // !accumulate
											ss.LogTstTrl(ss.TstTrlLog)
										}
										// This function has been heavily modified from the TestTrial function  to dynamically change the External and Internal world state as a function of the model's behavior and exogenous changes.  
										// WARNING!!! THIS CODE IS A WORK IN PROGRESS. IT HAS NOT BEEN FINISHED OR TESTED.
										// World file: first row represents initial State of the World.  Subsequent rows initially represent Exogenous Changes to the world.
										// Will calculate new state of the world (next row, t) by recording: 1) previous Environment, previous Interostate, (t-1)
										// 2) Behavior (t-1), then 3) calculate new Environment and Interostate (t) by calculating a) changes in both Environment and Interostate
										// due to behavior, b) changes over time (e.g., getting hungry), and c) exogenous changes in the world represented in row t of World, and 
										// 4) write all of this into row t of World.  
										
										// 
										//
										// WARNING!! Dynamics only runs one trial. Need to have something like TestAll call Dynamics until stopping point reached.
										//
										//
										
										func (ss *Sim) Dynamics(returnOnChg bool) {
											ss.TestEnv.Step()
											// TestTrial is called periodically during training, by default.  TestInterval is currently set to -1 to turn that off.
											
											if ss.TestEnv.Trial.Cur == 0 {
												ss.TestEnv.Table = etable.NewIdxView(ss.World) // define a World Struc item and then use "ss.World"
												// datatable for World is read in by OpenPats function
												// row 0 of World is initial state of World (Enviro and Interostate)
												// subsequent rows will initially have Changes, and these Changes will be used in conjunction with other info
												// to calculate new state of the World
												ss.TestEnv.NewOrder()	
												
												// basic idea is to call this if you read in a new testing file and then have it update the indexes and the order
												// so that they will correspond to the size of the new file
												// does ss.TestEnv.NewOrder() work properly as long as sequential is set in the ConfigEnv function?
												
											}
											// Query counters FIRST
											_, _, chg := ss.TestEnv.Counter(env.Epoch)
											if chg {
												if ss.ViewOn && ss.TestUpdt > leabra.AlphaCycle {
													ss.UpdateView(false)
												}
												ss.LogTstEpc(ss.TstEpcLog)
												if returnOnChg {
													return
												}
											}
											
											en := ss.TestEnv   // would this line be correct? should this be set up with the &
											
											// should we put the if statement here and then the function call that updates the WorldState?
											if ss.TestEnv.Trial.Cur < 1 { 
												// Code to read input values for Environment and Interoceptive State
												//		ss.ValsTsr("Envp") = en.State("Environment") // previous Environment
												//		ss.ValsTsr("Intp") = en.State("InteroState")	  // previous InteroState
												
												// TODO: proceduralize this section.  shouldn't be too hard once you finish figuring out how you are going to order the layer arrays.
												
												ss.tsrsStack.Update(REPLACE_Val, en.State("Enviro").(*etensor.Float32), FIND_Key, "EnvpTsr") // previous Environment
												ss.tsrsStack.Update(REPLACE_Val, en.State("Intero").(*etensor.Float32), FIND_Key, "IntpTsr") // previous Environment
												
												ss.ApplyInputs(&ss.TestEnv)
												ss.AlphaCyc(false)   // !train
												
												// Code to read Behavior activations after AlphaCyc applied
												
												beh := ss.Net.LayerByName("Behavior").(leabra.LeabraLayer).AsLeabra()
												bh := ss.ValsTsr("Behavior") // see ra25 example for this method for a reusable map of tensors
												beh.UnitValsTensor(bh, "ActM") // read the actMs into tensor
												ss.TrialStats(false) // !accumulate
												ss.LogTstTrl(ss.TstTrlLog)
												
												} else {
													
													// ! main emergentstack section:
													
													// take previous inputs (cur-1), retrieve current values(cur) and behavior activations (see below)
													// calculate new state of Environment and InteroState and put in TestEnv Table
													// do below 
													// NEED TO READ UP ON APPROPRIATE WAY TO USE ASSIGNMENT STATEMENTS VERSUS =
													// take envp, intp, and beh and then with the following, calculate new World State
													// 		en := &ss.TestEnv  // already assigned in "if statement" 
													
													// envc := en.State("Environment").(*etensor.Float32) // exogenous Changes to current Environment
													// intc := en.State("InteroState").(*etensor.Float32)	  // exogenous Changes to current InteroState
													// Are there exogenous changes to Interostate?  maybe we don't need intc.
													
													// TODO: discuss envp and intp tensors: do we need them if we're plotting to a csv each frame?
													// TODO: discuss role of envc/intc vs EnviroTsr and InteroTsr
													// envp := ss.EnvpTsr
													// intp := ss.IntpTsr
													
													// TODO: replace this with layerTensors; see EnviroTsr initialization to-do comment
													enviro := ss.tsrsStack.Get(FIND_Key, "enviro").Val.(*etensor.Float32)  // This is used to create new Environment representation
													intero := ss.tsrsStack.Get(FIND_Key, "intero").Val.(*etensor.Float32)   // This is used to create new InteroState representation
													
													layerTensors := MakeStack(Layers.ToArray(RETURN_Keys), []*etensor.Float32{enviro, intero}) // * dummy stack until you can figure out how to procedurally add EnviroTsr and InteroTsr
													
													SetupModel() // TODO: reorder setup model contingent on how you decide to deal with layerTensors (current implementation doesn't work sicne Layers isn't initialized yet)
													
													// Then calculate new enviro and new intero which will be written to the World datatable
													// This code takes the tensor from the Behavior layer and then finds the index of the most strongly activated behavior.
													
													bh := ss.ValsTsr("Behavior")
													_, _, _, maxBHIdx := bh.Range() // TODO: have behavior-generating minimum threshold, rather than just grab the max
													
													// then set the Behavior tensor to all zeros and then set the value at the index to 1.
													// so this tensor identifies which behavior is Performed or Enacted
													
													bh.SetZeros()
													bh.SetFloat1D(maxBHIdx, 1.0) // selected behavior => 1, rest are set to 0, found from idx "maxBHIdx"
													
													// * Update Parameters
													for _, card := range Parameters.Cards {
														
														// initialize variables
														parameterName := card.Key.(string)
														parameterData := card.Val.(*Stack)
														
														layers := parameterData.Get(FIND_Key, "layerValues").Val.(*Stack)
														initialLayers := MakeStack()
														
														dx_ui := parameterData.Get(FIND_Key, "dx_ui").Val.(float32)
														tprev_s := parameterData.Get(FIND_Key, "tprev_s").Val.(float32)
														tcur_s := parameterData.Get(FIND_Key, "tcur_s").Val.(float32)
														dt_s := parameterData.Get(FIND_Key, "dt_s").Val.(float32)
														
														timeIncrements := parameterData.Get(FIND_Key, "timeIncrements").Val.(*Stack) // get a stack containing each timeincrement function, corresponding to the desired incremented layer
														
														actions := parameterData.Get(FIND_Key, "actions").Val.(*Stack)
														relations := parameterData.Get(FIND_Key, "relations").Val.(*Stack)
														
														// update time increments
														for _, layerName := range layers.ToArray(RETURN_Keys) {
															
															x := layers.Get(FIND_Key, layerName).Val.(*float32)
															initialLayers.Add(layerName, gogenerics.CloneObject(*x)) // keep a save of the original layers so you can detect how much they've changed by then end of increments and actions
															
															thisTimeIncrement := timeIncrements.Get(FIND_Key, layerName).Val.(func(float32, float32, float32, float32, float32) float32)
															
															// update x/thisLayerValue to new float32
															*x = thisTimeIncrement(*x, dx_ui, tprev_s, tcur_s, dt_s)
															
														}
														
														// perform actions
														PerformActions(actions, parameterName)
														
														// update relations
														for _, _relation := range relations.ToArray() {
															
															relation := _relation.(*Relation)
															deltaX := *initialLayers.Get(FIND_Key, relation.ThisLayer).Val.(*float32) - *layers.Get(FIND_Key, relation.ThisLayer).Val.(*float32)
															
															// multiply the other parameter's value by the rate of change for the other parameter times the amount by which this parameter changed
															// TODO: update math, write it out in LaTeX
															*Parameters.Get(FIND_Key, relation.OtherParameter).Val.(*Stack).Get(FIND_Key, relation.OtherLayer).Val.(*float32) += (relation.Dx * deltaX)
															
														}
														
													}
													
													// * Perform ComplexActions
													PerformActions(ComplexActions, "")
													
													trl := ss.TestEnv.Trial.Cur
													row := trl
													
													for _, tsrCard := range layerTensors.Cards {
														
														tsrName := tsrCard.Key.(string)
														tsrVal := tsrCard.Val.(*etensor.Float32)
														
														ss.World.SetCellTensor(tsrName, row, tsrVal)
														
													}
													
													// Apply new Environment and InteroState to network
													
													ss.ApplyInputs(&ss.TestEnv)
													ss.AlphaCyc(false)   // !train
													
													// Code to read Behavior activations after AlphaCyc applied
													
													beh := ss.Net.LayerByName("Behavior").(leabra.LeabraLayer).AsLeabra()
													bh = ss.ValsTsr("Behavior") // see ra25 example for this method for a reusable map of tensors
													beh.UnitValsTensor(bh, "ActM") // read the actMs into tensor
													ss.tsrsStack.Update(REPLACE_Val, enviro.Clone(), FIND_Key, "enviroPrev")  // saves current environment  as previous Environment for next step
													ss.tsrsStack.Update(REPLACE_Val, enviro.Clone(), FIND_Key, "interoPrev")	  // saves current InteroState as previous InteroState for next step
													ss.TrialStats(false) // !accumulate
													ss.LogTstTrl(ss.TstTrlLog)
													
												}
											}
											// TODO: Implement this section if necessary for a first working version of the model?  Ask Dr. Read
											/*
											func (ss &Sim) UpdateWorld() {
												// Code to read previous or t-1 Environment and InteroState values from Network
												enviro := ss.Net.LayerByName(“Environment”).(leabra.LeabraLayer).AsLeabra()
												envp := ss.ValsTsr(“Environment”) // see ra25 example for this method for a reusable map of tensors
												enviro.UnitValsTensor(envp, "Act”) // read the acts into tensor
												intero := ss.Net.LayerByName(“InteroState”).(leabra.LeabraLayer).AsLeabra()
												intp := ss.ValsTsr(“InteroState) // see ra25 example for this method for a reusable map of tensors
												intero.UnitValsTensor(intp, "Act”) // read the acts into tensor
											}
											*/
											
											// TestItem tests given item which is at given index in test item list
											func (ss *Sim) TestItem(idx int) {
												cur := ss.TestEnv.Trial.Cur
												ss.TestEnv.Trial.Cur = idx
												ss.TestEnv.SetTrialName()
												ss.ApplyInputs(&ss.TestEnv)
												ss.AlphaCyc(false)   // !train
												ss.TrialStats(false) // !accumulate
												ss.TestEnv.Trial.Cur = cur
											}
											// Need some version of TestAll and RunTestAll to run the testing with a new environment
											
											// TestAll runs through the full set of testing items
											func (ss *Sim) TestAll() {
												ss.TestEnv.Init(ss.TrainEnv.Run.Cur)
												//	ss.World = Clone(ss.WorldChanges)
												//	ss.TestEnv.Table = etable.NewIdxView(ss.World)
												
												for {
													ss.TestTrial(true) // return on change -- don't wrap
													_, _, chg := ss.TestEnv.Counter(env.Epoch)
													if chg || ss.StopNow {
														break
													}
												}
											}
											
											// RunTestAll runs through the full set of testing items, has stop running = false at end -- for gui
											func (ss *Sim) RunTestAll() {
												ss.StopNow = false
												ss.TestAll()
												ss.Stopped()
											}
											
											/////////////////////////////////////////////////////////////////////////
											//   Params setting
											
											// ParamsName returns name of current set of parameters
											func (ss *Sim) ParamsName() string {
												if ss.ParamSet == "" {
													return "Base"
												}
												return ss.ParamSet
											}
											
											// SetParams sets the params for "Base" and then current ParamSet.
											// If sheet is empty, then it applies all avail sheets (e.g., Network, Sim)
											// otherwise just the named sheet
											// if setMsg = true then we output a message for each param that was set.
											func (ss *Sim) SetParams(sheet string, setMsg bool) error {
												if sheet == "" {
													// this is important for catching typos and ensuring that all sheets can be used
													ss.Params.ValidateSheets([]string{"Network", "Sim"})
												}
												err := ss.SetParamsSet("Base", sheet, setMsg)
												if ss.ParamSet != "" && ss.ParamSet != "Base" {
													err = ss.SetParamsSet(ss.ParamSet, sheet, setMsg)
												}
												return err
											}
											
											// SetParamsSet sets the params for given params.Set name.
											// If sheet is empty, then it applies all avail sheets (e.g., Network, Sim)
											// otherwise just the named sheet
											// if setMsg = true then we output a message for each param that was set.
											func (ss *Sim) SetParamsSet(setNm string, sheet string, setMsg bool) error {
												pset, err := ss.Params.SetByNameTry(setNm)
												if err != nil {
													return err
												}
												if sheet == "" || sheet == "Network" {
													netp, ok := pset.Sheets["Network"]
													if ok {
														ss.Net.ApplyParams(netp, setMsg)
													}
												}
												
												if sheet == "" || sheet == "Sim" {
													simp, ok := pset.Sheets["Sim"]
													if ok {
														simp.Apply(ss, setMsg)
													}
												}
												// note: if you have more complex environments with parameters, definitely add
												// sheets for them, e.g., "TrainEnv", "TestEnv" etc
												return err
											}
											
											// TODO: also discuss this section with Dr. Read
											/*
											Following function configures a data table to fit the structure of the network.
											func (ss *Sim) ConfigPats() {
												dt := ss.Pats
												dt.SetMetaData("name", "TrainPats")
												dt.SetMetaData("desc", "Training patterns")
												dt.SetFromSchema(etable.Schema{
													{"Name", etensor.STRING, nil, nil},
													{"Environment", etensor.FLOAT32, []int{1, 7}, []string{"Y", "X"}},
													{"InteroState", etensor.FLOAT32, []int{1, 7}, []string{"Y", "X"}},
													{"Approach", etensor.FLOAT32, []int{1, 5}, []string{"Y", "X"}},
													{"Avoid", etensor.FLOAT32, []int{1, 2}, []string{"Y", "X"}},
													{"Behavior", etensor.FLOAT32, []int{1, 12}, []string{"Y", "X"}},
													{"MotiveBias", etensor.FLOAT32, []int{1, 7}, []string{"Y", "X"}},
													{"VTA_DA", etensor.FLOAT32, []int{1, 1}, []string{"Y", "X"}},
													}, 25)
													patgen.PermutedBinaryRows(dt.Cols[1], 1, 1, 0)
													patgen.PermutedBinaryRows(dt.Cols[2], 1, 1, 0)
													patgen.PermutedBinaryRows(dt.Cols[3], 1, 1, 0)
													patgen.PermutedBinaryRows(dt.Cols[4], 1, 1, 0)
													patgen.PermutedBinaryRows(dt.Cols[5], 1, 1, 0)
													patgen.PermutedBinaryRows(dt.Cols[6], 1, 1, 0)
													patgen.PermutedBinaryRows(dt.Cols[7], 1, 1, 0)
													
													dt.SaveCSV("PersonalityTraining.tsv", etable.Tab, etable.Headers)
												}
												*/
												
												func (ss *Sim) OpenPats() {
													ss.Instr.OpenCSV("instr.tsv", etable.Tab) // Instrumental training data
													// ss.OpenPatAsset(ss.Hard, "hard.tsv", "Hard", "Hard Training patterns")
													ss.Pvlv.OpenCSV("pvlv.tsv", etable.Tab) // Pavlovian training data
													// ss.OpenPatAsset(ss.Impossible, "impossible.tsv", "Impossible", "Impossible Training patterns")
													ss.Trn.OpenCSV("InstrThenPvlv.tsv", etable.Tab) // Order of training and number of epochs for each, 
													// Pavlov first or Instrumental first. Should eventually create menu to choose.
													//
													ss.World.OpenCSV("World.tsv", etable.Tab) // Current state of the World	
													// ss.WorldChanges.OpenCSV("WorldChanges.tsv", etable.Tab) // Exogenous changes in State of the World
													
													// Currently the program is set up so that when first read in, the first row of World represents the initial state of the world
													// and subsequent rows represent changes in the World.  
													// When the network starts behaving, it will modify each row of the network after the first one
													// so that it takes into account all changes in both External and Internal state that will be input
													// to the network.  
												}
												
												////////////////////////////////////////////////////////////////////////////////////////////
												// 		Logging
												
												// ValsTsr gets value tensor of given name, creating if not yet made
												func (ss *Sim) ValsTsr(name string) *etensor.Float32 {
													if ss.ValsTsrs == nil {
														ss.ValsTsrs = make(map[string]*etensor.Float32)
													}
													tsr, ok := ss.ValsTsrs[name]
													if !ok {
														tsr = &etensor.Float32{}
														ss.ValsTsrs[name] = tsr
													}
													return tsr
												}
												
												// RunName returns a name for this run that combines Tag and Params -- add this to
												// any file names that are saved.
												func (ss *Sim) RunName() string {
													if ss.Tag != "" {
														return ss.Tag + "_" + ss.ParamsName()
														} else {
															return ss.ParamsName()
														}
													}
													
													// RunEpochName returns a string with the run and epoch numbers with leading zeros, suitable
													// for using in weights file names.  Uses 3, 5 digits for each.
													func (ss *Sim) RunEpochName(run, epc int) string {
														return fmt.Sprintf("%03d_%05d", run, epc)
													}
													
													// WeightsFileName returns default current weights file name
													func (ss *Sim) WeightsFileName() string {
														return ss.Net.Nm + "_" + ss.RunName() + "_" + ss.RunEpochName(ss.TrainEnv.Run.Cur, ss.TrainEnv.Epoch.Cur) + ".wts"
													}
													
													// LogFileName returns default log file name
													func (ss *Sim) LogFileName(lognm string) string {
														return ss.Net.Nm + "_" + ss.RunName() + "_" + lognm + ".csv"
													}
													
													//////////////////////////////////////////////
													//  TrnEpcLog
													
													// LogTrnEpc adds data from current epoch to the TrnEpcLog table.
													// computes epoch averages prior to logging.
													func (ss *Sim) LogTrnEpc(dt *etable.Table) {
														row := dt.Rows
														dt.SetNumRows(row + 1)
														
														epc := ss.TrainEnv.Epoch.Prv          // this is triggered by increment so use previous value
														nt := float64(len(ss.TrainEnv.Order)) // number of trials in view
														
														ss.EpcSSE = ss.SumSSE / nt
														ss.SumSSE = 0
														ss.EpcAvgSSE = ss.SumAvgSSE / nt
														ss.SumAvgSSE = 0
														ss.EpcPctErr = float64(ss.SumErr) / nt
														ss.SumErr = 0
														ss.EpcPctCor = 1 - ss.EpcPctErr
														ss.EpcCosDiff = ss.SumCosDiff / nt
														ss.SumCosDiff = 0
														if ss.FirstZero < 0 && ss.EpcPctErr == 0 {
															ss.FirstZero = epc
														}
														if ss.EpcPctErr == 0 {
															ss.NZero++
															} else {
																ss.NZero = 0
															}
															
															if ss.LastEpcTime.IsZero() {
																ss.EpcPerTrlMSec = 0
																} else {
																	iv := time.Now().Sub(ss.LastEpcTime)
																	ss.EpcPerTrlMSec = float64(iv) / (nt * float64(time.Millisecond))
																}
																ss.LastEpcTime = time.Now()
																
																dt.SetCellFloat("Run", row, float64(ss.TrainEnv.Run.Cur))
																dt.SetCellFloat("Epoch", row, float64(epc))
																dt.SetCellFloat("SSE", row, ss.EpcSSE)
																dt.SetCellFloat("AvgSSE", row, ss.EpcAvgSSE)
																dt.SetCellFloat("PctErr", row, ss.EpcPctErr)
																dt.SetCellFloat("PctCor", row, ss.EpcPctCor)
																dt.SetCellFloat("CosDiff", row, ss.EpcCosDiff)
																dt.SetCellFloat("PerTrlMSec", row, ss.EpcPerTrlMSec)
																
																for _, lnm := range ss.LayStatNms {
																	ly := ss.Net.LayerByName(lnm).(leabra.LeabraLayer).AsLeabra()
																	dt.SetCellFloat(ly.Nm+" ActAvg", row, float64(ly.Pools[0].ActAvg.ActPAvgEff))
																}
																
																// note: essential to use Go version of update when called from another goroutine
																ss.TrnEpcPlot.GoUpdate()
																if ss.TrnEpcFile != nil {
																	if ss.TrainEnv.Run.Cur == 0 && epc == 0 {
																		dt.WriteCSVHeaders(ss.TrnEpcFile, etable.Tab)
																	}
																	dt.WriteCSVRow(ss.TrnEpcFile, row, etable.Tab)
																}
															}
															
															func (ss *Sim) ConfigTrnEpcLog(dt *etable.Table) {
																dt.SetMetaData("name", "TrnEpcLog")
																dt.SetMetaData("desc", "Record of performance over epochs of training")
																dt.SetMetaData("read-only", "true")
																dt.SetMetaData("precision", strconv.Itoa(LogPrec))
																
																sch := etable.Schema{
																	{"Run", etensor.INT64, nil, nil},
																	{"Epoch", etensor.INT64, nil, nil},
																	{"SSE", etensor.FLOAT64, nil, nil},
																	{"AvgSSE", etensor.FLOAT64, nil, nil},
																	{"PctErr", etensor.FLOAT64, nil, nil},
																	{"PctCor", etensor.FLOAT64, nil, nil},
																	{"CosDiff", etensor.FLOAT64, nil, nil},
																	{"PerTrlMSec", etensor.FLOAT64, nil, nil},
																}
																for _, lnm := range ss.LayStatNms {
																	sch = append(sch, etable.Column{lnm + " ActAvg", etensor.FLOAT64, nil, nil})
																}
																dt.SetFromSchema(sch, 0)
															}
															
															func (ss *Sim) ConfigTrnEpcPlot(plt *eplot.Plot2D, dt *etable.Table) *eplot.Plot2D {
																plt.Params.Title = "Leabra Random Associator 25 Epoch Plot"
																plt.Params.XAxisCol = "Epoch"
																plt.SetTable(dt)
																// order of params: on, fixMin, min, fixMax, max
																plt.SetColParams("Run", eplot.Off, eplot.FixMin, 0, eplot.FloatMax, 0)
																plt.SetColParams("Epoch", eplot.Off, eplot.FixMin, 0, eplot.FloatMax, 0)
																plt.SetColParams("SSE", eplot.Off, eplot.FixMin, 0, eplot.FloatMax, 0)
																plt.SetColParams("AvgSSE", eplot.Off, eplot.FixMin, 0, eplot.FloatMax, 0)
																plt.SetColParams("PctErr", eplot.On, eplot.FixMin, 0, eplot.FixMax, 1) // default plot
																plt.SetColParams("PctCor", eplot.On, eplot.FixMin, 0, eplot.FixMax, 1) // default plot
																plt.SetColParams("CosDiff", eplot.Off, eplot.FixMin, 0, eplot.FixMax, 1)
																plt.SetColParams("PerTrlMSec", eplot.Off, eplot.FixMin, 0, eplot.FloatMax, 0)
																
																for _, lnm := range ss.LayStatNms {
																	plt.SetColParams(lnm+" ActAvg", eplot.Off, eplot.FixMin, 0, eplot.FixMax, .5)
																}
																return plt
															}
															
															//////////////////////////////////////////////
															//  TstTrlLog
															
															// LogTstTrl adds data from current trial to the TstTrlLog table.
															// log always contains number of testing items
															func (ss *Sim) LogTstTrl(dt *etable.Table) {
																epc := ss.TrainEnv.Epoch.Prv // this is triggered by increment so use previous value
																// TODO: proceduralize this section.  shouldn't be too hard.
																enviro := ss.Net.LayerByName("Environment").(leabra.LeabraLayer).AsLeabra()
																intero := ss.Net.LayerByName("InteroState").(leabra.LeabraLayer).AsLeabra()
																app := ss.Net.LayerByName("Approach").(leabra.LeabraLayer).AsLeabra()
																av := ss.Net.LayerByName("Avoid").(leabra.LeabraLayer).AsLeabra()
																beh := ss.Net.LayerByName("Behavior").(leabra.LeabraLayer).AsLeabra()
																trl := ss.TestEnv.Trial.Cur
																row := trl
																
																if dt.Rows <= row {
																	dt.SetNumRows(row + 1)
																}
																
																dt.SetCellFloat("Run", row, float64(ss.TrainEnv.Run.Cur))
																dt.SetCellFloat("Epoch", row, float64(epc))
																dt.SetCellFloat("Trial", row, float64(trl))
																dt.SetCellString("TrialName", row, ss.TestEnv.TrialName.Cur)
																dt.SetCellFloat("Err", row, ss.TrlErr)
																dt.SetCellFloat("SSE", row, ss.TrlSSE)
																dt.SetCellFloat("AvgSSE", row, ss.TrlAvgSSE)
																dt.SetCellFloat("CosDiff", row, ss.TrlCosDiff)
																
																for _, lnm := range ss.LayStatNms {
																	ly := ss.Net.LayerByName(lnm).(leabra.LeabraLayer).AsLeabra()
																	dt.SetCellFloat(ly.Nm+" ActM.Avg", row, float64(ly.Pools[0].ActM.Avg))
																}
																envt := ss.ValsTsr("Environment")
																intt := ss.ValsTsr("InteroState")
																appt := ss.ValsTsr("Approach")
																avt := ss.ValsTsr("Avoid")
																beht := ss.ValsTsr("Behavior")
																
																enviro.UnitValsTensor(envt, "Act")
																dt.SetCellTensor("EnviroAct", row, envt)
																intero.UnitValsTensor(intt, "Act")
																dt.SetCellTensor("InteroAct", row, intt)
																app.UnitValsTensor(appt, "Act")
																dt.SetCellTensor("AppAct", row, appt)
																av.UnitValsTensor(avt, "Act")
																dt.SetCellTensor("AvAct", row, avt)
																
																beh.UnitValsTensor(beht, "ActM")
																dt.SetCellTensor("BehActM", row, beht)
																
																beh.UnitValsTensor(beht, "ActP")
																dt.SetCellTensor("BehActP", row, beht)
																
																// note: essential to use Go version of update when called from another goroutine
																ss.TstTrlPlot.GoUpdate()
															}
															
															func (ss *Sim) ConfigTstTrlLog(dt *etable.Table) {
																// TODO: proceduralize this section.  shouldn't be too hard.
																enviro := ss.Net.LayerByName("Environment").(leabra.LeabraLayer).AsLeabra()
																intero := ss.Net.LayerByName("InteroState").(leabra.LeabraLayer).AsLeabra()
																app := ss.Net.LayerByName("Approach").(leabra.LeabraLayer).AsLeabra()
																av := ss.Net.LayerByName("Avoid").(leabra.LeabraLayer).AsLeabra()
																beh := ss.Net.LayerByName("Behavior").(leabra.LeabraLayer).AsLeabra()
																
																dt.SetMetaData("name", "TstTrlLog")
																dt.SetMetaData("desc", "Record of testing per input pattern")
																dt.SetMetaData("read-only", "true")
																dt.SetMetaData("precision", strconv.Itoa(LogPrec))
																
																nt := ss.TestEnv.Table.Len() // number in view
																sch := etable.Schema{
																	{"Run", etensor.INT64, nil, nil},
																	{"Epoch", etensor.INT64, nil, nil},
																	{"Trial", etensor.INT64, nil, nil},
																	{"TrialName", etensor.STRING, nil, nil},
																	{"Err", etensor.FLOAT64, nil, nil},
																	{"SSE", etensor.FLOAT64, nil, nil},
																	{"AvgSSE", etensor.FLOAT64, nil, nil},
																	{"CosDiff", etensor.FLOAT64, nil, nil},
																}
																for _, lnm := range ss.LayStatNms {
																	sch = append(sch, etable.Column{lnm + " ActM.Avg", etensor.FLOAT64, nil, nil})
																}
																sch = append(sch, etable.Schema{
																	{"EnviroAct", etensor.FLOAT64, enviro.Shp.Shp, nil},
																	{"InteroAct", etensor.FLOAT64, intero.Shp.Shp, nil},
																	{"AppActM", etensor.FLOAT64, app.Shp.Shp, nil},
																	{"AppActP", etensor.FLOAT64, app.Shp.Shp, nil},
																	{"AvActM", etensor.FLOAT64, av.Shp.Shp, nil},
																	{"AvActP", etensor.FLOAT64, av.Shp.Shp, nil},
																	{"BehActM", etensor.FLOAT64, beh.Shp.Shp, nil},
																	{"BehActP", etensor.FLOAT64, beh.Shp.Shp, nil},
																	}...)
																	dt.SetFromSchema(sch, nt)
																}
																
																func (ss *Sim) ConfigTstTrlPlot(plt *eplot.Plot2D, dt *etable.Table) *eplot.Plot2D {
																	plt.Params.Title = "Personality Model Test Trial Plot"
																	plt.Params.XAxisCol = "Trial"
																	plt.SetTable(dt)
																	// order of params: on, fixMin, min, fixMax, max
																	plt.SetColParams("Run", eplot.Off, eplot.FixMin, 0, eplot.FloatMax, 0)
																	plt.SetColParams("Epoch", eplot.Off, eplot.FixMin, 0, eplot.FloatMax, 0)
																	plt.SetColParams("Trial", eplot.Off, eplot.FixMin, 0, eplot.FloatMax, 0)
																	plt.SetColParams("TrialName", eplot.Off, eplot.FixMin, 0, eplot.FloatMax, 0)
																	plt.SetColParams("Err", eplot.Off, eplot.FixMin, 0, eplot.FloatMax, 0)
																	plt.SetColParams("SSE", eplot.Off, eplot.FixMin, 0, eplot.FloatMax, 0)
																	plt.SetColParams("AvgSSE", eplot.On, eplot.FixMin, 0, eplot.FloatMax, 0)
																	plt.SetColParams("CosDiff", eplot.On, eplot.FixMin, 0, eplot.FixMax, 1)
																	
																	for _, lnm := range ss.LayStatNms {
																		plt.SetColParams(lnm+" ActM.Avg", eplot.Off, eplot.FixMin, 0, eplot.FixMax, .5)
																	}
																	
																	plt.SetColParams("EnviroAct", eplot.Off, eplot.FixMin, 0, eplot.FixMax, 1)
																	plt.SetColParams("InteroAct", eplot.Off, eplot.FixMin, 0, eplot.FixMax, 1)
																	plt.SetColParams("BehActM", eplot.Off, eplot.FixMin, 0, eplot.FixMax, 1)
																	plt.SetColParams("BehActP", eplot.Off, eplot.FixMin, 0, eplot.FixMax, 1)
																	return plt
																}
																
																//////////////////////////////////////////////
																//  TstEpcLog
																
																func (ss *Sim) LogTstEpc(dt *etable.Table) {
																	row := dt.Rows
																	dt.SetNumRows(row + 1)
																	
																	trl := ss.TstTrlLog
																	tix := etable.NewIdxView(trl)
																	epc := ss.TrainEnv.Epoch.Prv // ?
																	
																	// note: this shows how to use agg methods to compute summary data from another
																	// data table, instead of incrementing on the Sim
																	dt.SetCellFloat("Run", row, float64(ss.TrainEnv.Run.Cur))
																	dt.SetCellFloat("Epoch", row, float64(epc))
																	dt.SetCellFloat("SSE", row, agg.Sum(tix, "SSE")[0])
																	dt.SetCellFloat("AvgSSE", row, agg.Mean(tix, "AvgSSE")[0])
																	dt.SetCellFloat("PctErr", row, agg.Mean(tix, "Err")[0])
																	dt.SetCellFloat("PctCor", row, 1-agg.Mean(tix, "Err")[0])
																	dt.SetCellFloat("CosDiff", row, agg.Mean(tix, "CosDiff")[0])
																	
																	trlix := etable.NewIdxView(trl)
																	trlix.Filter(func(et *etable.Table, row int) bool {
																		return et.CellFloat("SSE", row) > 0 // include error trials
																	})
																	ss.TstErrLog = trlix.NewTable()
																	
																	// TODO: proceduralize this section.  shouldn't be too hard.
																	allsp := split.All(trlix)
																	split.Agg(allsp, "SSE", agg.AggSum)
																	split.Agg(allsp, "AvgSSE", agg.AggMean)
																	split.Agg(allsp, "EnviroAct", agg.AggMean)
																	split.Agg(allsp, "InteroAct", agg.AggMean)
																	split.Agg(allsp, "BehActM", agg.AggMean)
																	split.Agg(allsp, "BehActP", agg.AggMean)q
																	
																	ss.TstErrStats = allsp.AggsToTable(etable.AddAggName)
																	
																	// note: essential to use Go version of update when called from another goroutine
																	ss.TstEpcPlot.GoUpdate()
																}
																
																func (ss *Sim) ConfigTstEpcLog(dt *etable.Table) {
																	dt.SetMetaData("name", "TstEpcLog")
																	dt.SetMetaData("desc", "Summary stats for testing trials")
																	dt.SetMetaData("read-only", "true")
																	dt.SetMetaData("precision", strconv.Itoa(LogPrec))
																	
																	dt.SetFromSchema(etable.Schema{
																		{"Run", etensor.INT64, nil, nil},
																		{"Epoch", etensor.INT64, nil, nil},
																		{"SSE", etensor.FLOAT64, nil, nil},
																		{"AvgSSE", etensor.FLOAT64, nil, nil},
																		{"PctErr", etensor.FLOAT64, nil, nil},
																		{"PctCor", etensor.FLOAT64, nil, nil},
																		{"CosDiff", etensor.FLOAT64, nil, nil},
																		}, 0)
																	}
																	
																	func (ss *Sim) ConfigTstEpcPlot(plt *eplot.Plot2D, dt *etable.Table) *eplot.Plot2D {
																		plt.Params.Title = "Personality Model Testing Epoch Plot"
																		plt.Params.XAxisCol = "Epoch"
																		plt.SetTable(dt)
																		// order of params: on, fixMin, min, fixMax, max
																		plt.SetColParams("Run", eplot.Off, eplot.FixMin, 0, eplot.FloatMax, 0)
																		plt.SetColParams("Epoch", eplot.Off, eplot.FixMin, 0, eplot.FloatMax, 0)
																		plt.SetColParams("SSE", eplot.Off, eplot.FixMin, 0, eplot.FloatMax, 0)
																		plt.SetColParams("AvgSSE", eplot.Off, eplot.FixMin, 0, eplot.FloatMax, 0)
																		plt.SetColParams("PctErr", eplot.On, eplot.FixMin, 0, eplot.FixMax, 1) // default plot
																		plt.SetColParams("PctCor", eplot.On, eplot.FixMin, 0, eplot.FixMax, 1) // default plot
																		plt.SetColParams("CosDiff", eplot.Off, eplot.FixMin, 0, eplot.FixMax, 1)
																		return plt
																	}
																	
																	//////////////////////////////////////////////
																	//  TstCycLog
																	
																	// LogTstCyc adds data from current trial to the TstCycLog table.
																	// log just has 100 cycles, is overwritten
																	func (ss *Sim) LogTstCyc(dt *etable.Table, cyc int) {
																		if dt.Rows <= cyc {
																			dt.SetNumRows(cyc + 1)
																		}
																		
																		dt.SetCellFloat("Cycle", cyc, float64(cyc))
																		for _, lnm := range ss.LayStatNms {
																			ly := ss.Net.LayerByName(lnm).(leabra.LeabraLayer).AsLeabra()
																			dt.SetCellFloat(ly.Nm+" Ge.Avg", cyc, float64(ly.Pools[0].Inhib.Ge.Avg))
																			dt.SetCellFloat(ly.Nm+" Act.Avg", cyc, float64(ly.Pools[0].Inhib.Act.Avg))
																		}
																		
																		if cyc%10 == 0 { // too slow to do every cyc
																			// note: essential to use Go version of update when called from another goroutine
																			ss.TstCycPlot.GoUpdate()
																		}
																	}
																	
																	func (ss *Sim) ConfigTstCycLog(dt *etable.Table) {
																		dt.SetMetaData("name", "TstCycLog")
																		dt.SetMetaData("desc", "Record of activity etc over one trial by cycle")
																		dt.SetMetaData("read-only", "true")
																		dt.SetMetaData("precision", strconv.Itoa(LogPrec))
																		
																		np := 100 // max cycles
																		sch := etable.Schema{
																			{"Cycle", etensor.INT64, nil, nil},
																		}
																		for _, lnm := range ss.LayStatNms {
																			sch = append(sch, etable.Column{lnm + " Ge.Avg", etensor.FLOAT64, nil, nil})
																			sch = append(sch, etable.Column{lnm + " Act.Avg", etensor.FLOAT64, nil, nil})
																		}
																		dt.SetFromSchema(sch, np)
																	}
																	
																	func (ss *Sim) ConfigTstCycPlot(plt *eplot.Plot2D, dt *etable.Table) *eplot.Plot2D {
																		plt.Params.Title = "Personality Model Test Cycle Plot"
																		plt.Params.XAxisCol = "Cycle"
																		plt.SetTable(dt)
																		// order of params: on, fixMin, min, fixMax, max
																		plt.SetColParams("Cycle", eplot.Off, eplot.FixMin, 0, eplot.FloatMax, 0)
																		for _, lnm := range ss.LayStatNms {
																			plt.SetColParams(lnm+" Ge.Avg", true, true, 0, true, .5)
																			plt.SetColParams(lnm+" Act.Avg", true, true, 0, true, .5)
																		}
																		return plt
																	}
																	
																	//////////////////////////////////////////////
																	//  RunLog
																	
																	// LogRun adds data from current run to the RunLog table.
																	func (ss *Sim) LogRun(dt *etable.Table) {
																		run := ss.TrainEnv.Run.Cur // this is NOT triggered by increment yet -- use Cur
																		row := dt.Rows
																		dt.SetNumRows(row + 1)
																		
																		epclog := ss.TrnEpcLog
																		epcix := etable.NewIdxView(epclog)
																		// compute mean over last N epochs for run level
																		nlast := 5
																		if nlast > epcix.Len()-1 {
																			nlast = epcix.Len() - 1
																		}
																		epcix.Idxs = epcix.Idxs[epcix.Len()-nlast:]
																		
																		params := ss.RunName() // includes tag
																		
																		dt.SetCellFloat("Run", row, float64(run))
																		dt.SetCellString("Params", row, params)
																		dt.SetCellFloat("FirstZero", row, float64(ss.FirstZero))
																		dt.SetCellFloat("SSE", row, agg.Mean(epcix, "SSE")[0])
																		dt.SetCellFloat("AvgSSE", row, agg.Mean(epcix, "AvgSSE")[0])
																		dt.SetCellFloat("PctErr", row, agg.Mean(epcix, "PctErr")[0])
																		dt.SetCellFloat("PctCor", row, agg.Mean(epcix, "PctCor")[0])
																		dt.SetCellFloat("CosDiff", row, agg.Mean(epcix, "CosDiff")[0])
																		
																		runix := etable.NewIdxView(dt)
																		spl := split.GroupBy(runix, []string{"Params"})
																		split.Desc(spl, "FirstZero")
																		split.Desc(spl, "PctCor")
																		ss.RunStats = spl.AggsToTable(etable.AddAggName)
																		
																		// note: essential to use Go version of update when called from another goroutine
																		ss.RunPlot.GoUpdate()
																		if ss.RunFile != nil {
																			if row == 0 {
																				dt.WriteCSVHeaders(ss.RunFile, etable.Tab)
																			}
																			dt.WriteCSVRow(ss.RunFile, row, etable.Tab)
																		}
																	}
																	
																	func (ss *Sim) ConfigRunLog(dt *etable.Table) {
																		dt.SetMetaData("name", "RunLog")
																		dt.SetMetaData("desc", "Record of performance at end of training")
																		dt.SetMetaData("read-only", "true")
																		dt.SetMetaData("precision", strconv.Itoa(LogPrec))
																		
																		dt.SetFromSchema(etable.Schema{
																			{"Run", etensor.INT64, nil, nil},
																			{"Params", etensor.STRING, nil, nil},
																			{"FirstZero", etensor.FLOAT64, nil, nil},
																			{"SSE", etensor.FLOAT64, nil, nil},
																			{"AvgSSE", etensor.FLOAT64, nil, nil},
																			{"PctErr", etensor.FLOAT64, nil, nil},
																			{"PctCor", etensor.FLOAT64, nil, nil},
																			{"CosDiff", etensor.FLOAT64, nil, nil},
																			}, 0)
																		}
																		
																		func (ss *Sim) ConfigRunPlot(plt *eplot.Plot2D, dt *etable.Table) *eplot.Plot2D {
																			plt.Params.Title = "Leabra Random Associator 25 Run Plot"
																			plt.Params.XAxisCol = "Run"
																			plt.SetTable(dt)
																			// order of params: on, fixMin, min, fixMax, max
																			plt.SetColParams("Run", eplot.Off, eplot.FixMin, 0, eplot.FloatMax, 0)
																			plt.SetColParams("FirstZero", eplot.On, eplot.FixMin, 0, eplot.FloatMax, 0) // default plot
																			plt.SetColParams("SSE", eplot.Off, eplot.FixMin, 0, eplot.FloatMax, 0)
																			plt.SetColParams("AvgSSE", eplot.Off, eplot.FixMin, 0, eplot.FloatMax, 0)
																			plt.SetColParams("PctErr", eplot.Off, eplot.FixMin, 0, eplot.FixMax, 1)
																			plt.SetColParams("PctCor", eplot.Off, eplot.FixMin, 0, eplot.FixMax, 1)
																			plt.SetColParams("CosDiff", eplot.Off, eplot.FixMin, 0, eplot.FixMax, 1)
																			return plt
																		}
																		
																		////////////////////////////////////////////////////////////////////////////////////////////
																		// 		Gui
																		
																		// ADDING UNIT LABELS TO NETVIEW
																		
																		func (ss *Sim) ConfigNetView(nv *netview.NetView) {
																			// TODO: proceduralize this section.  shouldn't be too hard.  discuss considerations with dr. read if necessary.
																			
																			nv.ViewDefaults()
																			nv.Scene().Camera.Pose.Pos.Set(0.1, 1.5, 4)
																			nv.Scene().Camera.LookAt(mat32.Vec3{0.1, 0.1, 0}, mat32.Vec3{0, 1, 0})
																			labs := []string{"Frnd Desk Food Mate Bed SocSit Dngr", "IAff IAch IHngr ISex ISlp ISanx Fear", "WAff WAch WFood WSex WSlp", " WAvRej WAvHrm", "Hngt Stdy Eat Sex Sleep AvSoc Lve SHngt SStdy SEat SSlp SSex"}
																			
																			nv.ConfigLabels(labs)
																			
																			lays := []string{"Environment", "InteroState", "Approach", "Avoid","Behavior"}
																			
																			for li, lnm := range lays {
																				ly := nv.LayerByName(lnm)
																				lbl := nv.LabelByName(labs[li])
																				lbl.Pose = ly.Pose
																				lbl.Pose.Pos.Y += .1
																				lbl.Pose.Pos.Z += .02
																				lbl.Pose.Scale.SetMul(mat32.Vec3{.4, .08, 0.4})
																			}
																		}
																		
																		// ConfigGui configures the GoGi gui interface for this simulation,
																		func (ss *Sim) ConfigGui() *gi.Window {
																			width := 1600
																			height := 1200
																			
																			// gi.WinEventTrace = true
																			
																			gi.SetAppName("Person")
																			gi.SetAppAbout(`This is the Personality Dynamics Model.</p>`)
																			
																			win := gi.NewMainWindow("Person", "Personality Model", width, height)
																			ss.Win = win
																			
																			vp := win.WinViewport2D()
																			updt := vp.UpdateStart()
																			
																			mfr := win.SetMainFrame()
																			
																			tbar := gi.AddNewToolBar(mfr, "tbar")
																			tbar.SetStretchMaxWidth()
																			ss.ToolBar = tbar
																			
																			split := gi.AddNewSplitView(mfr, "split")
																			split.Dim = mat32.X
																			split.SetStretchMax()
																			
																			sv := giv.AddNewStructView(split, "sv")
																			sv.SetStruct(ss)
																			
																			tv := gi.AddNewTabView(split, "tv")
																			
																			nv := tv.AddNewTab(netview.KiT_NetView, "NetView").(*netview.NetView)
																			nv.Var = "Act"
																			// nv.Params.ColorMap = "Jet" // default is ColdHot
																			// which fares pretty well in terms of discussion here:
																			// https://matplotlib.org/tutorials/colors/colormaps.html
																			nv.SetNet(ss.Net)
																			ss.NetView = nv
																			ss.ConfigNetView(nv) // add labels etc
																			nv.Scene().Camera.Pose.Pos.Set(0, 1, 2.75) // more "head on" than default which is more "top down"
																			nv.Scene().Camera.LookAt(mat32.Vec3{0, 0, 0}, mat32.Vec3{0, 1, 0})
																			
																			plt := tv.AddNewTab(eplot.KiT_Plot2D, "TrnEpcPlot").(*eplot.Plot2D)
																			ss.TrnEpcPlot = ss.ConfigTrnEpcPlot(plt, ss.TrnEpcLog)
																			
																			plt = tv.AddNewTab(eplot.KiT_Plot2D, "TstTrlPlot").(*eplot.Plot2D)
																			ss.TstTrlPlot = ss.ConfigTstTrlPlot(plt, ss.TstTrlLog)
																			
																			plt = tv.AddNewTab(eplot.KiT_Plot2D, "TstCycPlot").(*eplot.Plot2D)
																			ss.TstCycPlot = ss.ConfigTstCycPlot(plt, ss.TstCycLog)
																			
																			plt = tv.AddNewTab(eplot.KiT_Plot2D, "TstEpcPlot").(*eplot.Plot2D)
																			ss.TstEpcPlot = ss.ConfigTstEpcPlot(plt, ss.TstEpcLog)
																			
																			plt = tv.AddNewTab(eplot.KiT_Plot2D, "RunPlot").(*eplot.Plot2D)
																			ss.RunPlot = ss.ConfigRunPlot(plt, ss.RunLog)
																			
																			split.SetSplits(.3, .7)
																			
																			tbar.AddAction(gi.ActOpts{Label: "Init", Icon: "update", Tooltip: "Initialize everything including network weights, and start over.  Also applies current params.", UpdateFunc: func(act *gi.Action) {
																				act.SetActiveStateUpdt(!ss.IsRunning)
																				}}, win.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
																					ss.Init()
																					vp.SetNeedsFullRender()
																				})
																				tbar.AddAction(gi.ActOpts{Label: "TrainPIT", Icon: "run", Tooltip: "This runs the Pavlovian and Instrumental training in the sequence and number of Epochs defined by a file that is read in.",
																				UpdateFunc: func(act *gi.Action) {
																					act.SetActiveStateUpdt(!ss.IsRunning)
																					}}, win.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
																						if !ss.IsRunning {
																							ss.IsRunning = true
																							tbar.UpdateActions()
																							// ss.Train()
																							go ss.TrainPIT()
																						}
																					})
																					tbar.AddAction(gi.ActOpts{Label: "Train", Icon: "run", Tooltip: "Starts the network training, picking up from wherever it may have left off.  If not stopped, training will complete the specified number of Runs through the full number of Epochs of training, with testing automatically occuring at the specified interval.",
																					UpdateFunc: func(act *gi.Action) {
																						act.SetActiveStateUpdt(!ss.IsRunning)
																						}}, win.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
																							if !ss.IsRunning {
																								ss.IsRunning = true
																								tbar.UpdateActions()
																								// ss.Train()
																								go ss.Train()
																							}
																						})
																						
																						tbar.AddAction(gi.ActOpts{Label: "Stop", Icon: "stop", Tooltip: "Interrupts running.  Hitting Train again will pick back up where it left off.", UpdateFunc: func(act *gi.Action) {
																							act.SetActiveStateUpdt(ss.IsRunning)
																							}}, win.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
																								ss.Stop()
																							})
																							
																							tbar.AddAction(gi.ActOpts{Label: "Step Trial", Icon: "step-fwd", Tooltip: "Advances one training trial at a time.", UpdateFunc: func(act *gi.Action) {
																								act.SetActiveStateUpdt(!ss.IsRunning)
																								}}, win.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
																									if !ss.IsRunning {
																										ss.IsRunning = true
																										ss.TrainTrial()
																										ss.IsRunning = false
																										vp.SetNeedsFullRender()
																									}
																								})
																								
																								tbar.AddAction(gi.ActOpts{Label: "Step Epoch", Icon: "fast-fwd", Tooltip: "Advances one epoch (complete set of training patterns) at a time.", UpdateFunc: func(act *gi.Action) {
																									act.SetActiveStateUpdt(!ss.IsRunning)
																									}}, win.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
																										if !ss.IsRunning {
																											ss.IsRunning = true
																											tbar.UpdateActions()
																											go ss.TrainEpoch()
																										}
																									})
																									
																									tbar.AddAction(gi.ActOpts{Label: "Step Run", Icon: "fast-fwd", Tooltip: "Advances one full training Run at a time.", UpdateFunc: func(act *gi.Action) {
																										act.SetActiveStateUpdt(!ss.IsRunning)
																										}}, win.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
																											if !ss.IsRunning {
																												ss.IsRunning = true
																												tbar.UpdateActions()
																												go ss.TrainRun()
																											}
																										})
																										
																										tbar.AddSeparator("test")
																										
																										tbar.AddAction(gi.ActOpts{Label: "Test Trial", Icon: "step-fwd", Tooltip: "Runs the next testing trial.", UpdateFunc: func(act *gi.Action) {
																											act.SetActiveStateUpdt(!ss.IsRunning)
																											}}, win.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
																												if !ss.IsRunning {
																													ss.IsRunning = true
																													ss.TestTrial(false) // don't return on change -- wrap
																													ss.IsRunning = false
																													vp.SetNeedsFullRender()
																												}
																											})
																											
																											tbar.AddAction(gi.ActOpts{Label: "Test Item", Icon: "step-fwd", Tooltip: "Prompts for a specific input pattern name to run, and runs it in testing mode.", UpdateFunc: func(act *gi.Action) {
																												act.SetActiveStateUpdt(!ss.IsRunning)
																												}}, win.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
																													gi.StringPromptDialog(vp, "", "Test Item",
																													gi.DlgOpts{Title: "Test Item", Prompt: "Enter the Name of a given input pattern to test (case insensitive, contains given string."},
																													win.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
																														dlg := send.(*gi.Dialog)
																														if sig == int64(gi.DialogAccepted) {
																															val := gi.StringPromptDialogValue(dlg)
																															idxs := ss.TestEnv.Table.RowsByString("Name", val, etable.Contains, etable.IgnoreCase)
																															if len(idxs) == 0 {
																																gi.PromptDialog(nil, gi.DlgOpts{Title: "Name Not Found", Prompt: "No patterns found containing: " + val}, gi.AddOk, gi.NoCancel, nil, nil)
																																} else {
																																	if !ss.IsRunning {
																																		ss.IsRunning = true
																																		fmt.Printf("testing index: %v\n", idxs[0])
																																		ss.TestItem(idxs[0])
																																		ss.IsRunning = false
																																		vp.SetNeedsFullRender()
																																	}
																																}
																															}
																														})
																													})
																													
																													tbar.AddAction(gi.ActOpts{Label: "Test All", Icon: "fast-fwd", Tooltip: "Tests all of the testing trials.", UpdateFunc: func(act *gi.Action) {
																														act.SetActiveStateUpdt(!ss.IsRunning)
																														}}, win.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
																															if !ss.IsRunning {
																																ss.IsRunning = true
																																tbar.UpdateActions()
																																go ss.RunTestAll()
																															}
																														})
																														
																														tbar.AddSeparator("log")
																														
																														tbar.AddAction(gi.ActOpts{Label: "Reset RunLog", Icon: "reset", Tooltip: "Reset the accumulated log of all Runs, which are tagged with the ParamSet used"}, win.This(),
																														func(recv, send ki.Ki, sig int64, data interface{}) {
																															ss.RunLog.SetNumRows(0)
																															ss.RunPlot.Update()
																														})
																														
																														tbar.AddSeparator("misc")
																														
																														tbar.AddAction(gi.ActOpts{Label: "New Seed", Icon: "new", Tooltip: "Generate a new initial random seed to get different results.  By default, Init re-establishes the same initial seed every time."}, win.This(),
																														func(recv, send ki.Ki, sig int64, data interface{}) {
																															ss.NewRndSeed()
																														})
																														
																														tbar.AddAction(gi.ActOpts{Label: "README", Icon: "file-markdown", Tooltip: "Opens your browser on the README file that contains instructions for how to run this model."}, win.This(),
																														func(recv, send ki.Ki, sig int64, data interface{}) {
																															gi.OpenURL("https://github.com/emer/leabra/blob/master/examples/ra25/README.md")
																														})
																														
																														vp.UpdateEndNoSig(updt)
																														
																														// main menu
																														appnm := gi.AppName()
																														mmen := win.MainMenu
																														mmen.ConfigMenus([]string{appnm, "File", "Edit", "Window"})
																														
																														amen := win.MainMenu.ChildByName(appnm, 0).(*gi.Action)
																														amen.Menu.AddAppMenu(win)
																														
																														emen := win.MainMenu.ChildByName("Edit", 1).(*gi.Action)
																														emen.Menu.AddCopyCutPaste(win)
																														
																														// note: Command in shortcuts is automatically translated into Control for
																														// Linux, Windows or Meta for MacOS
																														// fmen := win.MainMenu.ChildByName("File", 0).(*gi.Action)
																														// fmen.Menu.AddAction(gi.ActOpts{Label: "Open", Shortcut: "Command+O"},
																														// 	win.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
																															// 		FileViewOpenSVG(vp)
																															// 	})
																															// fmen.Menu.AddSeparator("csep")
																															// fmen.Menu.AddAction(gi.ActOpts{Label: "Close Window", Shortcut: "Command+W"},
																															// 	win.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
																																// 		win.Close()
																																// 	})
																																
																																inQuitPrompt := false
																																gi.SetQuitReqFunc(func() {
																																	if inQuitPrompt {
																																		return
																																	}
																																	inQuitPrompt = true
																																	gi.PromptDialog(vp, gi.DlgOpts{Title: "Really Quit?",
																																	Prompt: "Are you <i>sure</i> you want to quit and lose any unsaved params, weights, logs, etc?"}, gi.AddOk, gi.AddCancel,
																																	win.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
																																		if sig == int64(gi.DialogAccepted) {
																																			gi.Quit()
																																			} else {
																																				inQuitPrompt = false
																																			}
																																		})
																																	})
																																	
																																	// gi.SetQuitCleanFunc(func() {
																																		// 	fmt.Printf("Doing final Quit cleanup here..\n")
																																		// })
																																		
																																		inClosePrompt := false
																																		win.SetCloseReqFunc(func(w *gi.Window) {
																																			if inClosePrompt {
																																				return
																																			}
																																			inClosePrompt = true
																																			gi.PromptDialog(vp, gi.DlgOpts{Title: "Really Close Window?",
																																			Prompt: "Are you <i>sure</i> you want to close the window?  This will Quit the App as well, losing all unsaved params, weights, logs, etc"}, gi.AddOk, gi.AddCancel,
																																			win.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
																																				if sig == int64(gi.DialogAccepted) {
																																					gi.Quit()
																																					} else {
																																						inClosePrompt = false
																																					}
																																				})
																																			})
																																			
																																			win.SetCloseCleanFunc(func(w *gi.Window) {
																																				go gi.Quit() // once main window is closed, quit
																																			})
																																			
																																			win.MainMenuUpdated()
																																			return win
																																		}
																																		
																																		// These props register Save methods so they can be used
																																		var SimProps = ki.Props{
																																			"CallMethods": ki.PropSlice{
																																				{"SaveWeights", ki.Props{
																																					"desc": "save network weights to file",
																																					"icon": "file-save",
																																					"Args": ki.PropSlice{
																																						{"File Name", ki.Props{
																																							"ext": ".wts,.wts.gz",
																																						}},
																																					},
																																				}},
																																			},
																																		}
																																		
																																		func (ss *Sim) CmdArgs() {
																																			ss.NoGui = true
																																			var nogui bool
																																			var saveEpcLog bool
																																			var saveRunLog bool
																																			var note string
																																			flag.StringVar(&ss.ParamSet, "params", "", "ParamSet name to use -- must be valid name as listed in compiled-in params or loaded params")
																																			flag.StringVar(&ss.Tag, "tag", "", "extra tag to add to file names saved from this run")
																																			flag.StringVar(&note, "note", "", "user note -- describe the run params etc")
																																			flag.IntVar(&ss.MaxRuns, "runs", 10, "number of runs to do (note that MaxEpcs is in paramset)")
																																			flag.BoolVar(&ss.LogSetParams, "setparams", false, "if true, print a record of each parameter that is set")
																																			flag.BoolVar(&ss.SaveWts, "wts", false, "if true, save final weights after each run")
																																			flag.BoolVar(&saveEpcLog, "epclog", true, "if true, save train epoch log to file")
																																			flag.BoolVar(&saveRunLog, "runlog", true, "if true, save run epoch log to file")
																																			flag.BoolVar(&nogui, "nogui", true, "if not passing any other args and want to run nogui, use nogui")
																																			flag.Parse()
																																			ss.Init()
																																			
																																			if note != "" {
																																				fmt.Printf("note: %s\n", note)
																																			}
																																			if ss.ParamSet != "" {
																																				fmt.Printf("Using ParamSet: %s\n", ss.ParamSet)
																																			}
																																			
																																			if saveEpcLog {
																																				var err error
																																				fnm := ss.LogFileName("epc")
																																				ss.TrnEpcFile, err = os.Create(fnm)
																																				if err != nil {
																																					log.Println(err)
																																					ss.TrnEpcFile = nil
																																					} else {
																																						fmt.Printf("Saving epoch log to: %v\n", fnm)
																																						defer ss.TrnEpcFile.Close()
																																					}
																																				}
																																				if saveRunLog {
																																					var err error
																																					fnm := ss.LogFileName("run")
																																					ss.RunFile, err = os.Create(fnm)
																																					if err != nil {
																																						log.Println(err)
																																						ss.RunFile = nil
																																						} else {
																																							fmt.Printf("Saving run log to: %v\n", fnm)
																																							defer ss.RunFile.Close()
																																						}
																																					}
																																					if ss.SaveWts {
																																						fmt.Printf("Saving final weights per run\n")
																																					}
																																					fmt.Printf("Running %d Runs\n", ss.MaxRuns)
																																					ss.Train()
																																				}
																																				