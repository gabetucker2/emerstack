Motivational Dynamics in the Everyday Life of a College Student

5/16/2019

Looks like it now learns the weights properly regardless of whether it is Pavlov then Instrumental or Instrumental then Pavlov. It looks like it may not save the weights for a layer if it is lesioned.

5/14/2019

Deleted Hidden2

04/18/2019

TO DO 
This currently does separate training for Instrumental and Pavlovian learning.

Make sure that the weights are saved between types of training, so that I have the correct weights for the overall network
NOTE: 
In this version of the model, we are doing training that parallels separate Pavlovian and Instrumental training. One run through trains the Pavlovian associations between Environment and Interoceptive State, and Motivation. A separate run through trains the Instrumental relationship between Motives and Behavior. Borrows from the code that Peter Wang used to train his PIT model where the type of the different layers is switched so that in Pavlovian, Environment and Interoceptive are Input Layers and Motive layers are Target. Then for Instrumental training, Motives are Input and Behavior is Target. Hidden layers remain set to Hidden and all other layers are set to Output, which does not participate in Learning.

Code has been written to programatically control the order in which the two kinds of learning are done by reading a DataTable that controls the order of learning and the number of epochs for each type of learning.

In this version of the model the weights from Environment and Interostate to Wanting are learned.

There is a direct connection from Interoceptive State to the Hidden2 Layer.

NOTE: 
This version of the model does not include the following characteristics of other versions of the model:

training events for seeking behaviors. Seeking connections are hard to learn and make it harder to model sequences of behavior over time.
values for the increment and decrement parameters for Seeking in the model. The variables are still there, but they are assigned to 0.
Goal of Model

The goal of this project is to model the dynamics of motivation and behavior over time. Researchers such as Berridge have argued that the degree to which we Want a particular Reward or Want to Avoid a particular Aversive state is a major driver of choice and behavior. Wanting is argued to be the multiplicative result of the current relevant bodily or interoceptive state (e.g, Hunger) and the availability of relevant Motive Affordances (e.g, Food) in the situation.

Dopamine. In addition, Wanting is also argued to be a function of Dopamine levels in the NAcc/Amygdala circuit. Higher dopamine, both phasic and tonic should lead to higher levels of Wanting. In the current model, Dopamine levels are modeled by inputs from a VTA_DA system.

Following this idea, the current network models the following things:

The joint, multiplicative role of Environmental Motive Affordances and Bodily or Interoceptive state in generating Wanting
The impact of differences in Tonic (and Phasic) dopmaine (DA) levels.
The role of Wanting in generating Behavior, where different Needs compete with each other in "selecting" a behavior.
The impact of the selected Behavior on:
Environmental Affordances and
Interoceptive State
The effect of resulting changes in Affordances and Interoceptive State on subsequent Wanting and Behavior
The concrete domain we use is the everyday life of a college student. We represent the major motives that would drive the behavior of a typical college student and the major Situational/Environmental Affordances that a typical college student would encounter in everyday life.
Structure of Network

The Environmental Affordances, Interoceptive States, Types of Wanting, and Behaviors in the current network are listed below

Environment

Frnd - Friend

Lbry - Library

Food - Food

Mate - Mate

Bed - Bed

SocSit - Social Situation

Dngr - Danger

Interoceptive State

nAff -nAffiliation

nAch - nAchievement

Hngr - Hunger

Sex - Sex

Slp -Sleep

SAnx - Social Anxiety

Fear - Fear

Approach

AFF - Affiliation

ACH - Achievement

HNGR - Hunger

SEX - Sex

SLP - Sleep

Avoid

REJ - Social Rejection

HRM - Physical Harm

Behavior

Hngt - Hangout

Stdy - Study

Eat - Eat

AvSoc - Avoid Socializing

Sex - Have Sex

Sleep - Sleep

Lve - Leave

SHngt - Seek Hangout

SStdy - Seek Study

SEat - Seek Eat

SSlp - Seek Sleep

SSex - Seek Sex

Structure of Program

Training

Training programs in LeabraAll_Std and relevant Training Data are used to teach the network six general things:

In the present version it learns to calculate a roughly multiplicative relationship between the relevant environmental feature and the corresponding interoceptive state in predicting the relevant behavior.
In this version the activation of WANTING nodes can be modified by a tonic input from DA/VTA neurons
There is a MotiveBias layer that can be used to manipulate individual differences in the baseline important of different Motives.
It learns a relationship between the strength of WANTING and the relevant Behavior.
NOTE, THIS IS NOT CURRENTLY IMPLEMENTED. When relevant environmental affordance is absent, but the Interoceptive State is high, it will teach the network to enact the relevant Seeking behavior.
NOT IMPLEMENTED. Over time, when the environmental cue is strong the network can learn "habitual" behavior that is not sensitive to current goal state.
Training Data file is called: 
WantToBeh -- trains weights between Motives and Behavior only. 
Training_Data_NoSeekingTraining (previously used Training_Data_Current). This file is structured to teach the relationship between Inputs and Motives.

Modeling Motivation and Behavior Over Time

Currently, all the logic for how Behavior changes Environmental Affordances and Interoceptive State is in the LeabraEpochTestKen program under LeabraAll_Test. Need to spend some time on the logic of this to see if I agree with everything as Ken implemented it.

Event is presented to the network from the input_data Table, currently named World_State.

Once the network settles, the Program identifies the mostly highly activated behavior from the Behavior layer and based on that behavior updates both the Environment, if relevant, and the corresponding Interoceptive state that is represented in World_State.

These modified values for the Environment and the Interoceptive State are then fed back into the network for the next event.

Environment and Interoceptive State are changed in three ways.

Some Interoceptive states simply change over time. For example, with the passage of time one gets hungrier or has a greater need for social affiliation. The program in EpochTest changes these states with each time step as a function of an increment that can be modified by the user. Different increments can be set for different Interoceptive States.
Environment and Interoceptive State can be changed by an enacted behavior. Eating changes both Environment and Interoceptive State. Hanging out with friends reduces the need for Affiliation. This is handled by the program using a parameter that is set by the user. Different parameters for each one.
Environments sometime changes because of State changes: Friend enters or leaves the situation. Weather changes. Fire alarm goes off. Intruder shows up with gun, etc.
Oftentimes there can be a delay in the impact of an action on Interoceptive state. For example, it takes a while for eating to change level of Hunger. This delay can be set indepependently for each Interoceptive State.
Currently, a time step in the program is equivalent to the presentation of a single state of the world to the network and the generation of a behavior.

This is the logic for presenting a sequence of events to the network.

DataTable assigned to input_data (World_State) is the one directly presented to the network. It represents the current State of the World, both internal and external. One can think of this as a blank slate at the beginning, where the initial values of all event are copied from input_data_load (WorldState_Changes) [see next line]
DataTable assigned to input_data_load (WorldState_Change) represents CHANGES in the State of the World over time. The first event in this table is the initial or starting state of the world. Rest of the table represents changes from this initial state.
WorldState_Changes DataTable represents when a feature appears in the environment and when it disappears from the environment. The feature is NOT repeated in the DataTable as long as it exists. This table ONLY represents CHANGE of STATE. If a feature enters, it is represented by a positive value between 0 and 1. If a feature leaves, such as a friend leaving the environment, then this is explicitly indicated by a -1, which reduces the presence of this feature down to 0 (Program limits values to range of 0 to 1). Thus, this table represents CHANGES in the state of the world. Represents when a feature appears and when it disappears. Once a feature appears and has been flagged by an entry in this table, there is no further entry until the feature disappears which is indicated by a -1.
Results from the current selected behavior and information from the input_data_load table about changes of State are integrated and then copied to the next row of the input_data table World_State, which represents the current State of The World. This is the next event seen by the network.
NEED TO CHECK MORE CAREFULLY INTO THE SPECIFICS OF HOW THIS IS DONE.
Information about Changes in State from the input_data_load table is copied and integrated with information about changes in Environment and Interoceptive as a result of the previous behavior.
This new information about current Environment and Interoceptive State is then copied to the next row of the input_data table.
This is the next event seen by the network.
The program uses the following variables, which can be set in the interface:

Variable descriptions

_c -- Counter variable (time). Passage of time. 
_d -- Delay variable. Amount of time before changes caused in Interostate by Behavior. For example, how long does it take for Hunger to start decreasing after the consumption of food? 
_delta -- amount of change in Interostate (or Environment) per relevant behavior (e.g., consummation). Could be either positive or negative. 
_incr -- amount of change in either Environment or Interostate as a result of passage of time. Typically positive, although some events could have negative effect. 
From Ken : _c is just a counter that for the passage of time allows the _d to function properly. It just tells you how close you are towards when the delay condition is met (delay condition is met when _c is at least equal to _d, so really it's the difference between _c and _d that is counting the time).

After delay is met you add either the interoceptive state or environment and _delta together.

Delta just means the change of either IS or ES after delay depending on if it's consumption or Seeking, and we agree that for consumption, only IS is changed after delay, and for SEEKing, only ES is changed. The only exception is Eating behavior, where both may change. Hunger should decrease only after delay, but availability of Food in the Environment should decrease without delay, which is something that I need to modify in the program, as right now both of them decrease only after delay.

_incr indicates increases in interceptive states that are not a direct result of behavior. Instead they occur simply due to the passage of time. They increase regardless of whether the delay condition is met. So for example, studying would decrease need for achievement after delay, but other needs should steadily increase. So when you start to study, needs like hunger, sleep, need for affiliation would steadily increase, and they should increase before delay is met. The only needs that won't steadily increase right now are avoid harm and social anxiety.